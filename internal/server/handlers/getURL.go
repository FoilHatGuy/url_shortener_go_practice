package handlers

import (
	"context"
	"errors"
	"net/http"

	"shortener/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "shortener/internal/server/pb"

	"shortener/internal/cfg"

	"github.com/gin-gonic/gin"
)

// GetURL
// Get the original url by the short url
func (s *ServerHTTP) GetURL(ctx *gin.Context) {
	inputURL := ctx.Params.ByName("shortURL")
	result, err := getURL(ctx, s.Database, s.Config, inputURL)
	switch {
	case errors.Is(err, errDeleted):
		ctx.Status(http.StatusGone)
	case errors.Is(err, errInvalidInput):
		ctx.Status(http.StatusBadRequest)
	default:
		ctx.Redirect(307, result)
	}
}

// GetURL
// Get the original url by the short url
func (s *ServerGRPC) GetURL(ctx context.Context, in *pb.GetURLIn) (out *pb.GetURLOut, errRPC error) {
	inputURL := in.GetInputURL()

	result, err := getURL(ctx, s.Database, s.Config, inputURL)

	switch {
	case errors.Is(err, errDeleted):
		errRPC = status.Errorf(codes.NotFound, "")
	case errors.Is(err, errInvalidInput):
		errRPC = status.Errorf(codes.InvalidArgument, "")
	default:
		out = &pb.GetURLOut{ResultURL: result}
	}
	return
}

func getURL(
	ctx context.Context,
	database storage.DatabaseORM,
	config *cfg.ConfigT,
	inputURL string,
) (result string, err error) {
	if len(inputURL) != config.Shortener.URLLength {
		return "", errInvalidInput
	}

	result, ok, err := database.GetURL(ctx, inputURL)
	if err != nil {
		return "", errInvalidInput
	}
	if result == "" && ok {
		return "", errDeleted
	}
	return result, nil
}
