package errors

import (
	"fmt"
	"net/http"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("ERROR: %v\n", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
