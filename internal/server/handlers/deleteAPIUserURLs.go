package handlers

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "shortener/internal/server/pb"

	"github.com/gin-gonic/gin"
)

// DeleteURLs
// Make a url unavailable. Can only delete owned urls.
// Owner is being calculated via cookie of the requester
func (s *ServerHTTP) DeleteURLs(ctx *gin.Context) {
	owner, ok := ctx.Get(string(ownerCtxKey))
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}
	var urls []string
	if err := ctx.BindJSON(&urls); err != nil {
		ctx.Status(http.StatusInternalServerError)
	}

	ctx.Status(http.StatusAccepted)
	go func() {
		err := s.Database.Delete(ctx, urls, owner.(string))
		if err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}()
}

// DeleteURLs
// Make a url unavailable. Can only delete owned urls.
// Owner is being calculated via cookie of the requester
func (s *ServerGRPC) DeleteURLs(ctx context.Context, in *pb.DeleteURLIn) (out *pb.Empty, errRPC error) {
	owner := ctx.Value(ownerCtxKey)
	urls := in.GetInputURLs()

	err := s.Database.Delete(ctx, urls, owner.(string))
	if err != nil {
		errRPC = status.Errorf(codes.Internal, "")
		return
	}
	return
}
