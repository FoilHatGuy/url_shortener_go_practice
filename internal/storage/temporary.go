package storage

import (
	"errors"
	"fmt"
	"shortener/internal/urlgenerator"
)

// type longURLs []string
// type shortURLs []string
type dataT map[string]string

type storage struct {
	data dataT
}

type DatabaseORM interface {
	AddURL(string, int) string
	GetURL(string) (string, error)
}

var Database DatabaseORM = storage{data: make(dataT)}

func (s storage) AddURL(url string, urlLength int) string {
	if s.data == nil {
		s.data = make(dataT)
	}
	short := urlgenerator.RandSeq(urlLength)
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
