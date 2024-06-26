package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"user-service/models"
	"user-service/services"
)

// SignUp creates a new user account in the system with the provided details
func SignUp(c *gin.Context) {
	// Create a new user
	var user models.SignUpUserInput
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userCreated, err := services.SignUpService(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userResponse := models.UserResponse{
		Name:          userCreated.Name,
		LastName:      userCreated.LastName,
		Email:         userCreated.Email,
		Address:       userCreated.Address,
		PaymentMethod: userCreated.PaymentMethod,
		Role:          userCreated.Role,
	}
	c.JSON(http.StatusCreated, gin.H{"data": userResponse})
}

// SignIn authenticates a user by email and password
func SignIn(c *gin.Context) {
	// Get the email and password from the request
	var user models.SignInUserInput
	jwtService := services.NewJWTService()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userAuthenticated, err := services.SignInService(user)
	generateToken, _ := jwtService.GenerateToken(userAuthenticated.ID, userAuthenticated.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userResponse := models.UserResponse{
		Name:          userAuthenticated.Name,
		LastName:      userAuthenticated.LastName,
		Email:         userAuthenticated.Email,
		Address:       userAuthenticated.Address,
		PaymentMethod: userAuthenticated.PaymentMethod,
		Role:          userAuthenticated.Role,
	}
	c.JSON(http.StatusOK, gin.H{"data": userResponse, "token": generateToken})
}

// DeleteUserByID deletes a user from the system by its ID
func DeleteUserByID(c *gin.Context) {
	jwt := c.GetHeader("Authorization")
	jwtService := services.NewJWTService()
	jwtDecode, _ := jwtService.DecodeToken(strings.Replace(jwt, "Bearer ", "", 1))
	id := c.Param("id")
	if jwtDecode.ID == "" || jwtDecode.ID != id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	err := services.DeleteUserByIDService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
