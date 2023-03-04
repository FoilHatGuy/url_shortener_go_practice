package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"net/url"
	"regexp"
	"shortener/internal/cfg"
	"shortener/internal/urlgenerator"
)

type databaseT struct {
	database *pgx.Conn
}

func databaseInitialize() DatabaseORM {
	if cfg.Storage.DatabaseDSN == "" {
		return nil
	}
	r := regexp.MustCompile(`dbname=[a-zA-Z]*\\s?`)
	initAddress := r.ReplaceAllString(cfg.Storage.DatabaseDSN, "")
	initDB, err := pgx.Connect(context.Background(), initAddress)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	exec, err := initDB.Exec(context.Background(), `
		SELECT 'CREATE DATABASE shortener'
		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'shortener')
	`)
	fmt.Println(exec)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	db, err := pgx.Connect(context.Background(), cfg.Storage.DatabaseDSN)
	if err != nil {
		return nil
	}
	return databaseT{database: db}
}
func (d databaseT) Initialize() {
	exec, err := d.database.Exec(context.Background(), `
	CREATE TABLE IF 	NOT EXISTS urls (
	    short_url 		text 	UNIQUE NOT NULL, 
	    original_url 	text 	UNIQUE NOT NULL,
	    PRIMARY KEY (short_url)
	                                )
`)
	fmt.Println(exec)
	if err != nil {
		fmt.Println(err)
		return
	}
	exec, err = d.database.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users (
	    user_id 		text 	NOT NULL, 
	    url 			text 	UNIQUE NOT NULL,
	    FOREIGN KEY (url)
	    	REFERENCES urls (short_url)
	    	ON DELETE CASCADE
	                                )
`)
	fmt.Println(exec)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (d databaseT) AddURL(url string, owner string, ctx context.Context) (string, bool, error) {
	short := urlgenerator.RandSeq(cfg.Shortener.URLLength)
	added := false
	var shortURL, originalURL string

	_, err := d.database.Exec(ctx, `
	INSERT INTO urls VALUES($1, $2) 
	ON CONFLICT DO NOTHING
`, short, url)
	if err != nil {
		fmt.Println("ERR", err)
		return "", false, err
	}
	err = d.database.QueryRow(ctx, `
		SELECT * FROM urls
		WHERE original_url = $1
	`, url).Scan(&shortURL, &originalURL)
	if err != nil {
		fmt.Println("ERR", err)
		return "", false, err
	}

	if short != shortURL {
		fmt.Println("RESULT:", shortURL, originalURL)

	} else {
		_, err = d.database.Exec(ctx, `
			INSERT INTO users VALUES($1, $2) 
		`, owner, short)
		if err != nil {
			fmt.Println("ERR", err)
			return "", false, err
		}
		added = true
	}
	return shortURL, added, nil
}

func (d databaseT) GetURL(short string, ctx context.Context) (string, error) {
	var originalURL string
	err := d.database.QueryRow(ctx, `
		SELECT original_url FROM urls
		WHERE short_url = $1
	`, short).Scan(&originalURL)
	if err != nil {
		fmt.Println("ERR", err)
		return "", err
	}
	return originalURL, err
}

func (d databaseT) GetURLByOwner(owner string, ctx context.Context) ([]URLOfOwner, error) {
	rows, err := d.database.Query(ctx, `
		SELECT short_url, original_url FROM urls, users
		WHERE user_id = $1 AND short_url = url
	`, owner)
	if err != nil {
		fmt.Println("ERR", err)
		return nil, err
	}
	defer rows.Close()
	fmt.Println(rows)
	var originalURL, shortURL string
	var result []URLOfOwner
	for rows.Next() {
		err := rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}
		fullAddr, err := url.JoinPath(cfg.Server.BaseURL, shortURL)
		if err != nil {
			return nil, err
		}
		result = append(result, URLOfOwner{fullAddr, originalURL})
	}

	return result, err
}

func (d databaseT) Ping(ctx context.Context) bool {
	err := d.database.Ping(ctx)
	return err == nil
}
