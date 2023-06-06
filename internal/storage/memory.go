package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/juliangruber/go-intersect"

	"shortener/internal/cfg"
)

type dataTVal struct {
	Original string
	Deleted  bool
}

func getMemoryController(config *cfg.ConfigT) DatabaseORM {
	return &storage{
		Data:   sync.Map{},
		Owners: sync.Map{},
		config: config,
	}
}

type storage struct {
	Data   sync.Map `json:"data"`   // map[string]dataTVal
	Owners sync.Map `json:"owners"` // map[string][]string
	config *cfg.ConfigT
}

// Initialize
// Performs initial setup of main operating variable using configuration from cfg.ConfigT
func (s *storage) Initialize() {
	s.loadData()
	// if there is a need to use autosave function:
	// interval := s.config.Storage.AutosaveInterval
	// if interval > 0 {
	//	go func(s *storage, inter int) {
	//		err := s.saveData()
	//		if err != nil {
	//			time.Sleep(10 * time.Second)
	//		}
	//		time.Sleep(time.Duration(inter) * time.Second)
	//	}(s, interval)
	//}
}

// func (s *storage) saveData() error {
//   validateFolder(s.config)
//   fmt.Print("SAVING\n")
//   if data, err := json.Marshal(s); err == nil {
//   	err := os.WriteFile(s.config.Storage.SavePath, data, os.ModePerm)
//   	if err != nil {
//   		return err
//   	}
//   }
//   return nil
// }

func (s *storage) loadData() {
	validateFolder(s.config)
	fmt.Printf("DATA LOADING\n")
	if file, err := os.ReadFile(s.config.Storage.SavePath); err == nil {
		err := json.Unmarshal(file, &s)
		fmt.Printf("LOADED URLS\n")
		if err != nil {
			return
		}
	}
}

func validateFolder(config *cfg.ConfigT) {
	if _, err := os.Stat(config.Storage.SavePath); os.IsNotExist(err) {
		fmt.Println("FOLDER DOESN'T EXIST, ")
		err := os.MkdirAll(filepath.Dir(config.Storage.SavePath), os.ModePerm)
		if err != nil {
			return
		}
	}
}

// AddURL adds a new entry to storage if it wasn't already added.
// Additionally, stores user key and all urls saved by each user
func (s *storage) AddURL(_ context.Context, original, short, user string) (ok bool, existing string, err error) {
	res := dataTVal{original, false}
	s.Data.Store(short, res)
	arr, ok := s.Owners.Load(user)
	if ok {
		s.Owners.Store(user, append(arr.([]string), short))
	} else {
		s.Owners.Store(user, []string{short})
	}
	return true, existing, nil
}

// GetURL retrieves original URL by its shortened form
func (s *storage) GetURL(_ context.Context, url string) (original string, ok bool, err error) {
	v, ok := s.Data.Load(url)
	if ok {
		val := v.(dataTVal)
		if val.Deleted {
			return "", true, nil
		}
		return val.Original, true, nil
	}
	return "", false, errors.New("no url")
}

// GetURLByOwner returns slice of URLOfOwner by user's uid
func (s *storage) GetURLByOwner(_ context.Context, owner string) (arrayURLs []URLOfOwner, err error) {
	var result []URLOfOwner
	user, ok := s.Owners.Load(owner)
	if !ok {
		return nil, nil
	}
	for _, address := range user.([]string) {
		fullAddr, err := url.JoinPath(s.config.Server.BaseURL, address)
		if err != nil {
			return nil, fmt.Errorf("while memory.GetURLByOwner %w", err)
		}
		v, ok := s.Data.Load(address)
		if ok {
			origURL := v.(dataTVal).Original
			result = append(result, URLOfOwner{fullAddr, origURL})
		}
	}
	return result, nil
}

// Delete marks url as deleted, and it will no longer be accessible by GetURL
func (s *storage) Delete(_ context.Context, urls []string, owner string) error {
	v, _ := s.Owners.Load(owner)
	if v == nil {
		return fmt.Errorf("no user")
	}
	userData := v.([]string)
	items := intersect.Hash(userData, urls)

	for _, i := range items {
		v, _ := s.Data.Load(i.(string))
		val := v.(dataTVal)
		val.Deleted = true
		s.Data.Store(i.(string), val)
	}

	return nil
}

// Ping checks the database availability
func (s *storage) Ping(_ context.Context) bool {
	return true
}
