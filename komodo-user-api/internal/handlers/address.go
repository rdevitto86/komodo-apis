package handlers

import (
	"net/http"
)

type GetAddressesRequest struct {
	UserID string `json:"user_id"`
}

type GetAddressesResponse struct {
	Addresses []string `json:"addresses"`
}

// Returns all addresses for the authenticated user
func GetAddresses(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

type AddAddressRequest struct {
	UserID string `json:"user_id"`
	Address string `json:"address"`
}

type AddAddressResponse struct {
	AddressID string `json:"address_id"`
}


// Adds a new address for the authenticated user
func AddAddress(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

type UpdateAddressRequest struct {
	AddressID string `json:"address_id"`
	Address string `json:"address"`
}

type UpdateAddressResponse struct {
	AddressID string `json:"address_id"`
}

// Updates an address for the authenticated user
func UpdateAddress(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}

type DeleteAddressRequest struct {
	AddressID string `json:"address_id"`
}

type DeleteAddressResponse struct {
	AddressID string `json:"address_id"`
}

// Deletes an address for the authenticated user
func DeleteAddress(wtr http.ResponseWriter, req *http.Request) {
	// TODO: implement
}
