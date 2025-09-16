package presenter

import (
	"github.com/projeto-estudos/api-golang/internal/payment/dto"
	"github.com/projeto-estudos/api-golang/internal/payment/entity"
)

type Presenter struct {
}

func Build() *Presenter {
	return &Presenter{}
}

func (p *Presenter) FromEntityToResponseDTO(payment entity.Payment) dto.PaymentResponseDTO {
	return dto.PaymentResponseDTO{
		ID:      payment.ID,
		OrderID: payment.OrderID,
		QrCode:  payment.QrCode,
		Status:  payment.Status,
	}
}
