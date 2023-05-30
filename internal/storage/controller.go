package storage

import (
	"context"

	"shortener/internal/cfg"
)

// Controller is a main operating variable. To use, Initialize it first
var Controller DatabaseORM

// Initialize
// Performs initial setup of main operating variable using configuration from cfg.ConfigT
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

// URLOfOwner is a structure returned by DatabaseORM.GetURLByOwner method
type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// DatabaseORM
// Interface for realization of all used methods that need the database interactions. Can support multiple realizations.
type DatabaseORM interface {
	Initialize()
	AddURL(ctx context.Context, original string, short string, user string) (ok bool, existing string, err error)
	GetURL(ctx context.Context, short string) (original string, ok bool, err error)
	GetURLByOwner(ctx context.Context, owner string) (arrayURLs []URLOfOwner, err error)
	Ping(ctx context.Context) (ok bool)
	Delete(ctx context.Context, stringArray []string, owner string) (err error)
}
