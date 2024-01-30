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
		ListOrdersEndpoint       endpoint.Endpoint
		DeleteOrderEndpoint      endpoint.Endpoint
		OrderCheckoutEndpoint    endpoint.Endpoint
		GetOrderByPaymentID      endpoint.Endpoint
	}
)

func MakeOrdersEndpoint(svc service.OrdersService) OrdersEndpoint {
	return OrdersEndpoint{
		GetOrderEndpoint:         makeGetOrderEndpoint(svc),
		CreateOrderEndpoint:      makeCreateOrderEndpoint(svc),
		UpdateOrderItemsEndpoint: makeUpdateOrderItemsEndpoint(svc),
		DeleteOrderEndpoint:      makeDeleteOrderEndpoint(svc),
		OrderCheckoutEndpoint:    makeOrderCheckoutEndpoint(svc),
		GetOrderByPaymentID:      makeGetOrderByPaymentIDEndpoint(svc),
		ListOrdersEndpoint:       makeListOrdersEndpoint(svc),
	}
}

func makeListOrdersEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ListOrderRequest)

		svcOut, err := svc.ListOrders(ctx, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}

		var orders []OrderResponse
		for _, o := range svcOut.Orders {
			orders = append(orders, OrderResponseFromModel(o))
		}

		return OrderList{
			Orders: orders,
			Limit:  req.Limit,
			Offset: req.Offset,
			Total:  len(orders),
		}, nil
	}
}

func makeGetOrderByPaymentIDEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetOrderByPaymentIDRequest)

		pID, err := uuid.Parse(req.PaymentID)
		if err != nil {
			return nil, err
		}

		order, err := svc.GetOrderByPaymentID(ctx, pID)
		if err != nil {
			return nil, err
		}

		return OrderResponseFromModel(order), nil
	}
}

func makeOrderCheckoutEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CheckoutOrderRequest)

		oID, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		checkout, err := svc.Checkout(ctx, oID)
		if err != nil {
			return nil, err
		}

		return OrderResponseFromModel(checkout), nil
	}
}

func makeDeleteOrderEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DeleteOrderRequest)

		oID, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		out := svc.DeleteOrder(ctx, oID)
		if out != nil {
			return DeleteOrderResponse{Deleted: "false"}, err
		}

		return DeleteOrderResponse{Deleted: "true"}, nil
	}
}

func makeUpdateOrderItemsEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateOrderRequest)

		oID, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		ids := req.ProductsIDs
		var prods []models.Product

		for _, v := range ids {
			prodID, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}
			prods = append(prods, models.Product{ID: prodID})
		}

		items, err := svc.UpdateOrderItems(ctx, oID, prods)
		if err != nil {
			return nil, err
		}

		return OrderResponseFromModel(items), nil

	}
}

func makeCreateOrderEndpoint(svc service.OrdersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateOrderRequest)

		uid, _ := uuid.Parse(req.UserID)

		ids := req.ProductsIDs
		var prods []models.Product

		for _, v := range ids {
			prodID, err := uuid.Parse(v)
			if err != nil {
				return nil, err
			}
			prods = append(prods, models.Product{ID: prodID})
		}

		order, err := svc.CreateOrder(ctx, prods, uid)
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
