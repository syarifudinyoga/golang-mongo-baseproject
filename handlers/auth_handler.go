package handlers

import (
	"context"
	"net/http"

	"golang-mongodb/config"
	"golang-mongodb/models"
	"golang-mongodb/utils"

	"github.com/gin-gonic/gin"
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
	var user models.UserRegis
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	user.ID = primitive.NewObjectID()

	if user.Role == "" {
		user.Role = "peserta"
	}
	_, err := config.GetCollection("baseproject").InsertOne(context.TODO(), user)
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
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.UserLogin
	err := config.GetCollection("baseproject").FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&user)
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
