package dto

import (
	"time"

	"github.com/projeto-estudos/api-golang/internal/admin/entity"
	gormEntity "github.com/projeto-estudos/api-golang/internal/shared/entity"
	"github.com/google/uuid"
)

type AdminRequestDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminDAO struct {
	gormEntity.Entity
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

func ToAdminDAO(admin entity.Admin) AdminDAO {
	return AdminDAO{
		Entity: gormEntity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    admin.Email,
		Password: admin.Password,
	}
}

func FromAdminDAO(dao AdminDAO) entity.Admin {
	return entity.Admin{
		Id:       dao.ID,
		Email:    dao.Email,
		Password: dao.Password,
	}
}

func FromAdminRequestDTO(dto AdminRequestDTO) entity.Admin {
	return entity.Admin{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
