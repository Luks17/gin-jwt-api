package helpers

import (
	"context"
	"golang-jwt/config"
	"golang-jwt/db"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Name     string
	UserId   string
	UserType string
	jwt.RegisteredClaims
}

var user_collection *mongo.Collection = db.OpenCollection(db.Client, "user")

var SECRET_KEY string = config.GetSecret()

func ValidateToken(signed_token string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signed_token,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		msg = "Token is expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signed_token string, signed_refresh_token string, user_id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var update_obj primitive.D

	update_obj = append(update_obj, bson.E{Key: "token", Value: signed_token})
	update_obj = append(update_obj, bson.E{Key: "refresh_token", Value: signed_refresh_token})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update_obj = append(update_obj, bson.E{Key: "updated_at", Value: updated_at})

	filter := bson.D{{Key: "user_id", Value: user_id}}
	opt := options.Update().SetUpsert(true)

	_, err := user_collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: update_obj},
		},
		opt,
	)

	if err != nil {
		log.Panic(err)
		return
	}
}

func GenerateAllTokens(email string, name string, user_type string, user_id string) (signed_token string, signed_refresh_token string, err error) {
	claims := &SignedDetails{
		Email:    email,
		Name:     name,
		UserType: user_type,
		UserId:   user_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24)),
		},
	}

	refresh_claims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 168)),
		},
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims).SignedString([]byte(SECRET_KEY))

	return token, refresh_token, err
}
