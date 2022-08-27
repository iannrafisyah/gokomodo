package cmd

import (
	"context"

	"github.com/iannrafisyah/gokomodo/database/postgres"
	"github.com/iannrafisyah/gokomodo/module"
	"github.com/iannrafisyah/gokomodo/package/logger"
	"github.com/iannrafisyah/gokomodo/router"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var app = &cobra.Command{
	Use:   "start",
	Short: "Running service",
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(
			fx.Provide(router.NewRouter),
			fx.Provide(postgres.NewPostgres),
			fx.Provide(logger.NewLogRus),
			module.BundleRepository,
			module.BundleLogic,
			module.BundleRoute,
			fx.Invoke(registerHooks),
		).Run()
	},
}

func init() {
	rootCmd.AddCommand(app)
}

func registerHooks(lifecycle fx.Lifecycle,
	e *router.Router,
	db *postgres.DB,
	log *logger.LogRus) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go e.Start(":8081")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := e.Shutdown(ctx); err != nil {
					log.Fatalf("registerHooks", err.Error())
					return err
				}
				defer db.Sql.Close()
				return nil
			},
		},
	)
}
