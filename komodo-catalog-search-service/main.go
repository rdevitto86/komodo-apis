package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// fetches a catalog item's information
func getCatalogItem(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] processing getCatalogItem request")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n\n", r.Method, r.URL, r.Header)

	// pathParams := mux.Vars(r)
	// var err error

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad item ID
	// TODO - check cache service before fetching from DB

	// if val, ok := pathParams["userID"]; ok {
	// 	userID, err = strconv.Atoi(val)
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		w.Write([]byte(`{"message": "need a number"}`))
	// 		return
	// 	}
	// }

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

// gets all catalog items for a given category
func getCategoryItems(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] processing getCategoryItems request")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n\n", r.Method, r.URL, r.Header)

	w.Header().Set("Content-Type", "application/json")

	var err error

	pathParams := mux.Vars(r)
	catID, err := strconv.Atoi(pathParams["id"])

	// TODO - do param santization
	// TODO - do param type validation

	if err != nil {
		log.Printf("[ERROR] %+v", catID)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{
			"status": 500,
			"message": "failed to parse category id"
		}`))
		return
	}

	// TODO - check cache service before fetching from DB

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"catID": %d,
		"size": 1,
		"items": [
			{
				"id": 1,
				"catID": %d,
				"price": 3.99,
				"desc": "A test JSON using Golang APIs",
				"overview": "A test JSON",
				"rating": 4,
				"enableRating": true,
				"enableReviews": false,
				"sku": "ABCDEF",
				"stock": 2,
				"enablePromotions": false
			}
		]
	}`, catID, catID)))
}

// searches for a catalog item with given params
func search(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] processing catalog item search request")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n\n", r.Method, r.URL, r.Header)

	w.Header().Set("Content-Type", "application/json")

	// queryParams := r.URL.Query()

	// TODO - write error handler for bad params
	// TODO - check cache service before fetching from DB
	// TODO - implement some form of search algorithm
}

// fetches specific catalog item review
func getItemReview(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] processing getItemReview request")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n\n", r.Method, r.URL, r.Header)

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad params
	// TODO - check cache service before fetching from DB
}

// fetches all catalog item reviews
func getItemReviews(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] processing getItemReviews request")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n\n", r.Method, r.URL, r.Header)

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad params
	// TODO - check cache service before fetching from DB
}

// adds/updates a user review for a given item
func submitReview(w http.ResponseWriter, r *http.Request) {
	log.Println("[INFO] processing submitReview request")
	log.Printf("[INFO]\n- %s: %s\n- Headers: %+v\n- Body: %+v\n\n", r.Method, r.URL, r.Header, r.Body)

	w.Header().Set("Content-Type", "application/json")

	// TODO - write error handler for bad body / header verification

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status": 200 }`))
}

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

// handles 400 exceptions
func badRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("[ERROR] malformed request received")
	log.Printf("[ERROR]\n- %s: %s\n- Headers: %+v\n- Body: %+v\n\n", r.Method, r.URL, r.Header, r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{
		"status": 400,
		"message": "malformed request received"
	}`))
}

// manages api request routing
func main() {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/catalog-search/v0.1").Subrouter()

	// 404 handlers
	api.HandleFunc("", badRequest)
	api.HandleFunc("/", badRequest)
	api.HandleFunc("/category", notFound)
	api.HandleFunc("/item", notFound)
	api.HandleFunc("/item/{itemID}/review", notFound)
	api.HandleFunc("/item/{itemID}/submit", notFound)

	// GET handers
	api.HandleFunc("/search", search).Methods(http.MethodGet)
	api.HandleFunc("/category/{id}", getCategoryItems).Methods(http.MethodGet)
	api.HandleFunc("/item/{id}", getCatalogItem).Methods(http.MethodGet)
	api.HandleFunc("/item/{itemID}/review/{reviewID}", getItemReview).Methods(http.MethodGet)
	api.HandleFunc("/item/{itemID}/reviews", getItemReviews).Methods(http.MethodGet)

	// POST handlers
	api.HandleFunc("/item/{itemID}/submit/review", submitReview).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// // helper function that fetches item stored in the catalog DB
// func _fetchFromDB() {
// 	// TODO - access data and fetch data
// 	// TODO - cache data after fetch (if applicable)
// }

// // helper function that fetches item stored in the catalog cache
// func _fetchFromCache(key string) {
// 	// TODO - check for cached item
// 	// TODO - return nil if item not found
// }
