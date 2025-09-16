package controller

import (
	"github.com/projeto-estudos/api-golang/internal/auth/entity"
	"github.com/projeto-estudos/api-golang/internal/auth/gateway"
	"github.com/projeto-estudos/api-golang/internal/auth/usecase"
)

type Controller struct {
	generateTokenUC *usecase.GenerateTokenUseCase
	validateTokenUC *usecase.ValidateTokenUseCase
}

func New(tokenGateway gateway.TokenGateway) *Controller {
	return &Controller{
		generateTokenUC: usecase.NewGenerateTokenUseCase(tokenGateway),
		validateTokenUC: usecase.NewValidateTokenUseCase(tokenGateway),
	}
}

func (c *Controller) GenerateToken(userID, userType string, additionalClaims map[string]any) (string, error) {
	return c.generateTokenUC.Execute(userID, userType, additionalClaims)
}

func (c *Controller) ValidateToken(tokenString string) (*entity.CustomClaims, error) {
	return c.validateTokenUC.Execute(tokenString)
}
