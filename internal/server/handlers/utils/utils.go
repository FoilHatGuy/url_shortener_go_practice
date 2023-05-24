package utils

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

func Shorten(inputURL string, owner string, ctx context.Context) (string, bool, error) {

	_, err := url.Parse(inputURL)
	if err != nil {
		return "", false, errors.New("bad URL")
	}
	shortURL, added, err := storage.Controller.AddURL(ctx, inputURL, owner)
	if err != nil {
		return "", added, err
	}

	fmt.Printf("Input url: %s\n", inputURL)
	fmt.Printf("Short url: %s\n\n", shortURL)

	result, _ := url.Parse(cfg.Server.BaseURL)
	result = result.JoinPath(shortURL)
	return result.String(), added, nil
}
