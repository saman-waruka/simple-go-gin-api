package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

)

// User represents a user entity with an ID, Name, and Email.
// ID is a unique identifier (string), Name is required, and Email is required and must be a valid email address.
type User struct {
	ID    string `json:"id" example:"51d4fe21-ef0b-416d-90f5-f3d3b7b94452"`
	Name  string `json:"name" binding:"required" example:"Jovan Bednar"`
	Email string `json:"email" binding:"required,email" example:"Jovan.Bednar@gggmail.com"`
}


// File path for users
const usersFile = "users.json"

/**
 * User represents a user entity with ID, Name, and Email fields.
 *
 * @property {string} ID    - The unique identifier for the user.
 * @property {string} Name  - The name of the user. Required.
 * @property {string} Email - The email address of the user. Required, must be a valid email.
 */
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

/**
 * writeUsersToFile writes the provided slice of User structs to the users.json file.
 *
 * @param users The slice of User objects to be written to file.
 * @return error Returns an error if marshalling or writing to the file fails, otherwise returns nil.
 */
func writeUsersToFile(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(usersFile, data, 0644)
}

// @Summary      Get all users
// @Description  Retrieve the list of all users
// @Tags         users
// @Produce      json
// @Success      200  {array}   User
// @Failure      401  {object}  models.ErrorUnauthorizedResponse
// @Failure      500  {object}  models.ErrorInternalServerResponse
// @Param        x-authorization  header    string  true  "Authorization token"
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	users, err := readUsersFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read users"})
		return
	}
	c.JSON(http.StatusOK, users)
}


// @Summary      Create a new user
// @Description  Add a new user to the list
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      User  true  "User to create"
// @Param        x-authorization  header    string  true  "Authorization token"
// @Success      201   {object}  User
// @Failure      400   {object}  models.ErrorBadRequestCreateUserResponse
// @Failure      401  {object} models.ErrorUnauthorizedResponse
// @Failure      500   {object} models.ErrorInternalServerResponse
// @Router       /users [post]
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
