package main

import (
	"context"
	"log"

	"github.com/Omramanuj/intern_task_server/database"
	"github.com/Omramanuj/intern_task_server/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	log.Printf("server starting ....")

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "Upgrade",
	}))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := database.ConnectDB()
	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	database.SeedTasks()

	// Set up routes (only once)
	app.Get("/api/tasks", handlers.GetTasks)
	app.Get("/api/tasks/:id", handlers.GetTaskByID)

	log.Println("Routes set up successfully")
	log.Fatal(app.Listen(":8080"))
}
