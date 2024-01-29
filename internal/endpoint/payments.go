package endpoint

import (
	"context"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/helpers"
	"github.com/SOAT1StackGoLang/msvc-orders/internal/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	PaymentsEndpoints struct {
		GetPaymentEndpoint endpoint.Endpoint
	}
)

func MakePaymentsEndpoint(svc service.PaymentsService) PaymentsEndpoints {
	return PaymentsEndpoints{
		GetPaymentEndpoint: makeGetPayment(svc),
	}
}

func makeGetPayment(svc service.PaymentsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetPaymentRequest)

		uid, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, err
		}

		payment, err := svc.GetPayment(ctx, uid)
		if err != nil {
			return nil, err
		}

		out := GetPaymentResponse{
			ID:        payment.ID.String(),
			CreatedAt: payment.CreatedAt.String(),
			UpdatedAt: payment.UpdatedAt.String(),
			Price:     helpers.ParseDecimalToString(payment.Price),
			OrderID:   payment.OrderID.String(),
			Status:    string(payment.Status),
		}

		return out, nil
	}
}

