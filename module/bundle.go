package module

import (
	//Route
	authRoute "github.com/iannrafisyah/gokomodo/module/auth/route"
	productRoute "github.com/iannrafisyah/gokomodo/module/product/route"
	transactionRoute "github.com/iannrafisyah/gokomodo/module/transaction/route"

	//Logic
	authLogic "github.com/iannrafisyah/gokomodo/module/auth/logic"
	productLogic "github.com/iannrafisyah/gokomodo/module/product/logic"
	transactionLogic "github.com/iannrafisyah/gokomodo/module/transaction/logic"
	userLogic "github.com/iannrafisyah/gokomodo/module/user/logic"

	//Repository
	productRepository "github.com/iannrafisyah/gokomodo/module/product/repository"
	transactionRepository "github.com/iannrafisyah/gokomodo/module/transaction/repository"
	userRepository "github.com/iannrafisyah/gokomodo/module/user/repository"

	"go.uber.org/fx"
)

// Register Route
var BundleRoute = fx.Options(
	fx.Invoke(transactionRoute.NewRoute),
	fx.Invoke(productRoute.NewRoute),
	fx.Invoke(authRoute.NewRoute),
)

// Register logic
var BundleLogic = fx.Options(
	fx.Provide(userLogic.NewLogic),
	fx.Provide(transactionLogic.NewLogic),
	fx.Provide(productLogic.NewLogic),
	fx.Provide(authLogic.NewLogic),
)

// Register Repository
var BundleRepository = fx.Options(
	fx.Provide(userRepository.NewRepository),
	fx.Provide(transactionRepository.NewRepository),
	fx.Provide(productRepository.NewRepository),
)
