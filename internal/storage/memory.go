package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/juliangruber/go-intersect"
	"net/url"
	"os"
	"path/filepath"
	"shortener/internal/cfg"
	"shortener/internal/urlgenerator"
	"sync"
)

type dataTVal struct {
	Original string
	Deleted  bool
}

var memoryController DatabaseORM

func getMemoryController() DatabaseORM {
	if memoryController == nil {
		memoryController = &storage{Data: sync.Map{}, Owners: sync.Map{}}
	}
	return memoryController
}

type storage struct {
	Data   sync.Map `json:"data"`   // map[string]dataTVal
	Owners sync.Map `json:"owners"` // map[string][]string
}

func (s *storage) Delete(_ context.Context, urls []string, owner string) error {
	user, _ := s.Owners.Load(owner)
	items := intersect.Hash(user, urls)

	for _, i := range items {
		v, _ := s.Data.Load(i.(string))
		val := v.(dataTVal)
		val.Deleted = true
		s.Data.Store(i.(string), val)
	}
	return nil
}

func (s *storage) Ping(_ context.Context) bool {
	return true
}

func (s *storage) Initialize() {
	s.loadData()
}

func (s *storage) saveData() error {
	if cfg.Storage.StorageType == "none" {
		return nil
	}
	validateFolder()
	fmt.Print("SAVING\n")
	if data, err := json.Marshal(s); err == nil {
		//fmt.Printf("WRITING %v\n", data)
		err := os.WriteFile(cfg.Storage.SavePath, data, os.ModePerm)
		if err != nil {
			return err
		}
		//fmt.Print("COMPLETE\n")
	}
	return nil
}
func (s *storage) loadData() {
	if cfg.Storage.StorageType == "none" {
		return
	}
	validateFolder()
	fmt.Printf("DATA LOADING\n")
	if file, err := os.ReadFile(cfg.Storage.SavePath); err == nil {
		err := json.Unmarshal(file, &s)
		fmt.Printf("LOADED URLS\n")
		if err != nil {
			return
		}
	}
}

func validateFolder() {
	if _, err := os.Stat(cfg.Storage.SavePath); os.IsNotExist(err) {
		fmt.Println("FOLDER DOESN'T EXIST, ")
		err := os.MkdirAll(filepath.Dir(cfg.Storage.SavePath), os.ModePerm)
		if err != nil {
			return
		}
	}
}

func (s *storage) AddURL(_ context.Context, url string, owner string) (short string, added bool, err error) {
	short = urlgenerator.RandSeq(cfg.Shortener.URLLength)
	res := dataTVal{url, false}
	s.Data.Store(short, res)
	arr, ok := s.Owners.Load(owner)
	if ok {
		s.Owners.Store(owner, append(arr.([]string), short))
	} else {
		s.Owners.Store(owner, []string{short})
	}
	return short, true, nil
}

func (s *storage) GetURL(_ context.Context, url string) (original string, ok bool, err error) {
	v, ok := s.Data.Load(url)
	val := v.(dataTVal)
	if ok {
		if val.Deleted {
			return "", true, nil
		}
		return val.Original, false, nil
	}
	return "", false, errors.New("no url")
}

func (s *storage) GetURLByOwner(_ context.Context, owner string) ([]URLOfOwner, error) {
	var result []URLOfOwner
	user, _ := s.Owners.Load(owner)
	for _, address := range user.([]string) {
		fullAddr, err := url.JoinPath(cfg.Server.BaseURL, address)
		if err != nil {
			return nil, err
		}
		v, ok := s.Data.Load(address)
		if ok {
			origURL := v.(dataTVal).Original
			result = append(result, URLOfOwner{fullAddr, origURL})
		}
	}
	return result, nil
}
