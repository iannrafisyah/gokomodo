package repository

import (
	"context"
	"time"

	"github.com/iannrafisyah/gokomodo/database/postgres"
	"github.com/iannrafisyah/gokomodo/model"
	"github.com/iannrafisyah/gokomodo/package/logger"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TransactionRepository
type ITransactionRepository interface {
	Create(context.Context, *model.Transactions, *gorm.DB) (*int, error)
	FindAll(context.Context, *model.Transactions) ([]*model.Transactions, error)
	Find(context.Context, *model.Transactions) (*model.Transactions, error)
	Update(context.Context, *model.Transactions, *gorm.DB) error
}

type TransactionRepository struct {
	fx.In
	Logger   *logger.LogRus
	Database *postgres.DB
}

// NewRepository :
func NewRepository(transactionRepository TransactionRepository) ITransactionRepository {
	return &transactionRepository
}

// Create
func (l *TransactionRepository) Create(ctx context.Context, reqData *model.Transactions, tx *gorm.DB) (*int, error) {
	if err := tx.WithContext(ctx).Create(&reqData).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return &reqData.ID, nil
}

// FindAll
func (l *TransactionRepository) FindAll(ctx context.Context, reqData *model.Transactions) ([]*model.Transactions, error) {
	transactions := []*model.Transactions{}

	if err := l.Database.Gorm.WithContext(ctx).Model(&model.Transactions{}).
		Preload("Seller").
		Preload("Buyer").
		Where(&model.Transactions{
			SellerID: reqData.SellerID,
			BuyerID:  reqData.BuyerID,
		}).
		Order("id desc").
		Find(&transactions).
		Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}

	return transactions, nil
}

// Find
func (l *TransactionRepository) Find(ctx context.Context, reqData *model.Transactions) (*model.Transactions, error) {
	transaction := new(model.Transactions)
	if err := l.Database.Gorm.WithContext(ctx).
		Where(&model.Transactions{
			ID:       reqData.ID,
			SellerID: reqData.SellerID,
		}).
		First(&transaction).Error; err != nil {
		l.Logger.Error(err)
		return nil, err
	}
	return transaction, nil
}

// Update
func (l *TransactionRepository) Update(ctx context.Context, reqData *model.Transactions, tx *gorm.DB) error {
	if err := tx.WithContext(ctx).Model(&model.Transactions{}).
		Where("id = ?", reqData.ID).
		Where("seller_id = ?", reqData.SellerID).
		Updates(model.Transactions{
			Status:    reqData.Status,
			UpdatedAt: time.Now(),
		}).Error; err != nil {
		l.Logger.Error(err)
		return err
	}
	return nil
}
