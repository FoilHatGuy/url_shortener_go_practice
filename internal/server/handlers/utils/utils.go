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

func Shorten(ctx context.Context, inputURL string, owner string) (string, bool, error) {
	if config == nil {
		config = cfg.Initialize()
	}

	_, err := url.Parse(inputURL)
	if err != nil {
		return "", false, errors.New("bad URL")
	}

	shortURL := RandSeq(config.Shortener.URLLength)
	added, err := storage.Controller.AddURL(ctx, inputURL, shortURL, owner)
	if err != nil {
		return "", added, err
	}

	fmt.Printf("Input url: %s\n", inputURL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	result, _ := url.Parse(config.Server.BaseURL)
	result = result.JoinPath(shortURL)
	return result.String(), added, nil
}
