package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/libs"
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/migration"
	"github.com/PhantomX7/dhamma/modules"
	"github.com/PhantomX7/dhamma/routes"
	customValidator "github.com/PhantomX7/dhamma/utility/validator"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// func args(db *gorm.DB) bool {
// 	if len(os.Args) > 1 {
// 		flag := command.Commands(db)
// 		if !flag {
// 			return false
// 		}
// 	}

// 	return true
// }

func main() {

	// load config from .env
	config.LoadEnv()

	// if !args(db) {
	// 	return
	// }
	app := fx.New(
		// fx.NopLogger, // disable logger for fx
		fx.Provide(
			setupDatabase,
			setUpServer,
			customValidator.New, // initiate custom validator
			middleware.New,      // initiate middleware
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

func startServer(
	lc fx.Lifecycle,
	server *gin.Engine,
	db *gorm.DB,
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
				log.Printf("api is available at %s\n", srv.Addr)

				// service connections
				if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("listen: %s\n", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")

			// Create a timeout context for shutdown
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// Shutdown the server
			if err := srv.Shutdown(ctx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
				return err
			}

			log.Println("Server shutdown completed")
			return nil
		},
	})

}

func setUpServer() *gin.Engine {
	if config.APP_ENV == constants.EnumRunProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	// server.Use(middleware.CORSMiddleware())

	server.Static("/assets", "./assets")

	return server
}

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
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				IgnoreRecordNotFoundError: true,
			},
		),
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// run migration
	if err = migration.RunMigration(db); err != nil {
		log.Println(err)
		panic(err)
	}

	// Database lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Ensuring database connection...")
			db, _ := db.DB()
			return db.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection...")
			db, _ := db.DB()
			return db.Close()
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
			log.Println("Starting Cron")

			cron.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping cron")

			return cron.Shutdown()
		},
	})
}
