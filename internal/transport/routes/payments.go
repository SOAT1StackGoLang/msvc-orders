package routes

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/endpoint"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	kittransport "github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"net/http"
)

func NewPaymentsRouter(svc service.PaymentsService, r *mux.Router, logger kitlog.Logger) *mux.Router {
	prodEndpoints := endpoint.MakePaymentsEndpoint(svc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(kittransport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods(http.MethodGet).Path("/payment/{id}").Handler(httptransport.NewServer(prodEndpoints.GetPaymentEndpoint,
		decodeGetPaymentsRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeGetPaymentsRequest(_ context.Context, r *http.Request) (request any, err error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return endpoint.GetPaymentRequest{ID: id}, nil
}
