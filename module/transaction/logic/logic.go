package logic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iannrafisyah/gokomodo/enum"
	"github.com/iannrafisyah/gokomodo/model"
	productDto "github.com/iannrafisyah/gokomodo/module/product/dto"
	productLogic "github.com/iannrafisyah/gokomodo/module/product/logic"
	"github.com/iannrafisyah/gokomodo/module/transaction/dto"
	"github.com/iannrafisyah/gokomodo/module/transaction/repository"
	userDto "github.com/iannrafisyah/gokomodo/module/user/dto"
	userLogic "github.com/iannrafisyah/gokomodo/module/user/logic"
	"github.com/iannrafisyah/gokomodo/package/logger"
	"github.com/iannrafisyah/gokomodo/static"
	"github.com/iannrafisyah/gokomodo/utilities"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

// TransactionLogic
type ITransactionLogic interface {
	CreateOrder(context.Context, *dto.CreateOrderRequest, *gorm.DB) error
	FindAll(context.Context, *dto.FindAllRequest) ([]*model.Transactions, error)
	AcceptOrder(context.Context, *dto.AcceptOrderRequest, *gorm.DB) error
}

type TransactionLogic struct {
	fx.In
	Logger          *logger.LogRus
	ProductLogic    productLogic.IProductLogic
	UserLogic       userLogic.IUserLogic
	TransactionRepo repository.ITransactionRepository
}

// NewLogic :
func NewLogic(transactionLogic TransactionLogic) ITransactionLogic {
	return &transactionLogic
}

// Create
func (l *TransactionLogic) CreateOrder(ctx context.Context, reqData *dto.CreateOrderRequest, tx *gorm.DB) error {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var (
		grandTotal   float64
		snapshotItem model.ItemsTransaction
	)

	// Find detail seller
	sellerDetail, err := l.UserLogic.Find(ctx, &userDto.FindRequest{
		ID: reqData.SellerID,
	})
	if err != nil {
		l.Logger.Error(err)
		return err
	}

	// Find detail buyer
	buyerDetail, err := l.UserLogic.Find(ctx, &userDto.FindRequest{
		ID: reqData.BuyerID,
	})
	if err != nil {
		l.Logger.Error(err)
		return err
	}

	// Validate product
	for _, productID := range reqData.Items {
		productDetail, err := l.ProductLogic.Find(ctx, &productDto.FindRequest{
			ID:       productID,
			SellerID: sellerDetail.ID,
		})
		if err != nil {
			l.Logger.Error(err)
			return err
		}

		snapshotItem = append(snapshotItem, model.ProductTransaction{
			ID:          productDetail.ID,
			Name:        productDetail.Name,
			Description: productDetail.Description,
			Price:       productDetail.Price,
		})

		grandTotal += productDetail.Price
	}

	if _, err := l.TransactionRepo.Create(ctx, &model.Transactions{
		BuyerID:     buyerDetail.ID,
		SellerID:    sellerDetail.ID,
		GrandTotal:  grandTotal,
		Status:      enum.TransactionStatusTypePending,
		Items:       snapshotItem,
		Origin:      sellerDetail.Address,
		Destination: buyerDetail.Address,
	}, tx); err != nil {
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}

// FindAll
func (l *TransactionLogic) FindAll(ctx context.Context, reqData *dto.FindAllRequest) ([]*model.Transactions, error) {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return nil, utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	var whereData = model.Transactions{}

	if reqData.RoleID == enum.RoleTypeBuyer {
		whereData = model.Transactions{
			BuyerID: reqData.UserID,
		}
	} else {
		whereData = model.Transactions{
			SellerID: reqData.UserID,
		}
	}

	transactions, err := l.TransactionRepo.FindAll(ctx, &whereData)
	if err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return nil, utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	for _, v := range transactions {
		v.StatusTransaction = v.Status.String()
	}

	return transactions, nil
}

// AcceptOrder
func (l *TransactionLogic) AcceptOrder(ctx context.Context, reqData *dto.AcceptOrderRequest, tx *gorm.DB) error {
	// Validate request data
	if err := reqData.Validate(); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusBadRequest)
	}

	// Find transaction
	if _, err := l.TransactionRepo.Find(ctx, &model.Transactions{
		ID:       reqData.TransactionID,
		SellerID: reqData.SellerID,
	}); err != nil {
		l.Logger.Error(err)
		if err == gorm.ErrRecordNotFound {
			return utilities.ErrorRequest(fmt.Errorf(static.DataNotFound, "transaksi"), http.StatusNotFound)
		}
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	if err := l.TransactionRepo.Update(ctx, &model.Transactions{
		ID:       reqData.TransactionID,
		SellerID: reqData.SellerID,
		Status:   enum.TransactionStatusTypeAccept,
	}, tx); err != nil {
		l.Logger.Error(err)
		return utilities.ErrorRequest(err, http.StatusInternalServerError)
	}

	return nil
}
