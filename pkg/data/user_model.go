package data

import (
	"time"

	"github.com/abhik-99/passwordless-login/pkg/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Pic   string             `bson:"profilePic" json:"profilePic"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Phone string             `bson:"phone" json:"phone"`
}

var (
	userCollection = config.Db.Collection("user-collection")
	ctx            = config.MongoCtx

	redisDb = config.Rdb
	rCtx    = config.RedisCtx
)

func CreateNewUser(user CreateUserDTO) (*mongo.InsertOneResult, error) {
	return userCollection.InsertOne(ctx, user)
}

func GetAllPublicUserProfiles() ([]PublicUserProfileDTO, error) {
	var users []PublicUserProfileDTO
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

func GetPrivateUserProfile(id string) (User, error) {
	var user User
	obId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": obId}).Decode(&user)
	return user, err
}

func UserLookupViaEmail(email string) (bool, string, error) {
	var result User
	filter := bson.D{{Key: "email", Value: email}}
	projection := bson.D{{Key: "_id", Value: 1}}

	if err := userCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		return false, "", err
	}
	return true, result.Id.Hex(), nil

}

func UserLookupViaPhone(phone string) (bool, string, error) {
	var result User
	filter := bson.D{{Key: "phone", Value: phone}}
	projection := bson.D{{Key: "secret", Value: 1}, {Key: "counter", Value: 1}, {Key: "_id", Value: 1}}

	if err := userCollection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result); err != nil {
		return false, "", err
	}
	return true, result.Id.Hex(), nil
}

func UpdateUserProfile(id string, u EditUserDTO) (*mongo.UpdateResult, error) {
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

func SetOTPForUser(userId string, otp string) error {
	return redisDb.Set(rCtx, userId, otp, 30*time.Minute).Err()
}

func CheckOTP(userId string, otp string) (bool, error) {
	if storedOtp, err := redisDb.Get(rCtx, userId).Result(); err != nil {
		return false, err
	} else {
		if storedOtp == otp {
			return true, nil
		}
	}
	return false, nil
}
