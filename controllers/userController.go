package controllers

import "github.com/gin-gonic/gin"

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func LogIn() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {

}

func VerifyPassword(userPassword, providePassword string) (bool, string) {

}
