package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/migration"
	"github.com/PhantomX7/dhamma/modules"
	"github.com/PhantomX7/dhamma/routes"
	customValidator "github.com/PhantomX7/dhamma/utility/validator"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		// fx.NopLogger, // disable logger
		fx.Provide(
			setupDatabase,
			setUpServer,
			customValidator.New,
			// initLibs,
		),
		modules.RepositoryModule,
		modules.ServiceModule,
		modules.ControllerModule,
		routes.Module,
		fx.Invoke(
			//startCron,
			//startQueue,
			startServer,
		),
	)
	app.Run()

	// server := gin.Default()
	// server.Use(middleware.CORSMiddleware())

	// routes
	// routes.User(server, userController, jwtService)

}

func startServer(lc fx.Lifecycle, server *gin.Engine, db *gorm.DB, cv customValidator.CustomValidator) {
	myFigure := figure.NewColorFigure("Phantom", "", "green", true)
	myFigure.Print()

	// register all custom validator here
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("unique",
			cv.Unique())
		if err != nil {
			log.Println("error when applying unique validator")
		}
		err = v.RegisterValidation("exist", cv.Exist())
		if err != nil {
			log.Println("error when applying exist validator")
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// start server

			return nil
		},
		OnStop: func(context.Context) error {
			// close database connection
			dbSQL, err := db.DB()
			if err != nil {
				panic(err)
			}
			dbSQL.Close()

			return nil
		},
	})

	if err := server.Run(fmt.Sprintf(":%s", config.PORT)); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func setUpServer() *gin.Engine {
	if config.APP_ENV == constants.ENUM_RUN_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	// server.Use(middleware.CORSMiddleware())

	server.Static("/assets", "./assets")

	return server
}

func setupDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		config.DATABASE_USERNAME,
		config.DATABASE_PASSWORD,
		config.DATABASE_HOST,
		config.DATABASE_PORT,
		config.DATABASE_NAME,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if err = migration.RunMigration(db); err != nil {
		log.Println(err)
		panic(err)
	}
	return db
}
