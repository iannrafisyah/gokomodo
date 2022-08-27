package dto

import (
	"fmt"

	"github.com/iannrafisyah/gokomodo/static"
)

type LoginRequest struct {
	Email    string
	Password string
}

func (d *LoginRequest) Validate() error {
	if d.Email == "" {
		return fmt.Errorf(static.EmptyValue, "email")
	}
	if d.Password == "" {
		return fmt.Errorf(static.EmptyValue, "password")
	}
	return nil
}

type Response struct {
	Token        string
	RefreshToken string
}
