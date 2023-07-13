package handlers

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shortener/internal/cfg"
	pb "shortener/internal/server/pb"
	"shortener/internal/storage"

	"github.com/gin-gonic/gin"

	"shortener/internal/server/handlers/utils"
)

type reqElement struct {
	LineID string `json:"correlation_id"`
	URL    string `json:"original_url"`
}
type resElement struct {
	LineID string `json:"correlation_id"`
	URL    string `json:"short_url"`
}

// BatchShorten
// Handler for batch shortening of urs.
// Accepts array of JSON objects containing:
// {"correlation_id": "id of url", "original_url": "url to be shortened"}
// returns array of following jsons:
// {"correlation_id": "id of url", "short_url": "url that was shortened"}
func (s *ServerHTTP) BatchShorten(ctx *gin.Context) {
	var newReqBody []*reqElement
	owner, ok := ctx.Get(string(ownerCtxKey))
	if !ok {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if err := ctx.BindJSON(&newReqBody); err != nil {
		return
	}

	newResBody, err := batchShorten(ctx, s.Database, s.Config, newReqBody, owner.(string))

	if errors.Is(err, errInvalidInput) {
		ctx.Status(http.StatusBadRequest)
	} else {
		ctx.IndentedJSON(http.StatusCreated, newResBody)
	}
}

// BatchShorten
// Handler for batch shortening of urs.
// Accepts array of JSON objects containing:
// {"correlation_id": "id of url", "original_url": "url to be shortened"}
// returns array of following jsons:
// {"correlation_id": "id of url", "short_url": "url that was shortened"}
func (s *ServerGRPC) BatchShorten(ctx context.Context, in *pb.BatchShortenIn) (out *pb.BatchShortenOut, errRPC error) {
	owner := ctx.Value(ownerCtxKey)
	reqBody := in.GetData()
	newReqBody := make([]*reqElement, 0, len(reqBody))
	for _, el := range reqBody {
		newReqBody = append(newReqBody, &reqElement{
			LineID: el.GetCorrelationID(),
			URL:    el.GetInputURL(),
		})
	}

	data, err := batchShorten(ctx, s.Database, s.Config, newReqBody, owner.(string))
	if errors.Is(err, errInvalidInput) {
		errRPC = status.Errorf(codes.InvalidArgument, "")
		return
	}

	newRes := make([]*pb.BatchShortenOutElement, 0, len(data))
	for _, el := range data {
		newRes = append(newRes, &pb.BatchShortenOutElement{
			CorrelationID: el.LineID,
			ResultURL:     el.LineID,
		})
	}
	out = &pb.BatchShortenOut{Data: newRes}
	return
}

func batchShorten(
	ctx context.Context,
	database storage.DatabaseORM,
	config *cfg.ConfigT,
	newReqBody []*reqElement,
	owner string,
) ([]*resElement, error) {
	newResBody := make([]*resElement, 0, len(newReqBody))
	for _, element := range newReqBody {
		result, _, err := utils.Shorten(ctx, database, element.URL, owner, config)
		if err != nil {
			return nil, errInvalidInput
		}
		newResBody = append(newResBody, &resElement{LineID: element.LineID, URL: result})
	}
	return newResBody, nil
}
