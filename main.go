package main

import (
	"go-backend/controllers"
	"go-backend/models"
	"go-backend/storage"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_DATABASE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal(err)
	}

	err = models.MigrateEvents(db)

	if err != nil {
		log.Fatal(err)
	}

	r := controllers.Repository{
		DB: db,
	}

	app := fiber.New()
	app.Use(cors.New())
	r.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
