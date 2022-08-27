package repository

import (
	"context"

	"github.com/iannrafisyah/gokomodo/database/postgres"
	"github.com/iannrafisyah/gokomodo/model"
	"github.com/iannrafisyah/gokomodo/package/logger"

	"go.uber.org/fx"
)

// UserRepository
type IUserRepository interface {
	Find(context.Context, *model.Users) (*model.Users, error)
}

type UserRepository struct {
	fx.In
	Logger   *logger.LogRus
	Database *postgres.DB
}

// NewRepository :
func NewRepository(userRepository UserRepository) IUserRepository {
	return &userRepository
}

// Find
func (l *UserRepository) Find(ctx context.Context, reqData *model.Users) (*model.Users, error) {
	product := new(model.Users)

	if err := l.Database.Gorm.WithContext(ctx).
		Where(&model.Users{
			ID:    reqData.ID,
			Email: reqData.Email,
		}).First(&product).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return product, nil
}
