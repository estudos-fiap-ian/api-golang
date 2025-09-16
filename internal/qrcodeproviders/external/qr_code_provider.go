package external

import (
	"context"
	"github.com/projeto-estudos/api-golang/internal/qrcodeproviders/dtos"
	"github.com/projeto-estudos/api-golang/internal/qrcodeproviders/entities"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request entities.GenerateQRCodeParams) (string, error)
	CheckPayment(ctx context.Context, requestUrl string) (dtos.ResponseVerifyOrderDTO, error)
}
