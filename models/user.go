package models

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRegis struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `bson:"email" json:"email" binding:"required,email"`
	Password string             `bson:"password" json:"password" binding:"required"`
	Role     string             `bson:"role" json:"role" bindings:"required"`
}

type UserLogin struct {
	Email    string `bson:"email" json:"email" binding:"required,email"`
	Password string `bson:"password" json:"password" binding:"required"`
}

// Custom error messages
func GetErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " tidak boleh kosong"
	case "email":
		return "format email tidak valid " + err.Field()
	default:
		return "Format input tidak sesuai"
	}
}
