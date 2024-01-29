package endpoint

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service/models"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	OrdersEndpoint struct {
		GetOrderEndpoint         endpoint.Endpoint
		CreateOrderEndpoint      endpoint.Endpoint
		UpdateOrderItemsEndpoint endpoint.Endpoint
		DeleteOrderEndpoint      endpoint.Endpoint
		CheckoutEndpoint         endpoint.Endpoint
		GetOrderByPaymentID      endpoint.Endpoint
	}
)

func MakeOrdersEndpoint(svc service.OrdersService) OrdersEndpoint {
	return OrdersEndpoint{
		GetOrderEndpoint:    makeGetOrderEndpoint(svc),
		CreateOrderEndpoint: makeCreateOrderEndpoint(svc),
	}
}

func makeCreateOrderEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateOrderRequest)

		ids := req.ProductsIDs
		var prods []models.Product

		for _, v := range ids {
			prodID, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}
			prods = append(prods, models.Product{ID: prodID})
		}

		order, err := svc.CreateOrder(ctx, prods)
		if err != nil {
			return nil, err
		}

		return OrderResponseFromModel(order), nil
	}
}

func makeGetOrderEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOrderRequest)

		uid, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		order, err := svc.GetOrder(ctx, uid)
		if err != nil {
			return nil, err
		}

		out := OrderResponseFromModel(order)
		return out, nil
	}
}
