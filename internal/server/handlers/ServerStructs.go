package handlers

import (
	"errors"
	"net/http"

	pb "shortener/internal/server/pb"

	"shortener/internal/auth"
	"shortener/internal/cfg"
	"shortener/internal/storage"
)

var (
	errDeleted      = errors.New("required field was deleted")
	errInvalidInput = errors.New("input data is erroneous")
)

type ServerHTTP struct {
	http.Server
	Database storage.DatabaseORM
	Security *auth.EngineT
	Config   *cfg.ConfigT
}

type ServerGRPC struct {
	pb.UnimplementedShortenerServer
	Database storage.DatabaseORM
	Security *auth.EngineT
	Config   *cfg.ConfigT
}

type ownerKeyT string

const ownerCtxKey ownerKeyT = "owner"
