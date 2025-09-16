package services

import (
	"context"

	orderinterfaces "github.com/projeto-estudos/api-golang/internal/order/interfaces"
	paymentinterfaces "github.com/projeto-estudos/api-golang/internal/payment/interfaces"
	"github.com/projeto-estudos/api-golang/internal/productorder/entity"
	"github.com/projeto-estudos/api-golang/internal/productorder/usecases"
)

type ProductOrderServiceGateway struct {
	productOrderUseCase *usecases.UseCases
}

func NewProductOrderServiceGateway(productOrderUseCase *usecases.UseCases) (
	orderinterfaces.ProductOrderService,
	paymentinterfaces.ProductOrderService,
) {
	adapter := &ProductOrderServiceGateway{
		productOrderUseCase: productOrderUseCase,
	}
	return adapter, adapter
}

func (a *ProductOrderServiceGateway) CreateBulk(ctx context.Context, productOrders []entity.ProductOrder) (int, error) {
	return a.productOrderUseCase.CreateBulk(ctx, productOrders)
}

func (a *ProductOrderServiceGateway) FindByOrderID(ctx context.Context, orderID string) ([]entity.ProductOrder, error) {
	return a.productOrderUseCase.FindByOrderID(ctx, orderID)
}
