package database

import (
	"context"
	"log"
	"os"
	"time"
	 "go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/Omramanuj/intern_task_server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TasksCollection *mongo.Collection

func ConnectDB() *mongo.Client {
	connString := os.Getenv("db_conn_string")

	clientOptions := options.Client().ApplyURI(connString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Ping the database to verify connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	log.Printf("database connected successfully :D")
	
	// Initialize the tasks collection
	TasksCollection = client.Database("intern_task").Collection("tasks")
	
	return client
}

func SeedTasks() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Check if tasks already exist
	count, err := TasksCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to count documents: %v\n", err)
	}
	
	if count > 0 {
		log.Println("Tasks already exist, skipping seeding")
		return
	}
	
	tasks := []interface{}{
		models.Task{ID: primitive.NewObjectID(), Task: "Read code"},
		models.Task{ID: primitive.NewObjectID(), Task: "Test API"},
		models.Task{ID: primitive.NewObjectID(), Task: "Fill the form"},
	}
	
	_, err = TasksCollection.InsertMany(ctx, tasks)
	if err != nil {
		log.Fatalf("Failed to seed tasks: %v\n", err)
	}
	
	log.Println("Database seeded with initial tasks")
}
