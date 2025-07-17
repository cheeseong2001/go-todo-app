package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/cheeseong2001/auth-service/models"
	"github.com/cheeseong2001/auth-service/repository"
	"github.com/cheeseong2001/auth-service/utils"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var registerCredentials Credentials
	if err := c.ShouldBindJSON(&registerCredentials); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := utils.HashPassword(registerCredentials.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	newUser := models.User{
		Email:    registerCredentials.Email,
		Password: hashedPassword,
	}

	err = repository.DB.QueryRow(`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`, newUser.Email, newUser.Password).Scan(&newUser.ID)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create new user",
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"user": newUser})
}

func Login(c *gin.Context) {
	var loginCredentials Credentials

	if err := c.ShouldBindJSON(&loginCredentials); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var hashedPassword string
	var userID int
	err := repository.DB.QueryRow(`SELECT id, password FROM users WHERE email = $1`, loginCredentials.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "Database error",
		})
		return
	}

	if !utils.CheckPasswordHash(loginCredentials.Password, hashedPassword) {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	token, err := utils.GenerateJWT(userID)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func ValidateToken(c *gin.Context) {
	var request models.ValidateTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, models.ValidateTokenResponse{
			Valid: false,
			Error: "Invalid request format",
		})
		return
	}

	userID, err := utils.ValidateJWT(request.Token)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, models.ValidateTokenResponse{
			Valid: false,
			Error: "Invalid or expired token",
		})
		return
	}

	c.IndentedJSON(http.StatusBadRequest, models.ValidateTokenResponse{
		UserID: userID,
		Valid:  true,
	})
}
