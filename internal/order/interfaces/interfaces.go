package interfaces

import (
	"context"

	paymententity "github.com/projeto-estudos/api-golang/internal/payment/entity"
	productentity "github.com/projeto-estudos/api-golang/internal/product/entity"
	productorderentity "github.com/projeto-estudos/api-golang/internal/productorder/entity"
)

type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]productentity.Product, error)
}

type ProductOrderService interface {
	CreateBulk(ctx context.Context, productOrders []productorderentity.ProductOrder) (int, error)
}

type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) (*paymententity.Payment, error)
}
