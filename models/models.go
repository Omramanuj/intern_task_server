package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Task string            `json:"task" bson:"task"`
}
