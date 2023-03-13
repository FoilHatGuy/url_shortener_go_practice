package storage

import (
	"context"
	"shortener/internal/cfg"
)

type controllerT struct {
	memory   DatabaseORM
	database DatabaseORM
}

var Controller DatabaseORM

func Initialize() {
	switch cfg.Storage.StorageType {
	case "none":
		fallthrough
	case "file":
		Controller = memory
		//fallthrough
	case "database":
		Controller = databaseInitialize()
	}
	//Controller = controllerT{
	//	memory:   memory,
	//	database: databaseInitialize(),
	//}
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
