package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"user-crud-api/config"
	"user-crud-api/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Create User
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	// Create user in DB
	result := config.DB.Create(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
		return
	}

	// Return the full user with real ID
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Get User by ID
// GetAllUsers returns all users in the database
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	// Fetch all users from DB
	result := config.DB.Find(&users)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
		return
	}

	// Return users as JSON
	json.NewEncoder(w).Encode(users)
}

// Update User
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if result := config.DB.First(&user, id); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}

	var updatedData models.User
	err = json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	user.Username = updatedData.Username
	if updatedData.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	config.DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

// Delete User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	if result := config.DB.Delete(&models.User{}, id); result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
