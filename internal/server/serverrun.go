package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "shortener/internal/server/pb"

	"shortener/internal/auth"
	"shortener/internal/storage"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"shortener/internal/cfg"
	"shortener/internal/server/handlers"
)

// RunHTTP
// Performs initial setup of server router and launches it.
func RunHTTP(config *cfg.ConfigT) {
	srv := handlers.ServerHTTP{
		Database: storage.New(config),
		Security: auth.New(config),
		Config:   config,
	}

	r := gin.Default()
	baseRouter := r.Group("")

	baseRouter.Use(handlers.Gzip())
	baseRouter.Use(handlers.Gunzip())
	baseRouter.Use(srv.Cooker)
	baseRouter.GET("/:shortURL", srv.GetURL)
	baseRouter.GET("/ping", srv.PingDatabase)
	baseRouter.POST("/", srv.PostURL)

	api := baseRouter.Group("/api")
	api.POST("/shorten", srv.PostAPIURL)
	api.POST("/shorten/batch", srv.BatchShorten)
	api.GET("/user/urls", srv.GetAllOwnedURL)
	api.DELETE("/user/urls", srv.DeleteURLs)
	api.POST("/internal/stats", srv.GetStats)

	pprof.Register(r)

	// end of handlers' declaration
	srv.Server = http.Server{
		Addr:              config.Server.Address,
		Handler:           r,
		ReadHeaderTimeout: time.Second,
	}

	go func() { // run server in separate goroutine
		fmt.Println("SERVER LISTENING ON", config.Server.Address)
		if !config.Server.IsHTTPS {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("listen: %s\n", err)
			}
		} // else

		certPEM, certKey, err := srv.Security.GetCertificate()
		if err != nil {
			log.Fatalf("certificate generation failed: %s\n", err)
		}
		if err := srv.ListenAndServeTLS(certPEM, certKey); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

func RunGRPC(config *cfg.ConfigT) {
	srv := handlers.ServerGRPC{
		UnimplementedShortenerServer: pb.UnimplementedShortenerServer{},
		Database:                     storage.New(config),
		Security:                     auth.New(config),
		Config:                       config,
	}
	lis, err := net.Listen("tcp", config.Server.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(srv.Cooker))
	pb.RegisterShortenerServer(s, &srv)
	log.Printf("GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
