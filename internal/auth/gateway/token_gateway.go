package gateway

import (
	"github.com/projeto-estudos/api-golang/internal/auth/entity"
)

type TokenGateway interface {
	GenerateToken(userID, userType string, additionalClaims map[string]any) (string, error)
	ValidateToken(tokenString string) (*entity.CustomClaims, error)
}
