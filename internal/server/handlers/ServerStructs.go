package handlers

import (
	"errors"
	"net/http"

	"shortener/internal/auth"
	"shortener/internal/cfg"
	pb "shortener/internal/server/pb"
	"shortener/internal/storage"
)

var (
	errDeleted      = errors.New("required field was deleted")
	errInvalidInput = errors.New("input data is erroneous")
)

// ServerHTTP is a structure containing all required services, as well as embedded server
type ServerHTTP struct {
	http.Server
	Database storage.DatabaseORM
	Security *auth.EngineT
	Config   *cfg.ConfigT
}

// ServerGRPC is a structure containing all required services, as well as embedded server
type ServerGRPC struct {
	pb.UnimplementedShortenerServer
	Database storage.DatabaseORM
	Security *auth.EngineT
	Config   *cfg.ConfigT
}

type ownerKeyT string

const ownerCtxKey ownerKeyT = "owner"
