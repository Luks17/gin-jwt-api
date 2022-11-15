package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name,omitempty" bson:"name" validate:"required,min=2,max=200"`
	Password     *string            `json:"password,omitempty" bson:"password" validate:"required"`
	Email        *string            `json:"email,omitempty" bson:"email" validate:"required"`
	Token        *string            `json:"token" bson:"token"`
	UserType     *string            `json:"user_type,omitempty" bson:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string            `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	UserId       string             `json:"user_id,omitempty" bson:"user_id"`
}
