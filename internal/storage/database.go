package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"shortener/internal/cfg"
)

type poolWrapper struct {
	pool *pgxpool.Pool
}

// Close is a wrapper function for the Close method of any available pg connection
func (p poolWrapper) Close() {
	p.pool.Close()
}

// Exec is a wrapper function for the Exec method of any available pg connection
func (p poolWrapper) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, arguments...) //nolint:wrapcheck
}

// Query is a wrapper function for the Query method of any available pg connection
func (p poolWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...) //nolint:wrapcheck
}

// QueryRow is a wrapper function for the QueryRow method of any available pg connection
func (p poolWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

// Ping is a wrapper function for the Ping method of any available pg connection
func (p poolWrapper) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx) //nolint:wrapcheck
}

type databaseI interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Ping(ctx context.Context) error
	Close()
}

type databaseT struct {
	database databaseI
	config   *cfg.ConfigT
}

// databaseInitialize
// Performs initial setup of main operating variable using configuration from cfg.ConfigT.
// Creates the database on postgres specified by cfg.StorageT DatabaseDSN field
func databaseInitialize(config *cfg.ConfigT) DatabaseORM {
	fmt.Println("CREATING DATABASE")
	r := regexp.MustCompile(`dbname=[a-zA-Z0-9]*`)
	initAddress := r.ReplaceAllString(config.Storage.DatabaseDSN, "")
	fmt.Println(initAddress)
	initDB, err := pgx.Connect(context.Background(), initAddress)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	dbConf, err := pgx.ParseConfig(config.Storage.DatabaseDSN)
	if err != nil {
		return nil
	}

	fmt.Println(dbConf.Config.Database)

	_, err = initDB.Exec(context.Background(), `
		CREATE DATABASE 
	`+dbConf.Config.Database)
	//
	//-- 		SELECT 'CREATE DATABASE shortener'
	//-- 		WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = $1)
	//
	if err != nil {
		fmt.Println(err)
		// return nil
	}
	err = initDB.Close(context.Background())
	if err != nil {
		return nil
	}

	db, err := pgxpool.New(context.Background(), config.Storage.DatabaseDSN)
	if err != nil && !errors.Is(err, new(pgconn.PgError)) {
		return nil
	}
	return &databaseT{
		database: &poolWrapper{db},
		config:   config,
	}
}

// Initialize creates all required tables and sets up relations
func (d *databaseT) Initialize() {
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

// AddURL adds a new entry to storage if it wasn't already added.
// users table stores user key and all urls saved by each user
func (d *databaseT) AddURL(ctx context.Context, original, short, user string) (ok bool, existing string, err error) {
	var shortURL, originalURL string

	_, err = d.database.Exec(ctx, `
	INSERT INTO urls VALUES($1, $2, FALSE) 
	ON CONFLICT DO NOTHING
`, short, original)
	if err != nil {
		fmt.Println("ERR", err)
		return false, "", fmt.Errorf("while database.AddURL %w", err)
	}
	err = d.database.QueryRow(ctx, `
		SELECT short_url, original_url FROM urls
		WHERE original_url = $1
	`, original).Scan(&shortURL, &originalURL)
	if err != nil {
		fmt.Println("ERR", err)
		return false, "", fmt.Errorf("while database.AddURL %w", err)
	}

	if short != shortURL {
		fmt.Println("RESULT:", shortURL, originalURL)
		return false, shortURL, nil
	}

	_, err = d.database.Exec(ctx, `
			INSERT INTO users VALUES($1, $2) 
		`, user, short)
	if err != nil {
		fmt.Println("ERR", err)
		return false, "", fmt.Errorf("while database.AddURL %w", err)
	}

	return true, "", nil
}

// GetURL retrieves original URL by its shortened form
func (d *databaseT) GetURL(ctx context.Context, short string) (original string, ok bool, err error) {
	var originalURL string
	var deleted bool
	err = d.database.QueryRow(ctx, `
		SELECT original_url, deleted FROM urls
		WHERE short_url = $1
	`, short).Scan(&originalURL, &deleted)
	fmt.Println(err, "\n", originalURL, "\n", short, "\n", deleted)
	if deleted {
		return "", true, nil
	}
	if err != nil {
		fmt.Println("ERR", err)
		return "", false, fmt.Errorf("while database.GetURL %w", err)
	}
	return originalURL, true, nil
}

// GetURLByOwner returns slice of URLOfOwner by user's uid
func (d *databaseT) GetURLByOwner(ctx context.Context, owner string) (arrayURLs []*URLOfOwner, err error) {
	rows, err := d.database.Query(ctx, `
		SELECT short_url, original_url FROM urls, users
		WHERE user_id = $1 AND short_url = url
	`, owner)
	if err != nil {
		fmt.Println("ERR", err)
		return nil, fmt.Errorf("while database.GetURLByOwner %w", err)
	}
	defer rows.Close()
	var originalURL, shortURL string
	for rows.Next() {
		err = rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, fmt.Errorf("while database.GetURLByOwner %w", err)
		}
		fullAddr, _ := url.JoinPath(d.config.Server.BaseURL, shortURL)
		arrayURLs = append(arrayURLs, &URLOfOwner{fullAddr, originalURL})
	}

	return arrayURLs, nil
}

// Delete marks url as deleted, and it will no longer be accessible by GetURL
func (d *databaseT) Delete(ctx context.Context, stringArray []string, owner string) error {
	_, err := d.database.Exec(ctx, fmt.Sprintf(`
		UPDATE urls
		SET deleted = TRUE
		WHERE short_url IN ('%s') AND short_url IN
		(
		    SELECT url FROM users WHERE user_id = $1
		)
	`, strings.Join(stringArray, "', '")), owner)
	if err != nil {
		fmt.Println("ERR", err)
		return fmt.Errorf("while database.Delete %w", err)
	}
	return nil
}

// GetStats returns the number of urls and users
func (d *databaseT) GetStats(ctx context.Context) (stats StatsT, err error) {
	var countURLs, countUsers int64

	err = d.database.QueryRow(ctx, `
	SELECT COUNT(short_url) FROM urls
`).Scan(&countURLs)
	if err != nil {
		return StatsT{}, fmt.Errorf("while database.Stats %w", err)
	}

	err = d.database.QueryRow(ctx, `
	SELECT COUNT(DISTINCT user_id) FROM owners
`).Scan(&countUsers)
	if err != nil {
		return StatsT{}, fmt.Errorf("while database.Stats %w", err)
	}

	return StatsT{
		URLs:  countURLs,
		Users: countUsers,
	}, nil
}

// Ping checks the database availability
func (d *databaseT) Ping(ctx context.Context) bool {
	err := d.database.Ping(ctx)
	return err == nil
}
