package utils

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

var config *cfg.ConfigT

// Shorten is a common function used by handlers that perform adding URLs to database.
// Takes original URL and uid and returns the URL by which user can access their URL.
func Shorten(ctx context.Context, inputURL string, owner string) (result string, added bool, err error) {
	if config == nil {
		config = cfg.Initialize()
	}

	_, err = url.Parse(inputURL)
	if err != nil {
		return "", false, errors.New("bad URL")
	}

	shortURL := RandSeq(config.Shortener.URLLength)
	added, existing, err := storage.Controller.AddURL(ctx, inputURL, shortURL, owner)
	if err != nil {
		return "", added, err
	}
	if !added {
		shortURL = existing
	}

	fmt.Printf("Input url: %s\n", inputURL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	resultURL, _ := url.Parse(config.Server.BaseURL)
	result = resultURL.JoinPath(shortURL).String()
	return result, added, nil
}
