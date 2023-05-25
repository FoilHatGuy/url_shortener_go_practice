package storage

import (
	"context"
	"shortener/internal/cfg"
)

var Controller DatabaseORM

func Initialize() {
	switch cfg.Storage.StorageType {
	case "database":
		Controller = databaseInitialize()
		if Controller == nil {
			Controller = getMemoryController()
		}
	case "none":
		fallthrough
	case "file":
		Controller = getMemoryController()
	}
	Controller.Initialize()
}

type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type DatabaseORM interface {
	Initialize()
	AddURL(context.Context, string, string) (string, bool, error)
	GetURL(context.Context, string) (string, bool, error)
	GetURLByOwner(context.Context, string) ([]URLOfOwner, error)
	Ping(context.Context) bool
	Delete(context.Context, []string, string) error
}
