package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"shortener/internal/auth"
)

// Cooker manages sid cookies.
// Deciphers "user" cookie using auth.AuthEngine Validate to get key from sid.
// If the referrer doesn't possess the cookie, it generates a new sid and sets user's cookie.
// Either way, it allows the request handling.
func (s *ServerHTTP) Cooker(ctx *gin.Context) {
	cookie, err := ctx.Cookie("user")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	newCookie, key, err := cooker(cookie, s.Security)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	ctx.SetCookie("user", newCookie, s.Config.Server.CookieLifetime, "/",
		strings.Split(s.Config.Server.AddressHTTP, ":")[0], false, true)
	ctx.Set(string(ownerCtxKey), key)
	ctx.Next()
}

// Cooker manages sid cookies.
// Deciphers "user" cookie using auth.AuthEngine Validate to get key from sid.
// If the referrer doesn't possess the cookie, it generates a new sid and sets user's cookie.
// Either way, it allows the request handling.
func (s *ServerGRPC) Cooker(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (response interface{}, errRPC error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	cookie := md.Get("user")
	if len(cookie) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "missing 'user' field in metadata")
	}

	newCookie, key, err := cooker(cookie[0], s.Security)
	if err != nil {
		errRPC = status.Errorf(codes.Unauthenticated, "invalid token")
		return
	}
	ctx = context.WithValue(ctx, ownerCtxKey, key)
	metadata.AppendToOutgoingContext(ctx, "user", newCookie)

	response, errRPC = handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error: %v", err)
	}
	return
}

func cooker(cookie string, validator *auth.EngineT) (newCookie, key string, err error) {
	if cookie != "" {
		key, err = validator.Validate(cookie) // validate input cookie/UID
		if err == nil {                       // if ok => return same cookie
			newCookie = cookie
			return
		}
	}
	newCookie, key, err = validator.Generate() // else generate new
	if err != nil {                            // check for error during generation
		return "", "", fmt.Errorf("inside cookie checker: %w", err)
	}
	return // finally return new values
}
