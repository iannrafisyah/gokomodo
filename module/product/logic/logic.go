package logic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iannrafisyah/gokomodo/model"
	"github.com/iannrafisyah/gokomodo/module/product/dto"
	"github.com/iannrafisyah/gokomodo/module/product/repository"
	"github.com/iannrafisyah/gokomodo/package/logger"
	"github.com/iannrafisyah/gokomodo/static"
	"github.com/iannrafisyah/gokomodo/utilities"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProductLogic
type IProductLogic interface {
	Create(context.Context, *dto.CreateRequest, *gorm.DB) error
	FindAll(context.Context, *dto.FindAllRequest) ([]*model.Products, error)
	Find(context.Context, *dto.FindRequest) (*model.Products, error)
}

type ProductLogic struct {
	fx.In
	Logger      *logger.LogRus
	ProductRepo repository.ISellerRepository
}

// NewLogic :
func NewLogic(productLogic ProductLogic) IProductLogic {
	return &productLogic
}

// Create
func (l *ProductLogic) Create(ctx context.Context, reqData *dto.CreateRequest, tx *gorm.DB) error {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	if _, err := l.ProductRepo.Create(ctx, &model.Products{
		SellerID:    reqData.SellerID,
		Name:        reqData.Name,
		Description: reqData.Description,
		Price:       reqData.Price,
	}, tx); err != nil {
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}

// FindAll
func (l *ProductLogic) FindAll(ctx context.Context, reqData *dto.FindAllRequest) ([]*model.Products, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	products, err := l.ProductRepo.FindAll(ctx, &model.Products{
		SellerID: reqData.SellerID,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "produk"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return products, nil
}

// FindByID
func (l *ProductLogic) Find(ctx context.Context, reqData *dto.FindRequest) (*model.Products, error) {
	product, err := l.ProductRepo.Find(ctx, &model.Products{
		ID:       reqData.ID,
		SellerID: reqData.SellerID,
	})
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "product"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}
	return product, nil
}
