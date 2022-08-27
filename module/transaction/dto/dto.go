package dto

import (
	"errors"
	"fmt"

	"github.com/iannrafisyah/gokomodo/enum"
	"github.com/iannrafisyah/gokomodo/static"
)

type CreateOrderRequest struct {
	BuyerID  int
	SellerID int
	Items    []int
	RoleID   enum.RoleType
}

func (d *CreateOrderRequest) Validate() error {
	if d.BuyerID <= 0 {
		return fmt.Errorf(static.EmptyValue, "BuyerID")
	}
	if d.SellerID <= 0 {
		return fmt.Errorf(static.EmptyValue, "SellerID")
	}
	if len(d.Items) == 0 {
		return fmt.Errorf(static.EmptyValue, "Product")
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	if d.RoleID != enum.RoleTypeBuyer {
		return errors.New(static.Authorization)
	}
	return nil
}

type FindAllRequest struct {
	UserID int
	Items  []int
	RoleID enum.RoleType
}

func (d *FindAllRequest) Validate() error {
	if d.UserID <= 0 {
		return fmt.Errorf(static.EmptyValue, "UserID")
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	return nil
}

type AcceptOrderRequest struct {
	SellerID      int
	TransactionID int
	RoleID        enum.RoleType
}

func (d *AcceptOrderRequest) Validate() error {
	if d.SellerID <= 0 {
		return fmt.Errorf(static.EmptyValue, "SellerID")
	}
	if d.TransactionID <= 0 {
		return fmt.Errorf(static.EmptyValue, "TransactionID")
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	if d.RoleID != enum.RoleTypeSeller {
		return errors.New(static.Authorization)
	}
	return nil
}
