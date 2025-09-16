package usecase

import (
	"github.com/projeto-estudos/api-golang/internal/auth/entity"
	"github.com/projeto-estudos/api-golang/internal/auth/gateway"
)

type ValidateTokenUseCase struct {
	tokenGateway gateway.TokenGateway
}

func NewValidateTokenUseCase(tokenGateway gateway.TokenGateway) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		tokenGateway: tokenGateway,
	}
}

func (uc *ValidateTokenUseCase) Execute(tokenString string) (*entity.CustomClaims, error) {
	return uc.tokenGateway.ValidateToken(tokenString)
}
