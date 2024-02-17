package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Input validation for required values and such
type CreateUserDTO struct {
	Pic   string `bson:"profilePic" json:"profilePic" validate:"uri"`
	Name  string `bson:"name" json:"name" validate:"required,max=30"`
	Email string `bson:"email" json:"email" validate:"required,email"`
	Phone string `bson:"phone" json:"phone" validate:"required,e164"`
}

type EditUserDTO struct {
	CreateUserDTO
}

// for returning User Profile if set to public
// for get req for all users
type PublicUserProfileDTO struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	Pic  string             `bson:"profilePic" json:"profilePic"`
	Name string             `bson:"name" json:"name"`
}

// for returning User's own profile
type PublicFullUserProfileDTO struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Pic   string             `bson:"profilePic" json:"profilePic"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Phone string             `bson:"phone" json:"phone"`
}
