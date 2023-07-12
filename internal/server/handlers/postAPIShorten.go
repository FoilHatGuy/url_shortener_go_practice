package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	utils "shortener/internal/server/handlers/utils"
	pb "shortener/internal/server/pb"
)

// PostAPIURL
// Handler for batch shortening of urs.
// Takes the field "url" from request body and returns the result in "url" field for accessing original url.
func (s *ServerHTTP) PostAPIURL(ctx *gin.Context) {
	var newReqBody struct {
		URL string `json:"url"`
	}
	owner, ok := ctx.Get(string(ownerCtxKey))
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := ctx.BindJSON(&newReqBody); err != nil {
		return
	}

	result, added, err := utils.Shorten(ctx, s.Database, newReqBody.URL, owner.(string), s.Config)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	newResBody := struct {
		Result string `json:"result"`
	}{result}
	if added {
		ctx.IndentedJSON(http.StatusCreated, newResBody)
	} else {
		ctx.IndentedJSON(http.StatusConflict, newResBody)
	}
}

// PostAPIURL
// Handler for batch shortening of urs.
// Takes the field "url" from request body and returns the result in "url" field for accessing original url.
// NOW IDENTICAL TO PostURL
func (s *ServerGRPC) PostAPIURL(ctx context.Context, in *pb.PostURLIn) (out *pb.PostURLOut, errRPC error) {
	inputURL := in.GetInputURL()
	owner := ctx.Value(ownerCtxKey)

	result, added, err := utils.Shorten(ctx, s.Database, inputURL, owner.(string), s.Config)
	if err != nil {
		errRPC = status.Errorf(codes.Internal, "while handling GRPC call: %v", err)
		return
	}

	out.ResultURL = result
	if added {
		errRPC = status.Errorf(codes.OK, "")
	} else {
		errRPC = status.Errorf(codes.AlreadyExists, "")
	}
	return
}
