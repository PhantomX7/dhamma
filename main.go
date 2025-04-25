package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/libs"
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/migration"
	"github.com/PhantomX7/dhamma/modules"
	"github.com/PhantomX7/dhamma/routes"
	"github.com/PhantomX7/dhamma/utility/logger"
	customValidator "github.com/PhantomX7/dhamma/utility/validator"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// initialize logger
	logger.NewLogger()

	// load config from .env
	config.LoadEnv()

	app := fx.New(
		fx.NopLogger, // disable logger for fx
		fx.Provide(
			setupDatabase,
			customValidator.New, // initiate custom validator
			middleware.New,      // initiate middleware
			setupServer,
		),
		modules.RepositoryModule,
		modules.ServiceModule,
		modules.ControllerModule,
		libs.Module,
		routes.Module,
		fx.Invoke(
			startCron,
			startServer,
		),
	)
	app.Run()
}

func setupServer(m *middleware.Middleware) *gin.Engine {
	// set gin mode
	if config.APP_ENV == constants.EnumRunProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()

	// Enable CORS middleware
	server.Use(m.CORS(), m.Logger())

	// register static files
	server.Static("/assets", "./assets")
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return server
}

// startServer initializes and starts the HTTP server.
func startServer(
	lc fx.Lifecycle,
	server *gin.Engine,
	cv customValidator.CustomValidator,
) {
	myFigure := figure.NewColorFigure("Phantom", "", "green", true)
	myFigure.Print()

	// list of custom validators
	validators := map[string]validator.Func{
		"unique": cv.Unique(),
		"exist":  cv.Exist(),
	}
	registerValidators(validators)

	// Initialize HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.PORT),
		Handler: server.Handler(),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Server lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				logger.Get().Info("api is available", zap.String("address", srv.Addr))

				// service connections
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Get().Fatal("listen error", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Get().Info("Shutting down HTTP server...")

			// Create a timeout context for shutdown
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Shutdown the server
			if err := srv.Shutdown(ctx); err != nil {
				logger.Get().Error("Server forced to shutdown", zap.Error(err))
				return err
			}

			logger.Get().Info("Server shutdown completed")
			// Flush any buffered log entries
			logger.Get().Sync()
			return nil
		},
	})

}

// setupDatabase initializes the database connection and runs migrations.
func setupDatabase(lc fx.Lifecycle) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FJakarta",
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_HOST,
		config.DATABASE_PORT,
		config.DATABASE_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		//DisableForeignKeyConstraintWhenMigrating: true,
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Keep standard log for GORM for now, or replace with a zap-compatible GORM logger if available/needed
			gormLogger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		logger.Get().Fatal("error initializing database", zap.Error(err))
	}

	// run migration
	if err = migration.RunMigration(db); err != nil {
		logger.Get().Fatal("error running migration", zap.Error(err))
	}

	// Database lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Get().Info("Ensuring database connection...")
			dbInstance, _ := db.DB()
			return dbInstance.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Get().Info("Closing database connection...")
			dbInstance, _ := db.DB()
			return dbInstance.Close()
		},
	})

	return db
}

func registerValidators(validators map[string]validator.Func) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register each validator
		for name, fn := range validators {
			if err := v.RegisterValidation(name, fn); err != nil {
				log.Printf("error when applying %s validator: %v", name, err)
			}
		}
	}
}

func startCron(lc fx.Lifecycle, cron gocron.Scheduler) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Get().Info("Starting Cron")

			cron.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Get().Info("Stopping cron")

			return cron.Shutdown()
		},
	})
}
