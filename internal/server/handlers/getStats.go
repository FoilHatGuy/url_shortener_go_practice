package handlers

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	pb "shortener/internal/server/pb"
)

// GetStats
// Ping server+database activity
func (s *ServerHTTP) GetStats(ctx *gin.Context) {
	k := ctx.GetHeader("X-Real-IP")
	if k == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ip := net.ParseIP(k)
	_, cidr, err := net.ParseCIDR(s.Config.Server.TrustedSubnet)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	valid := cidr.Contains(ip)
	if !valid {
		ctx.Status(http.StatusForbidden)
		return
	}

	data, err := s.Database.GetStats(ctx)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetStats
// Ping server+database activity
func (s *ServerGRPC) GetStats(ctx context.Context, _ *pb.Empty) (out *pb.Stats, errRPC error) {
	p, _ := peer.FromContext(ctx)
	k := p.Addr.String()
	if k == "" {
		errRPC = status.Errorf(codes.InvalidArgument, "")
		return
	}
	ip := net.ParseIP(k)
	_, cidr, err := net.ParseCIDR(s.Config.Server.TrustedSubnet)
	if err != nil {
		errRPC = status.Errorf(codes.InvalidArgument, "")
		return
	}
	valid := cidr.Contains(ip)
	if !valid {
		errRPC = status.Errorf(codes.PermissionDenied, "")
		return
	}

	data, err := s.Database.GetStats(ctx)
	if err != nil {
		errRPC = status.Errorf(codes.InvalidArgument, "")
		return
	}
	out = &pb.Stats{
		URLsCount:  data.URLs,
		UsersCount: data.Users,
	}
	return
}
