package handlers

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "shortener/internal/server/pb"

	"github.com/gin-gonic/gin"
)

// GetAllOwnedURL
// Get all owned urls.
// Owner is being calculated via cookie of the requester
func (s *ServerHTTP) GetAllOwnedURL(ctx *gin.Context) {
	owner, ok := ctx.Get(string(ownerCtxKey))
	if !ok {
		fmt.Println("NO OWNER CONTEXT")
		ctx.Status(http.StatusBadRequest)
		return
	}
	result, err := s.Database.GetURLByOwner(ctx, owner.(string))
	if err != nil {
		fmt.Println("ERROR WHILE GETTING DATA FROM DB")
		ctx.Status(http.StatusBadRequest)
		return
	}
	if result != nil {
		ctx.IndentedJSON(http.StatusOK, result)
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

// GetAllOwnedURL
// Get all owned urls.
// Owner is being calculated via cookie of the requester
func (s *ServerGRPC) GetAllOwnedURL(ctx context.Context, _ *pb.Empty) (out *pb.OwnedURLsOut, errRPC error) {
	owner := ctx.Value(ownerCtxKey)
	result, err := s.Database.GetURLByOwner(ctx, owner.(string))
	if err != nil {
		errRPC = status.Errorf(codes.Internal, "")
		return
	}
	if result == nil {
		errRPC = status.Errorf(codes.Internal, "")
		return
	}
	response := make([]*pb.URLPair, 0, len(result))
	for _, el := range result {
		response = append(response, &pb.URLPair{
			ShortURL:    el.ShortURL,
			OriginalURL: el.OriginalURL,
		})
	}
	out = &pb.OwnedURLsOut{Data: response}
	return
}
