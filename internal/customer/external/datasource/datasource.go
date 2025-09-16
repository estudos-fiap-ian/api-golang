package datasource

import (
	"context"

	"github.com/projeto-estudos/api-golang/internal/customer/dto"
)

type DataSource interface {
	Create(ctx context.Context, customer dto.CustomerDAO) (dto.CustomerDAO, error)
	FindByCPF(ctx context.Context, cpf string) (dto.CustomerDAO, error)
}
