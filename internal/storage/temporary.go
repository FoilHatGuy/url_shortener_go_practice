package storage

import (
	"errors"
	"fmt"
	"shortener/internal/cfg"
	"shortener/internal/urlgenerator"
)

type dataT map[string]string

type storage struct {
	data dataT
}

type DatabaseORM interface {
	AddURL(string) string
	GetURL(string) (string, error)
	GetData() dataT
}

var Database DatabaseORM = storage{data: make(dataT)}

func (s storage) GetData() dataT {
	return s.data
}

func (s storage) AddURL(url string) string {
	if s.data == nil {
		s.data = make(dataT)
	}
	short := urlgenerator.RandSeq(cfg.Shortener.UrlLength)
	s.data[short] = url
	//s.shortURLs = append(s.shortURLs, short)
	return short
}

func (s storage) GetURL(url string) (string, error) {
	if s.data == nil {
		s.data = make(dataT)
	}
	val, ok := s.data[url]
	fmt.Print(ok, "\n")
	if ok {
		return val, nil
	}
	return "", errors.New("no url")
}
