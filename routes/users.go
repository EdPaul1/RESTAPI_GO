package routes

import (
	"log"
	"net/http"

	"example.com/rest-api/models"
	"example.com/rest-api/utils"
	"github.com/gin-gonic/gin"
)

func signUp(context *gin.Context) {
	var user models.User

	// Attempt to bind JSON to the user struct
	err := context.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}

	// Attempt to save the user
	err = user.Save()
	if err != nil {
		log.Printf("Error saving user: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not save user"})
		return
	}

	// Respond with success
	context.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}
func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data"})
		return
	}
	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}
	authToken, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user"})

		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": authToken})
}
