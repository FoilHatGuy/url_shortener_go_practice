package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shortener/internal/auth"
	"shortener/internal/storage"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"shortener/internal/cfg"
	"shortener/internal/server/handlers"
	"shortener/internal/server/middleware"
)

// Run
// Performs initial setup of server router and launches it.
func Run(config *cfg.ConfigT) {
	dbController := storage.New(config)
	authEngine := auth.New(config)
	r := gin.Default()
	baseRouter := r.Group("")

	baseRouter.Use(middleware.Gzip())
	baseRouter.Use(middleware.Gunzip())
	baseRouter.Use(middleware.Cooker(config, authEngine))
	baseRouter.GET("/:shortURL",
		handlers.GetShortURL(dbController, config))
	baseRouter.GET("/ping", handlers.PingDatabase(dbController))
	baseRouter.POST("/", handlers.PostURL(dbController, config))

	api := baseRouter.Group("/api")
	api.POST("/shorten", handlers.PostAPIURL(dbController, config))
	api.POST("/shorten/batch", handlers.BatchShorten(dbController, config))
	api.GET("/user/urls", handlers.GetAllOwnedURL(dbController))
	api.DELETE("/user/urls", handlers.DeleteLine(dbController))

	pprof.Register(r)

	// end of handlers' declaration
	srv := &http.Server{
		Addr:              config.Server.Address,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}

	go func() { // run server in separate goroutine
		fmt.Println("SERVER LISTENING ON", config.Server.Address)
		if !config.Server.IsHTTPS {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		} // else

		certPEM, certKey, err := authEngine.GetCertificate()
		if err != nil {
			log.Fatalf("certificate generation failed: %w: %s\n", err)
		}
		if err := srv.ListenAndServeTLS(certPEM, certKey); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown setup
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-shutdown // wait for signal on channel
	log.Println("Initiating graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.Fatal("Server Shutdown:", err)
	}
	cancel()
	log.Println("Server exiting")
	os.Exit(0)
}
