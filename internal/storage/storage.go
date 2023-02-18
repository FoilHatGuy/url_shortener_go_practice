package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"shortener/internal/cfg"
	"shortener/internal/urlgenerator"
	"time"
)

type dataT map[string]string

type storage struct {
	Data dataT `json:"data"`
}

func RunAutosave() {
	t := time.NewTicker(time.Duration(cfg.Storage.AutosaveInterval) * time.Second)
	Database.LoadData()
	go func() {
		for range t.C {
			//fmt.Print("AUTOSAVE\n")
			Database.SaveData()
		}
	}()
}

type DatabaseORM interface {
	AddURL(string) string
	GetURL(string) (string, error)
	SaveData()
	LoadData()
}

var Database DatabaseORM = storage{Data: make(dataT)}

func (s storage) SaveData() {
	if _, err := os.Stat(cfg.Storage.SavePath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(cfg.Storage.SavePath), os.ModePerm)
		if err != nil {
			return
		}
	}
	//fmt.Print("MARSHALLING\n")
	if data, err := json.Marshal(s.Data); err == nil {
		//fmt.Printf("WRITING %v\n", data)
		err := os.WriteFile(cfg.Storage.SavePath, data, os.ModePerm)
		if err != nil {
			return
		}
		//fmt.Print("COMPLETE\n")
	}
}
func (s storage) LoadData() {
	validateFolder()
	if file, err := os.ReadFile(cfg.Storage.SavePath); err == nil {
		err := json.Unmarshal(file, &s.Data)
		if err != nil {
			return
		}
	}
}

func validateFolder() {
	if _, err := os.Stat(cfg.Storage.SavePath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(cfg.Storage.SavePath), os.ModePerm)
		if err != nil {
			return
		}
	}
}

func (s storage) AddURL(url string) string {
	if s.Data == nil {
		s.Data = make(dataT)
	}
	short := urlgenerator.RandSeq(cfg.Shortener.URLLength)
	s.Data[short] = url
	s.SaveData()
	//s.shortURLs = append(s.shortURLs, short)
	return short
}

func (s storage) GetURL(url string) (string, error) {
	if s.Data == nil {
		s.Data = make(dataT)
	}
	val, ok := s.Data[url]
	fmt.Print(ok, "\n")
	if ok {
		return val, nil
	}
	return "", errors.New("no url")
}
