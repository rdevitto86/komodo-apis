package main

// TODO - Future Enhancement - add Azure SDK / cache call using Redis

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

// in-memory cache implementation (5min exp)
// TODO - have expiration set by env props
var Cache = cache.New(5*time.Minute, 5*time.Minute)

// defines a cachable item
type CacheItem []struct {
	ID   int
	Name string
}

// fetches an item from the catalog cache
func getItem(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] fetching item from catalog cache")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n\n", r.Method, r.URL, r.Header)

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad body / header verification
	// TODO - handle 404 if item is not found

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"id": 1,
		"catID": 86,
		"price": 3.99,
		"desc": "A test JSON using Go APIs",
		"overview": "A test JSON",
		"rating": 4,
		"enableRating": true,
		"enableReviews": false,
		"sku": "ABCDEF",
		"stock": 2,
		"enablePromotions": false
	}`))
}

// adds an item to the catalog cache
func addItem(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] adding item to catalog cache")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n- Body: %+v\n\n", r.Method, r.URL, r.Header, r.Body)

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad body / header verification
	// TODO - handle 500 if item fails to cache

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status": 200 }`))
}

// removes an item from the catalog cache
func removeItem(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] removing item from catalog cache")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n- Body: %+v\n\n", r.Method, r.URL, r.Header, r.Body)

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad body / header verification
	// TODO - handle 404 if item is not found

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status": 200 }`))
}

// TODO - look into creating repo/library for error handling (import for all apis)
// handles 404 exceptions
func notFound(w http.ResponseWriter, r *http.Request) {
	log.Println("[ERROR] unable to find resource with given URI param(s)")
	log.Printf("[ERROR]\n- %s: %s\n- Headers: %+v\n- Body: %+v\n\n", r.Method, r.URL, r.Header, r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{
		"status": 404,
		"message": "unable to find resource with given URI param(s)"
	}`))
}

// manages api request routing
func main() {
	// TODO - handle api environment properties
	// TODO - handle reverse proxy / gateway requests
	router := mux.NewRouter()
	api := router.PathPrefix("/cache/catalog-cache/v0.1").Subrouter()

	// 404 handlers
	api.HandleFunc("", notFound)
	api.HandleFunc("/", notFound)

	// request handers
	api.HandleFunc("/{key}", getItem).Methods(http.MethodGet)
	api.HandleFunc("/{key}", addItem).Methods(http.MethodPost)
	api.HandleFunc("/{key}", removeItem).Methods(http.MethodDelete)

	// TODO - handle dyanmic port numbering for environments
	// creates api router + listener on the specified port. Kills process if exception thrown.
	log.Fatal(http.ListenAndServe(":8080", router))
}
