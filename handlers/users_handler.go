package handlers

import (
	"context"
	// "errors"
	// "fmt"
	"net/http"
	"os"

	"golang-mongodb/config"
	"golang-mongodb/models"

	// "golang-mongodb/utils"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// GetAllUsers godoc
// @Summary Get all users
// @Description Mendapatkan semua data user
// @Tags Users
// @Produce json
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.UserRegis
	cursor, err := config.GetCollection(os.Getenv("COLLECTION_USER")).Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	for cursor.Next(context.TODO()) {
		var user models.UserRegis
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding user"})
			return
		}
		users = append(users, user)
	}

	// c.JSON(http.StatusOK, gin.H{"users": users})
	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update user
// @Description Memperbarui data user berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body models.UserRegis true "User Data"
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	// Konversi ID dari string ke primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.UserRegis
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash password jika diupdate
	if user.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	// Update data user berdasarkan ID
	_, err = config.GetCollection(os.Getenv("COLLECTION_USER")).UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID}, // Gunakan objectID yang sudah dikonversi
		bson.M{"$set": user},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Menghapus data user berdasarkan ID
// @Tags Users
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Konversi ID dari string ke primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Hapus user berdasarkan ID
	_, err = config.GetCollection(os.Getenv("COLLECTION_USER")).DeleteOne(
		context.TODO(),
		bson.M{"_id": objectID},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
