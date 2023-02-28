package storage

import (
	"errors"
	"shortener/internal/cfg"
)

type controllerT struct {
	memory   DatabaseORM
	database DatabaseORM
}

var Controller DatabaseORM

func Initialize() {
	Controller = controllerT{
		memory:   memory,
		database: databaseInitialize(),
	}
	Controller.Initialize()
}

func (c controllerT) Ping() bool {
	switch cfg.Storage.StorageType {
	case "none":
		fallthrough
	case "file":
		//return c.memory.Ping()
		fallthrough
	case "database":
		return c.database.Ping()
	}
	return false
}

func (c controllerT) Initialize() {
	switch cfg.Storage.StorageType {
	case "none":
		fallthrough
	case "file":
		c.memory.Initialize()
		fallthrough
	case "database":
		//c.database =
		c.database.Initialize()
	}
}

func (c controllerT) AddURL(s string, s2 string) (string, error) {
	switch cfg.Storage.StorageType {
	case "database":
		fallthrough
		//return c.database.AddURL(s, s2)
	case "none":
		fallthrough
	case "file":
		return c.memory.AddURL(s, s2)
	}
	return "", errors.New("STORAGE_TYPE contains the value that is neither 'file' or 'database'")
}

func (c controllerT) GetURL(s string) (string, error) {
	switch cfg.Storage.StorageType {
	case "database":
		//return c.database.GetURL(s)
		fallthrough
	case "none":
		fallthrough
	case "file":
		return c.memory.GetURL(s)
	}
	return "", errors.New("STORAGE_TYPE contains the value that is neither 'file' or 'database'")
}

func (c controllerT) GetURLByOwner(s string) ([]URLOfOwner, error) {
	switch cfg.Storage.StorageType {
	case "database":
		fallthrough
		//return c.database.GetURLByOwner(s)
	case "none":
		fallthrough
	case "file":
		return c.memory.GetURLByOwner(s)
	}
	return nil, errors.New("STORAGE_TYPE contains the value that is neither 'file' or 'database'")
}

type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type DatabaseORM interface {
	Initialize()
	AddURL(string, string) (string, error)
	GetURL(string) (string, error)
	GetURLByOwner(string) ([]URLOfOwner, error)
	Ping() bool
}
