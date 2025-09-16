package controllers

import (
	"context"

	"github.com/projeto-estudos/api-golang/internal/payment/dto"
	"github.com/projeto-estudos/api-golang/internal/payment/presenter"
	"github.com/projeto-estudos/api-golang/internal/payment/usecases"
)

type Controller struct {
	paymentUseCase *usecases.UseCases
}

func Build(paymentUseCase *usecases.UseCases) *Controller {
	return &Controller{
		paymentUseCase: paymentUseCase,
	}
}

func (c *Controller) CreateByOrderID(ctx context.Context, orderID string) (dto.PaymentResponseDTO, error) {
	presenter := presenter.Build()

	payment, err := c.paymentUseCase.CreateByOrderID(ctx, orderID)
	if err != nil {
		return dto.PaymentResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(payment), nil
}

func (c *Controller) CheckPayment(ctx context.Context, requestUrl string) (interface{}, error) {
	return c.paymentUseCase.CheckPayment(ctx, requestUrl)
}
