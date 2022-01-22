package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management/databases"
	"restaurant-management/helpers"
	"restaurant-management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = databases.OpenCollection(databases.Client, "user")

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
			return
		}
		page, anotherErr := strconv.Atoi(c.Query("page"))
		if anotherErr != nil || page < 1 {
			page = 1
			return
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{
			{"$match", bson.D{
				{},
			}},
		}

		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 0},
				{"user_item", bson.D{{"$slice", []interface{}{"$data", startIndex}}}},
			}},
		}
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			projectStage,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var allUser []bson.M
		if err = result.All(ctx, &allUser); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUser[0])
	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")
		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking email"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while cheking phone"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or phone already exist"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Update_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken := helpers.GenerateAllTokens(*&user.Email, *&user.First_name, *&user.Last_name, *&user.User_id)
		user.Token = token
		user.Refresh_token = refreshToken

		resultInsert, InsertErr := userCollection.InsertOne(ctx, user)
		if InsertErr != nil {
			msg := fmt.Sprintf("User was not created")
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, resultInsert)
	}

}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		validPassword, err := VerifyPassword(*user.Password, *foundUser.Password)
		if validPassword != true {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		token, refreshToken := helpers.GenerateAllTokens(*&foundUser.Email, *&foundUser.First_name, *&foundUser.Last_name)
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword, providePassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}
