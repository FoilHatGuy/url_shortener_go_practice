package storage

import "github.com/FoilHatGuy/url_shortener_go_practice/urlGenerator"

type longUrls []string
type shortUrls []string
type Storage struct {
	longUrls
	shortUrls
}

func (s *Storage) AddUrl(url string, urlLength int) string {
	s.longUrls = append(s.longUrls, url)
	short := urlGenerator.RandSeq(urlLength)
	s.shortUrls = append(s.shortUrls, short)
	return short
}
