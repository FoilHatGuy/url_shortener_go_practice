package storage

import (
	"errors"
	"fmt"
	"github.com/FoilHatGuy/url_shortener_go_practice/cmd/internal/urlGenerator"
)

// type longUrls []string
// type shortUrls []string
type dataT map[string]string

type storage struct {
	//longUrls
	//shortUrls
	data dataT
}

var Database = new(storage)

//Database.data =

func (s *storage) AddUrl(url string, urlLength int) string {
	if Database.data == nil {
		Database.data = make(dataT)
	}
	short := urlGenerator.RandSeq(urlLength)
	s.data[short] = url
	//s.shortUrls = append(s.shortUrls, short)
	return short
}

func (s *storage) GetUrl(url string) (string, error) {
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
