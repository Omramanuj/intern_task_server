package handlers

import (
	"context"
	"time"

	"github.com/Omramanuj/intern_task_server/database"
	"github.com/Omramanuj/intern_task_server/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskResponse struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

func GetTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.TasksCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch tasks",
		})
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode tasks",
		})
	}

	response := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = TaskResponse{
			ID:   task.ID.Hex(),
			Task: task.Task,
		}
	}

	return c.JSON(response)
}

func GetTaskByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err = database.TasksCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	response := TaskResponse{
		ID:   task.ID.Hex(),
		Task: task.Task,
	}

	return c.JSON(response)
}

