package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shortener/internal/server/handlers/utils"
	pb "shortener/internal/server/pb"
)

// PostURL
// Handler for batch shortening of urs.
// Takes txt representation of request body and returns txt of url for accessing it.
func (s *ServerHTTP) PostURL(ctx *gin.Context) {
	buf, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	inputURL := string(buf)
	owner, ok := ctx.Get(string(ownerCtxKey))
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	result, added, err := utils.Shorten(ctx, s.Database, inputURL, owner.(string), s.Config)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if added {
		ctx.String(http.StatusCreated, "%v", result)
	} else {
		ctx.String(http.StatusConflict, "%v", result)
	}
}

// PostURL
// Handler for batch shortening of urs.
// Takes txt representation of request body and returns txt of url for accessing it.
func (s *ServerGRPC) PostURL(ctx context.Context, in *pb.PostURLIn) (out *pb.PostURLOut, errRPC error) {
	inputURL := in.GetInputURL()
	owner := ctx.Value(ownerCtxKey)

	result, added, err := utils.Shorten(ctx, s.Database, inputURL, owner.(string), s.Config)
	if err != nil {
		errRPC = status.Errorf(codes.Internal, "while handling GRPC call: %v", err)
		return
	}

	out = &pb.PostURLOut{ResultURL: result}
	if added {
		errRPC = status.Errorf(codes.OK, "")
	} else {
		errRPC = status.Errorf(codes.AlreadyExists, "")
	}
	return
}
