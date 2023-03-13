package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"net/url"
	"regexp"
	"shortener/internal/cfg"
	"shortener/internal/urlgenerator"
	"strings"
)

type databaseT struct {
	database *pgx.Conn
}

func databaseInitialize() DatabaseORM {
	if cfg.Storage.DatabaseDSN == "" {
		return nil
	}
	fmt.Println("CREATING DATABASE")
	r := regexp.MustCompile(`dbname=[a-zA-Z0-9]*`)
	initAddress := r.ReplaceAllString(cfg.Storage.DatabaseDSN, "")
	fmt.Println(initAddress)
	initDB, err := pgx.Connect(context.Background(), initAddress)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	config, err := pgx.ParseConfig(cfg.Storage.DatabaseDSN)
	if err != nil {
		return nil
	}

	fmt.Println(config.Config.Database)

	_, err = initDB.Exec(context.Background(), `
		CREATE DATABASE 
	`+config.Config.Database)
	//
	//-- 		SELECT 'CREATE DATABASE shortener'
	//-- 		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = $1)
	//
	if err != nil {
		fmt.Println(err)
		//return nil
	}
	err = initDB.Close(context.Background())
	if err != nil {
		return nil
	}

	db, err := pgx.Connect(context.Background(), cfg.Storage.DatabaseDSN)
	if err != nil && !errors.Is(err, new(pgconn.PgError)) {
		return nil
	}
	return databaseT{database: db}
}
func (d databaseT) Initialize() {
	exec, err := d.database.Exec(context.Background(), `
	CREATE TABLE IF 	NOT EXISTS urls (
	    short_url 		text 	UNIQUE NOT NULL PRIMARY KEY, 
	    original_url 	text 	UNIQUE NOT NULL,
	    deleted			bool	NOT NULL DEFAULT FALSE
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

func (d databaseT) AddURL(ctx context.Context, url string, owner string) (string, bool, error) {
	short := urlgenerator.RandSeq(cfg.Shortener.URLLength)
	added := false
	var shortURL, originalURL string

	_, err := d.database.Exec(ctx, `
	INSERT INTO urls VALUES($1, $2, FALSE) 
	ON CONFLICT DO NOTHING
`, short, url)
	if err != nil {
		fmt.Println("ERR", err)
		return "", false, err
	}
	err = d.database.QueryRow(ctx, `
		SELECT short_url, original_url FROM urls
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

func (d databaseT) GetURL(ctx context.Context, short string) (string, bool, error) {
	var originalURL string
	var deleted bool
	err := d.database.QueryRow(ctx, `
		SELECT original_url, deleted FROM urls
		WHERE short_url = $1
	`, short).Scan(&originalURL, &deleted)
	fmt.Println(err, "\n", originalURL, "\n", short, "\n", deleted)
	if deleted {
		return "", true, nil
	}
	if err != nil {
		fmt.Println("ERR", err)
		return "", false, err
	}
	return originalURL, true, err
}

func (d databaseT) GetURLByOwner(ctx context.Context, owner string) ([]URLOfOwner, error) {
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

func (d databaseT) Delete(ctx context.Context, stringArray []string, owner string) error {
	fmt.Println(stringArray)
	q, err := d.database.Exec(ctx, fmt.Sprintf(`
		UPDATE urls
		SET deleted = TRUE
		WHERE short_url IN ('%s') AND short_url IN
		(
		    SELECT url FROM users WHERE user_id = $1
		)
	`, strings.Join(stringArray, "', '")), owner)
	fmt.Println("QQQQ\n", q)
	if err != nil {
		fmt.Println("ERR", err)
		return err
	}
	return nil
}

func (d databaseT) Ping(ctx context.Context) bool {
	err := d.database.Ping(ctx)
	return err == nil
}
