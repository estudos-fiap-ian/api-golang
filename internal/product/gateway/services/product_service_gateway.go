package services

import (
	"context"

	"github.com/projeto-estudos/api-golang/internal/order/interfaces"
	"github.com/projeto-estudos/api-golang/internal/product/entity"
	"github.com/projeto-estudos/api-golang/internal/product/usecases"
)

type ProductServiceGateway struct {
	productUseCase *usecases.UseCases
}

func NewProductServiceGateway(productUseCase *usecases.UseCases) interfaces.ProductService {
	return &ProductServiceGateway{
		productUseCase: productUseCase,
	}
}

func (a *ProductServiceGateway) FindByIDs(ctx context.Context, productIDs []string) ([]entity.Product, error) {
	return a.productUseCase.FindByIDs(ctx, productIDs)
}
