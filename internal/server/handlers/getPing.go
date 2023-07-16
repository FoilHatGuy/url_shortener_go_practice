package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "shortener/internal/server/pb"
)

// PingDatabase
// Ping server+database activity
func (s *ServerHTTP) PingDatabase(ctx *gin.Context) {
	ping := s.Database.Ping(ctx)
	if ping {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}

// PingDatabase
// Ping server+database activity
func (s *ServerGRPC) PingDatabase(ctx context.Context, _ *pb.Empty) (out *pb.Empty, errRPC error) {
	ping := s.Database.Ping(ctx)
	if !ping {
		errRPC = status.Errorf(codes.Internal, "")
	}
	return
}
