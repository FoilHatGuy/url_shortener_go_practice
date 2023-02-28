package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5" // add this
	"shortener/internal/cfg"
)

type databaseT struct {
	memory   DatabaseORM
	database *pgx.Conn
}

func databaseInitialize() DatabaseORM {
	fmt.Println(cfg.Storage.DatabaseDSN)
	db, err := pgx.Connect(context.Background(), cfg.Storage.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	return databaseT{memory: memory, database: db}
}
func (d databaseT) Initialize() {
	return
}

func (d databaseT) AddURL(s string, s2 string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d databaseT) GetURL(s string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d databaseT) GetURLByOwner(s string) ([]URLOfOwner, error) {
	//TODO implement me
	panic("implement me")
}

func (d databaseT) Ping() bool {
	err := d.database.Ping(context.Background())
	if err != nil {
		return false
	}
	return true
}
