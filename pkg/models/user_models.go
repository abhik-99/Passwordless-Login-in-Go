package models

import (
	"passwordless-login/pkg/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Input validation for required values and such
type CreateUser struct {
	Pic     string `bson:"profilePic" json:"profilePic" validate:"uri"`
	Name    string `bson:"name" json:"name" validate:"required,max=30"`
	Email   string `bson:"email" json:"email" validate:"required,email"`
	Phone   string `bson:"phone" json:"phone" validate:"required,e164"`
	Address string `bson:"addresss" json:"address" validate:"min=5,max=180"`
	Public  bool   `bson:"public" json:"public" validate:"boolean" default:"false"`
	AuthVia uint8  `bson:"authVia" json:"authVia" validate:"required,min=0,max=2"`
}

// Omit Empty so that user's profile values are not set to default values during Patch req
type EditUser struct {
	Pic     string `bson:"profilePic,omitempty" json:"profilePic" validate:"uri"`
	Name    string `bson:"name,omitempty" json:"name" validate:"max=30"`
	Email   string `bson:"email,omitempty" json:"email" validate:"email"`
	Phone   string `bson:"phone,omitempty" json:"phone" validate:"e164"`
	Address string `bson:"addresss,omitempty" json:"address" validate:"min=5,max=180"`
	Public  bool   `bson:"public,omitempty" json:"public" validate:"boolean" default:"false"`
	AuthVia uint8  `bson:"authVia,omitempty" json:"authVia" validate:"min=0,max=2"`
}

// for returning User Profile if set to public
// for get req for all users
type PublicUserProfileResponse struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	Pic  string             `bson:"profilePic" json:"profilePic"`
	Name string             `bson:"name" json:"name"`
}

// for returning User's full Profile if set to public
// for get req with ID param
type PublicFullUserProfileResponse struct {
	Id      primitive.ObjectID `bson:"_id" json:"id"`
	Pic     string             `bson:"profilePic" json:"profilePic"`
	Name    string             `bson:"name" json:"name"`
	Email   string             `bson:"email" json:"email"`
	Phone   string             `bson:"phone" json:"phone"`
	Address string             `bson:"addresss" json:"address"`
}

// for returning user's own profile
type UserProfileResponse struct {
	Id      primitive.ObjectID `bson:"_id" json:"id"`
	Pic     string             `bson:"profilePic" json:"profilePic"`
	Name    string             `bson:"name" json:"name"`
	Email   string             `bson:"email" json:"email"`
	Phone   string             `bson:"phone" json:"phone"`
	Address string             `bson:"addresss" json:"address"`
	Public  bool               `bson:"public" json:"public"`
	AuthVia uint8              `bson:"authVia" json:"authVia"`
}

var (
	userCollection = config.Db.Collection("user-collection")
	ctx            = config.Ctx
)

func (u *CreateUser) CreateNewUser() (*mongo.InsertOneResult, error) {
	return userCollection.InsertOne(ctx, u)
}

func GetAllPublicUserProfiles() ([]PublicUserProfileResponse, error) {
	var users []PublicUserProfileResponse
	cursor, err := userCollection.Find(ctx, bson.M{"public": "true"})
	if err != nil {
		return users, err
	}

	if err = cursor.All(ctx, &users); err != nil {
		return users, nil
	} else {
		return users, err
	}
}

func GetFullPublicProfileById(id string) (PublicFullUserProfileResponse, error) {
	var user PublicFullUserProfileResponse
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: obId}, {Key: "public", Value: "true"}}).Decode(&user)
	return user, err
}

func GetPrivateUserProfile(id string) (UserProfileResponse, error) {
	var user UserProfileResponse
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": obId}).Decode(&user)
	return user, err
}

func UserLookupViaEmail(email string) (bool, error) {
	filter := bson.D{{Key: "email", Value: email}}

	if count, err := userCollection.CountDocuments(ctx, filter); err != nil {
		return count == 1, nil
	} else {
		return false, nil
	}
}

func UserLookupViaPhone(phone string) (bool, error) {
	filter := bson.D{{Key: "phone", Value: phone}}

	if count, err := userCollection.CountDocuments(ctx, filter); err != nil {
		return count == 1, nil
	} else {
		return false, nil
	}
}

func UpdateUserProfile(id string, u EditUser) (*mongo.UpdateResult, error) {
	if obId, err := primitive.ObjectIDFromHex(id); err != nil {
		update := bson.D{{Key: "$set", Value: u}}
		return userCollection.UpdateByID(ctx, obId, update)
	} else {
		return nil, err
	}
}

func DeleteUserProfile(id string) (*mongo.DeleteResult, error) {
	if obId, err := primitive.ObjectIDFromHex(id); err != nil {
		return userCollection.DeleteOne(ctx, bson.M{"_id": obId})
	} else {
		return nil, err
	}
}
