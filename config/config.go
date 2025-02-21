package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	APP_ENV string
	PORT    string

	DATABASE_HOST     string
	DATABASE_PORT     string
	DATABASE_NAME     string
	DATABASE_USERNAME string
	DATABASE_PASSWORD string

	JWT_SECRET string
)

func LoadEnv() {
	// load environment
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	APP_ENV = os.Getenv("APP_ENV")
	PORT = os.Getenv("PORT")

	DATABASE_HOST = os.Getenv("DATABASE_HOST")
	DATABASE_PORT = os.Getenv("DATABASE_PORT")
	DATABASE_NAME = os.Getenv("DATABASE_NAME")
	DATABASE_USERNAME = os.Getenv("DATABASE_USERNAME")
	DATABASE_PASSWORD = os.Getenv("DATABASE_PASSWORD")

	JWT_SECRET = os.Getenv("JWT_SECRET")
}
