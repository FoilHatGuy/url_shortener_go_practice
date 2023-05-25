package storage

import (
	"context"
	"shortener/internal/cfg"
)

var Controller DatabaseORM

func Initialize(config *cfg.ConfigT) {
	switch config.Storage.StorageType {
	case "database":
		Controller = databaseInitialize(config)
		if Controller == nil {
			Controller = getMemoryController(config)
		}
	default:
		Controller = getMemoryController(config)
	}
	Controller.Initialize()
}

type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type DatabaseORM interface {
	Initialize()
	AddURL(ctx context.Context, original string, short string, user string) (ok bool, existing string, err error)
	GetURL(ctx context.Context, short string) (original string, ok bool, err error)
	GetURLByOwner(ctx context.Context, owner string) (URLList []URLOfOwner, err error)
	Ping(ctx context.Context) (ok bool)
	Delete(ctx context.Context, stringArray []string, owner string) (err error)
}
