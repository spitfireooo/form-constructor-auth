package main

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/spitfireooo/form-constructor-auth/internal/config"
	"github.com/spitfireooo/form-constructor-auth/internal/database"
	"github.com/spitfireooo/form-constructor-auth/pkg/router"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error in load env file", err)
	}

	if err := config.ConfigInit(); err != nil {
		port += "8020"
		log.Fatal("Error in configuration init", err)
	} else {
		port += viper.GetString("http.auth_port")
	}

	if err := database.DatabaseInit(database.DBConfig{
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: viper.GetString("db.database"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		SSLMode:  viper.GetString("db.ssl_mode"),
	}); err != nil {
		log.Fatal("Error in database connection", err)
	}
}

var port = ":"

func main() {
	app := fiber.New()

	app.Static("/static", "./static")

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		//AllowOrigins: "http://localhost:8040",
		AllowOrigins: "*",
		AllowMethods: "Origin, Content-Type, Accept",
	}))
	app.Use(swagger.New(swagger.Config{
		BasePath: "/auth",
		FilePath: "./docs/swagger.json",
		Path:     "docs",
	}))

	router.Router(app)

	if err := app.Listen(port); err != nil {
		log.Fatal("Error in server started", err)
		return
	}
}
