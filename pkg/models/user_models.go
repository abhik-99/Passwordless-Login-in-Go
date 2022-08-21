package models

import (
	"passwordless-login/pkg/config"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Input validation for required values and such
type CreateUser struct {
	Pic           string `bson:"profilePic" json:"profilePic" validate:"uri"`
	Name          string `bson:"name" json:"name" validate:"required,max=30"`
	Email         string `bson:"email" json:"email" validate:"required,email"`
	Phone         string `bson:"phone" json:"phone" validate:"required,e164"`
	Address       string `bson:"addresss" json:"address" validate:"min=5,max=180"`
	Public        bool   `bson:"public" json:"public" validate:"boolean" default:"false"`
	HotpKey       string `bson:"secret"`
	SecretCounter uint64 `bson:"counter"`
}

// User to fetch user and create an OTP from the key
type CreateOtpForUser struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	HotpKey       string             `bson:"secret"`
	SecretCounter uint64             `bson:"counter"`
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
	ctx            = config.MongoCtx

	redisDb = config.Rdb
	rCtx    = config.RedisCtx
)

func (u *CreateUser) CreateNewUser() (*mongo.InsertOneResult, error) {
	key, err := hotp.Generate(hotp.GenerateOpts{Issuer: "Abhik Banerjee", AccountName: u.Name})
	if err != nil {
		return nil, err
	}
	u.HotpKey = key.String()
	u.SecretCounter = 0
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

func UserLookupViaEmail(email string) (bool, string, error) {
	var result CreateOtpForUser
	filter := bson.D{{Key: "email", Value: email}}
	projection := bson.D{{Key: "secret", Value: 1}, {Key: "counter", Value: 1}, {Key: "_id", Value: 1}}

	if err := userCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		return false, "", err
	} else {
		key, err := otp.NewKeyFromURL(result.HotpKey)
		if err != nil {
			return false, "", err
		}

		code, err := hotp.GenerateCode(key.String(), result.SecretCounter)
		if err != nil {
			return false, "", err
		}

		// Updating the counter for the next code
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "counter", Value: result.SecretCounter + 1}}}}
		if _, err = userCollection.UpdateByID(ctx, result.Id, update); err != nil {
			return false, "", err
		}

		return true, code, nil

	}
}

func UserLookupViaPhone(phone string) (bool, string, error) {
	var result CreateOtpForUser
	filter := bson.D{{Key: "phone", Value: phone}}
	projection := bson.D{{Key: "secret", Value: 1}, {Key: "counter", Value: 1}, {Key: "_id", Value: 1}}

	if err := userCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		return false, "", err
	} else {
		key, err := otp.NewKeyFromURL(result.HotpKey)
		if err != nil {
			return false, "", err
		}

		code, err := hotp.GenerateCode(key.String(), result.SecretCounter)
		if err != nil {
			return false, "", err
		}

		// Updating the counter for the next code
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "counter", Value: result.SecretCounter + 1}}}}
		if _, err = userCollection.UpdateByID(ctx, result.Id, update); err != nil {
			return false, "", err
		}

		return true, code, nil
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
