package storage

import (
	"errors"
	"fmt"
	"urlgenerator"
)

// type longURLs []string
// type shortURLs []string
type dataT map[string]string

type storage struct {
	//longURLs
	//shortURLs
	data dataT
}

var Database = new(storage)

//Database.data =

func (s *storage) AddURL(url string, urlLength int) string {
	if Database.data == nil {
		Database.data = make(dataT)
	}
	short := urlgenerator.RandSeq(urlLength)
	s.data[short] = url
	//s.shortURLs = append(s.shortURLs, short)
	return short
}

func (s *storage) GetURL(url string) (string, error) {
	if Database.data == nil {
		Database.data = make(dataT)
	}
	val, ok := s.data[url]
	fmt.Print(ok, "\n")
	if ok {
		return val, nil
	}
	return "", errors.New("no url")
}
