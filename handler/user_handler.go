package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}


// File path for users
const usersFile = "users.json"

// Read users from file
func readUsersFromFile() ([]User, error) {
	file, err := os.Open(usersFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var users []User
	err = json.NewDecoder(file).Decode(&users)
	return users, err
}

// Write users to file
func writeUsersToFile(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(usersFile, data, 0644)
}

// Remove the global users variable
// var users = ...

func GetUsers(c *gin.Context) {
	users, err := readUsersFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := readUsersFromFile()
	if err != nil {
		users = []User{} // If file doesn't exist, start with empty slice
	}

	newUser.ID = uuid.New().String()
	for _, user := range users {
		if user.Email == newUser.Email {
			c.JSON(http.StatusBadRequest, gin.H{"message": "This Name is already exists"})
			return
		}
	}

	users = append(users, newUser)

	if err := writeUsersToFile(users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
