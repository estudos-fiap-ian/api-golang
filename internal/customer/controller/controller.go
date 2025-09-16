package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/projeto-estudos/api-golang/internal/customer/dto"
	"github.com/projeto-estudos/api-golang/internal/customer/external/datasource"
	"github.com/projeto-estudos/api-golang/internal/customer/gateway"
	"github.com/projeto-estudos/api-golang/internal/customer/usecases"
	apperror "github.com/projeto-estudos/api-golang/internal/shared/errors"
)

type Controller struct {
	CustomerDataSource datasource.DataSource
	AuthGateway        gateway.AuthGateway
}

func Build(customerDataSource datasource.DataSource, authGateway gateway.AuthGateway) *Controller {
	return &Controller{
		CustomerDataSource: customerDataSource,
		AuthGateway:        authGateway,
	}
}

func (c *Controller) Create(ctx context.Context, customerRequest dto.CustomerRequestDTO) (string, error) {
	customerGateway := gateway.Build(c.CustomerDataSource)
	useCase := usecases.Build(*customerGateway)
	customer := dto.FromCustomerRequestDTO(customerRequest)
	customerId, err := useCase.Create(ctx, customer)

	if err != nil {
		return "", err
	}

	return customerId, nil

}

func (c *Controller) Identify(ctx context.Context, cpf string) (string, error) {

	if cpf == "" {
		return c.createAnonymousToken()
	}

	customerGateway := gateway.Build(c.CustomerDataSource)
	useCase := usecases.Build(*customerGateway)
	customerId, err := useCase.FindByCPF(ctx, cpf)

	if err != nil {
		return "", err
	}

	token, err := c.createToken(customerId, false)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (c *Controller) createAnonymousToken() (string, error) {
	anonymousID := uuid.NewString()

	token, err := c.createToken(anonymousID, true)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *Controller) createToken(id string, isAnonymous bool) (string, error) {
	additionalClaims := map[string]any{
		"is_anonymous": isAnonymous,
	}

	token, err := c.AuthGateway.GenerateToken(id, "customer", additionalClaims)
	if err != nil {
		return "", &apperror.InternalError{Msg: "Error creating token"}
	}

	return token, nil
}
