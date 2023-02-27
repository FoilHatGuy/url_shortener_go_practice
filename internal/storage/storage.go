package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"shortener/internal/cfg"
	"shortener/internal/urlgenerator"
)

type dataT map[string]string
type ownerT map[string][]string

type storage struct {
	Data   dataT  `json:"data"`
	Owners ownerT `json:"owners"`
}

func RunAutosave() {
	//t := time.NewTicker(time.Duration(cfg.Storage.AutosaveInterval) * time.Second)
	Database.LoadData()
	//go func() {
	//	for range t.C {
	//		//fmt.Print("AUTOSAVE\n")
	//		Database.SaveData()
	//	}
	//}()
}

type DatabaseORM interface {
	AddURL(string, string) string
	GetURL(string) (string, error)
	GetURLByOwner(string) ([]URLOfOwner, error)
	SaveData()
	LoadData()
}

var Database DatabaseORM = storage{Data: make(dataT), Owners: make(ownerT)}

func (s storage) SaveData() {
	if cfg.Storage.StorageType == "none" {
		return
	}
	validateStruct(s)
	validateFolder()
	fmt.Print("SAVING\n")
	if data, err := json.Marshal(s); err == nil {
		//fmt.Printf("WRITING %v\n", data)
		err := os.WriteFile(cfg.Storage.SavePath, data, os.ModePerm)
		if err != nil {
			return
		}
		//fmt.Print("COMPLETE\n")
	}
}
func (s storage) LoadData() {
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

func (s storage) AddURL(url string, owner string) string {
	validateStruct(s)
	short := urlgenerator.RandSeq(cfg.Shortener.URLLength)
	s.Data[short] = url
	s.Owners[owner] = append(s.Owners[owner], short)
	s.SaveData()
	//s.shortURLs = append(s.shortURLs, short)
	return short
}

func (s storage) GetURL(url string) (string, error) {
	validateStruct(s)
	val, ok := s.Data[url]
	if ok {
		return val, nil
	}
	return "", errors.New("no url")
}

type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s storage) GetURLByOwner(owner string) ([]URLOfOwner, error) {
	var result []URLOfOwner
	for _, address := range s.Owners[owner] {
		fullAddr, err := url.JoinPath(cfg.Server.BaseURL, address)
		if err != nil {
			return nil, err
		}
		result = append(result, URLOfOwner{fullAddr, s.Data[address]})
	}
	fmt.Println("RESULT IN DB", result)

	return result, nil

}
