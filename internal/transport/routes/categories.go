package routes

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func NewCategoriesRouter(svc service.CategoriesService, r *mux.Router, logger kitlog.Logger) *mux.Router {
	catEndpoints := endpoint.MakeCategoryEndpoints(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).Path("/category/all").Queries("limit", "{limit:[0-9]+}", "offset", "{offset:[0-9]+}").Handler(httptransport.NewServer(
		catEndpoints.ListCategoriesEndpoint,
		decodeListCategoriesRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodGet).Path("/category/{id}").Handler(httptransport.NewServer(
		catEndpoints.GetCategoryEndpoint,
		decodeGetCategoriesRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodPost).Path("/category").Handler(httptransport.NewServer(
		catEndpoints.InsertCategoryEndpoint,
		decodeInsertCategoriesRequest,
		encodeResponse,
		options...,
	))

	r.Methods(http.MethodDelete).Path("/category").Handler(httptransport.NewServer(
		catEndpoints.DeleteCategoryEndpoint,
		decodeDeleteCategoriesRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeListCategoriesRequest(_ context.Context, r *http.Request) (request any, err error) {
	query := r.URL.Query()
	limit := query.Get("limit")

	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, err
	}
	if limitInt == 0 {
		return nil, ErrBadRequest
	}

	offset := query.Get("offset")

	offsetInt, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		return nil, err
	}

	return endpoint.ListCategoriesRequest{
		Limit:  limitInt,
		Offset: offsetInt,
	}, nil

}

func decodeDeleteCategoriesRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.DeleteCategoryRequest{ID: id}, nil
}

func decodeInsertCategoriesRequest(_ context.Context, r *http.Request) (request any, err error) {
	var req endpoint.InsertCategoryRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest
	}

	return endpoint.InsertCategoryRequest{
		Name: req.Name,
	}, nil
}

func decodeGetCategoriesRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.GetCategoryRequest{ID: id}, nil
}
