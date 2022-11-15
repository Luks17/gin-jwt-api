package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name,omitempty" validate:"required,min=2,max=200"`
	Password     *string            `json:"password,omitempty" validate:"required"`
	Email        *string            `json:"email,omitempty" validate:"required"`
	Token        *string            `json:"token"`
	UserType     *string            `json:"user_type,omitempty" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserId       string             `json:"user_id"`
}
