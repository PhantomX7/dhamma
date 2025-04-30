package main

import (
	"fmt"
	"log"
	"os"

	"github.com/PhantomX7/dhamma/seeder/seed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/subosito/gotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/migration"
)

func main() {
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FJakarta",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		//DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	if err = migration.RunMigration(db); err != nil {
		panic(err)
	}

	userSeeder := seed.NewUserSeeder(db)

	err = userSeeder.GenerateRootUser()
	if err != nil {
		panic(err)
	}

	permissionSeeder := seed.NewPermissionSeeder(db)

	err = permissionSeeder.GenerateApiPermissions()
	if err != nil {
		panic(err)
	}

	err = permissionSeeder.SyncPermissions()
	if err != nil {
		panic(err)
	}

	// development seed only, will not run on production
	if os.Getenv("APP_ENV") == "development" {
		//// Generate 10 basic users
		//err = userSeeder.GenerateUsers(20)
		//if err != nil {
		//	panic(err)
		//}
		//
		//// Generate 5 users with custom options
		//err = userSeeder.GenerateUsers(5,
		//	seed.WithActiveStatus(false),
		//)
		//if err != nil {
		//	panic(err)
		//}
	}

	log.Println("finish seeding")
}
