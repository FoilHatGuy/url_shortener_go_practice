package storage

import (
	"context"

	"shortener/internal/cfg"
)

// New
// Performs initial setup of main operating variable using configuration from cfg.ConfigT
func New(config *cfg.ConfigT) DatabaseORM {
	var controller DatabaseORM
	if config.Storage.DatabaseDSN != "" {
		controller = databaseInitialize(config)
		if controller == nil {
			controller = getMemoryController(config)
		}
	} else {
		controller = getMemoryController(config)
	}
	controller.Initialize()
	return controller
}

// URLOfOwner is a structure returned by DatabaseORM.GetURLByOwner method
type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// StatsT is a structure returned by DatabaseORM.GetStats method
type StatsT struct {
	URLs  int64 `json:"urls"`
	Users int64 `json:"users"`
}

// DatabaseORM
// Interface for realization of all used methods that need the database interactions. Can support multiple realizations.
type DatabaseORM interface {
	Initialize()
	AddURL(ctx context.Context, original string, short string, user string) (ok bool, existing string, err error)
	GetURL(ctx context.Context, short string) (original string, ok bool, err error)
	GetURLByOwner(ctx context.Context, owner string) (arrayURLs []*URLOfOwner, err error)
	Ping(ctx context.Context) (ok bool)
	Delete(ctx context.Context, stringArray []string, owner string) (err error)
	GetStats(ctx context.Context) (stats StatsT, err error)
}
