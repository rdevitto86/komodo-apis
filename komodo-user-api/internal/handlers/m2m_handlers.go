package handlers

import (
	"net/http"
)

type GetUserByIDRequest struct {
	UserID string `json:"user_id"`
	Size   string `json:"size"` // "basic" | "minimal" | "full"
}

type GetUserByIDResponse struct {
	UserID string `json:"user_id"`
	Size   string `json:"size"` // "basic" | "minimal" | "full"
}

// GetUserByID retrieves a user by their ID
func GetUserByID(wtr http.ResponseWriter, req *http.Request) {	
	// TODO: implement
}

type CreateUserRequest struct {
	UserID string `json:"user_id"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

// Creates a new user
func CreateUser(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

type UpdateUserByIDRequest struct {
	UserID string `json:"user_id"`
}

type UpdateUserByIDResponse struct {
	UserID string `json:"user_id"`
}

// Updates a user by their ID
func UpdateUserByID(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

type DeleteUserByIDRequest struct {
	UserID string `json:"user_id"`
}

type DeleteUserByIDResponse struct {
	UserID string `json:"user_id"`
}

// Deletes a user by their ID
func DeleteUserByID(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}
