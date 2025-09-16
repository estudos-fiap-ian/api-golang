package gateway

import (
	"context"

	"github.com/projeto-estudos/api-golang/internal/admin/dto"
	"github.com/projeto-estudos/api-golang/internal/admin/entity"
	"github.com/projeto-estudos/api-golang/internal/admin/external/datasource"
	apperror "github.com/projeto-estudos/api-golang/internal/shared/errors"
)

type Gateway struct {
	Datasource datasource.DataSource
}

func Build(datasource datasource.DataSource) *Gateway {
	return &Gateway{
		Datasource: datasource,
	}
}

func (g *Gateway) Create(c context.Context, admin entity.Admin) error {
	var adminDAO = dto.ToAdminDAO(admin)
	err := g.Datasource.Create(c, adminDAO)

	if err != nil {
		return &apperror.InternalError{Msg: err.Error()}
	}

	return nil
}

func (g *Gateway) FindByEmail(c context.Context, email string) (entity.Admin, error) {
	adminDAO, err := g.Datasource.FindByEmail(c, email)

	if err != nil {
		return entity.Admin{}, &apperror.InternalError{Msg: err.Error()}
	}

	admin := dto.FromAdminDAO(adminDAO)

	return admin, nil
}
