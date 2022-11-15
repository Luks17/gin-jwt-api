package controllers

import (
	"context"
	"golang-jwt/db"
	"golang-jwt/helpers"
	"golang-jwt/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var user_collection *mongo.Collection = db.OpenCollection(db.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(user_password string, provided_password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(provided_password), []byte(user_password))

	check := true
	msg := ""

	if err != nil {
		msg = "Email or Password is incorrect"
		check = false
	}

	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User

		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validation_err := validate.Struct(user)
		if validation_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err := user_collection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error ocurred while counting documents"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		token, refresh_token, _ := helpers.GenerateAllTokens(*user.Email, *user.Name, *user.UserType, user.UserId)
		user.Token = &token
		user.RefreshToken = &refresh_token

		result_insertion_number, insert_error := user_collection.InsertOne(ctx, user)
		if insert_error != nil {
			msg := "User was not created"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}

		c.JSON(http.StatusOK, result_insertion_number)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		var found_user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := user_collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&found_user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or Password is incorrect"})
			return
		}

		passwd_is_valid, msg := VerifyPassword(*user.Password, *found_user.Password)

		if !passwd_is_valid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refresh_token, _ := helpers.GenerateAllTokens(*found_user.Email, *found_user.Name, *found_user.UserType, found_user.UserId)
		helpers.UpdateAllTokens(token, refresh_token, found_user.UserId)

		c.JSON(http.StatusOK, found_user)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		record_per_page, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || record_per_page < 1 {
			record_per_page = 10
		}

		page, err_ := strconv.Atoi(c.Query("page"))
		if err_ != nil || page < 1 {
			page = 1
		}

		start_index := (page - 1) * record_per_page
		start_index, _ = strconv.Atoi(c.Query("startIndex"))

		match_stage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		group_stage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}}, {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}}, {Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
		project_stage := bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "total_count", Value: 1}, {Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", start_index, record_per_page}}}}}}}

		result, err := user_collection.Aggregate(ctx, mongo.Pipeline{match_stage, group_stage, project_stage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error ocurred while listing user items"})
		}

		var all_users []bson.M
		if err = result.All(ctx, &all_users); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, all_users[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")

		if err := helpers.MatchUserTypeToUUID(c, user_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		err := user_collection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
