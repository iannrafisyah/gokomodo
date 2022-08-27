package main_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iannrafisyah/gokomodo/config"
	"github.com/iannrafisyah/gokomodo/database/postgres"
	"github.com/iannrafisyah/gokomodo/enum"
	"github.com/iannrafisyah/gokomodo/module"
	transactionRoute "github.com/iannrafisyah/gokomodo/module/transaction/route"
	"github.com/iannrafisyah/gokomodo/package/jwt"
	"github.com/iannrafisyah/gokomodo/package/logger"
	"github.com/iannrafisyah/gokomodo/router"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestMain(t *testing.T) {
	config.SetConfig()
	fx.New(
		fx.Provide(router.NewRouter),
		fx.Provide(postgres.NewPostgres),
		fx.Provide(logger.NewLogRus),
		module.BundleRepository,
		module.BundleLogic,
		module.BundleRoute,
		fx.Invoke(NewRouteTest),
	).Start(context.Background())
}

type RouteTest struct {
	fx.In
	TransactionHandler transactionRoute.Handler
}

var r RouteTest

func NewRouteTest(routeTest RouteTest) {
	r = routeTest
}

func TestTransaction(t *testing.T) {
	t.Run("FailedSellerCreateOrder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"SellerID":1,
			"Items":[1,2]
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 1,
			Role:   enum.RoleTypeSeller,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.CreateOrder(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("SuccessBuyerCreateOrder", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
			"SellerID":1,
			"Items":[1,2]
		}`))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		e := echo.New()
		c := e.NewContext(req, rec)

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, jwt.InternalClaimData{}, jwt.InternalClaimData{
			UserID: 2,
			Role:   enum.RoleTypeBuyer,
		})
		c.SetRequest(c.Request().WithContext(ctx))

		// Assertions
		if assert.NoError(t, r.TransactionHandler.CreateOrder(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}
