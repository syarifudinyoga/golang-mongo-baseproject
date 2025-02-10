package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang-mongodb/config"
	"golang-mongodb/models"
	"golang-mongodb/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary Register user
// @Description Mendaftarkan user baru dengan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.UserRegis true "User Data"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	var user models.UserRegis
	if err := c.ShouldBindJSON(&user); err != nil {
		// Ambil error dari validator
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := make(map[string]string)
			for _, fe := range ve {
				errorMessages[fe.Field()] = models.GetErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existingUser models.UserRegis
	err0 := config.GetCollection(os.Getenv("COLLECTION_USER")).FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err0 == nil {
		// Jika tidak ada error, berarti email sudah ada
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	user.ID = primitive.NewObjectID()

	if user.Role == "" {
		user.Role = "peserta"
	}
	_, err := config.GetCollection(os.Getenv("COLLECTION_USER")).InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Login godoc
// @Summary Login user
// @Description Login user dengan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.UserLogin true "User Login"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		// Ambil error dari validator
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := make(map[string]string)
			for _, fe := range ve {
				errorMessages[fe.Field()] = models.GetErrorMessage(fe)
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.UserLogin
	err := config.GetCollection(os.Getenv("COLLECTION_USER")).FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, _ := utils.GenerateToken(user.Email)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
