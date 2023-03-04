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

//
//func (c controllerT) Ping() bool {
//	switch cfg.Storage.StorageType {
//	case "none":
//		fallthrough
//	case "file":
//		return c.memory.Ping()
//		//fallthrough
//	case "database":
//		return c.database.Ping()
//	}
//	return false
//}
//
//func (c controllerT) Initialize() {
//	switch cfg.Storage.StorageType {
//	case "none":
//		fallthrough
//	case "file":
//		c.memory.Initialize()
//		//fallthrough
//	case "database":
//		c.database.Initialize()
//	}
//}
//
//func (c controllerT) AddURL(s string, s2 string) (string, bool, error) {
//	switch cfg.Storage.StorageType {
//	case "database":
//		return c.database.AddURL(s, s2)
//	case "none":
//		fallthrough
//	case "file":
//		return c.memory.AddURL(s, s2)
//	}
//	return "", false, errors.New("STORAGE_TYPE contains the value that is neither 'file' or 'database'")
//}
//
//func (c controllerT) GetURL(s string) (string, error) {
//	switch cfg.Storage.StorageType {
//	case "database":
//		return c.database.GetURL(s)
//	case "none":
//		fallthrough
//	case "file":
//		return c.memory.GetURL(s)
//	}
//	return "", errors.New("STORAGE_TYPE contains the value that is neither 'file' or 'database'")
//}
//
//func (c controllerT) GetURLByOwner(s string) ([]URLOfOwner, error) {
//	switch cfg.Storage.StorageType {
//	case "database":
//		return c.database.GetURLByOwner(s)
//	case "none":
//		fallthrough
//	case "file":
//		return c.memory.GetURLByOwner(s)
//	}
//	return nil, errors.New("STORAGE_TYPE contains the value that is neither 'file' or 'database'")
//}

type URLOfOwner struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
type DatabaseORM interface {
	Initialize()
	AddURL(string, string, context.Context) (string, bool, error)
	GetURL(string, context.Context) (string, error)
	GetURLByOwner(string, context.Context) ([]URLOfOwner, error)
	Ping(ctx context.Context) bool
}
