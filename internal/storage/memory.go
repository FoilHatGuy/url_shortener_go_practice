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
)

type dataTVal struct {
	Short   string
	Deleted bool
}
type dataT map[string]dataTVal
type ownerT map[string][]string

type storage struct {
	Data   dataT  `json:"data"`
	Owners ownerT `json:"owners"`
}

func (s storage) Delete(_ context.Context, urls []string, owner string) error {
	items := intersect.Hash(s.Owners[owner], urls)
	for _, i := range items {
		val := s.Data[i.(string)]
		val.Deleted = true
		s.Data[i.(string)] = val
	}
	return nil
}

func (s storage) Ping(_ context.Context) bool {
	return true
}

func (s storage) Initialize() {
	s.loadData()
}

var memory DatabaseORM = storage{Data: make(dataT), Owners: make(ownerT)}

func (s storage) saveData() error {
	if cfg.Storage.StorageType == "none" {
		return nil
	}
	validateStruct(s)
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
func (s storage) loadData() {
	if cfg.Storage.StorageType == "none" {
		return
	}
	validateStruct(s)
	validateFolder()
	fmt.Printf("DATA LOADING\n")
	if file, err := os.ReadFile(cfg.Storage.SavePath); err == nil {
		err := json.Unmarshal(file, &s)
		fmt.Printf("LOADED %d URLS\n", len(s.Data))
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
func validateStruct(s storage) {
	if s.Data == nil {
		s.Data = make(dataT)
	}
	if s.Owners == nil {
		s.Owners = make(ownerT)
	}
}

func (s storage) AddURL(_ context.Context, url string, owner string) (string, bool, error) {
	validateStruct(s)
	short := urlgenerator.RandSeq(cfg.Shortener.URLLength)
	res := dataTVal{url, false}
	s.Data[short] = res
	s.Owners[owner] = append(s.Owners[owner], short)
	//err := s.saveData()
	//if err != nil {
	//	return "", false, err
	//}
	//s.shortURLs = append(s.shortURLs, short)
	return short, true, nil
}

func (s storage) GetURL(_ context.Context, url string) (string, bool, error) {
	validateStruct(s)
	val, ok := s.Data[url]
	if ok {
		if val.Deleted {
			return "", true, nil
		}
		return val.Short, false, nil
	}
	return "", false, errors.New("no url")
}

func (s storage) GetURLByOwner(_ context.Context, owner string) ([]URLOfOwner, error) {
	var result []URLOfOwner
	for _, address := range s.Owners[owner] {
		fullAddr, err := url.JoinPath(cfg.Server.BaseURL, address)
		if err != nil {
			return nil, err
		}
		result = append(result, URLOfOwner{fullAddr, s.Data[address].Short})
	}

	return result, nil

}
