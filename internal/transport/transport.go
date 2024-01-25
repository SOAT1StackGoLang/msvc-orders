package transport

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func NewHTTPHandler(svc service.CategoriesService, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	catEndpoints := endpoint.MakeCategoryEndpoints(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).Path("/categories/{id}").Handler(httptransport.NewServer(
		catEndpoints.GetCategoryEndpoint,
		decodeGetCategoriesRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetCategoriesRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.GetCategoryRequest{ID: id}, nil
}
