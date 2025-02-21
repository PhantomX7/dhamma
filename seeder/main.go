package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/subosito/gotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/PhantomX7/dhamma/migration"
	"github.com/PhantomX7/dhamma/seeder/seed"
)

func main() {
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	if err = migration.RunMigration(db); err != nil {
		panic(err)
	}

	if err = seed.SeedRootUser(db); err != nil {
		panic(err)
	}
	log.Print("finish seeding user")

	if err = seed.SeedConfig(db); err != nil {
		panic(err)
	}
	log.Print("finish seeding config")

	// development seed only, will not run on production
	if os.Getenv("APP_ENV") == "development" {

	}

	log.Println("finish seeding")
}
