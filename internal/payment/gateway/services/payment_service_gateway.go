package services

import (
	"context"

	"github.com/projeto-estudos/api-golang/internal/order/interfaces"
	"github.com/projeto-estudos/api-golang/internal/payment/entity"
	"github.com/projeto-estudos/api-golang/internal/payment/usecases"
)

type PaymentServiceGateway struct {
	paymentUseCase *usecases.UseCases
}

func NewPaymentServiceGateway(paymentUseCase *usecases.UseCases) interfaces.PaymentService {
	return &PaymentServiceGateway{
		paymentUseCase: paymentUseCase,
	}
}

func (a *PaymentServiceGateway) CreateByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	payment, err := a.paymentUseCase.CreateByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
