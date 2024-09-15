package service

import (
	"fmt"
	"net/http"
)

// HandleProcessReceipts will receive a JSON and returns a JSON object with an ID and awards points
func HandleProcessReceipts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi")
}

// HandleGetPoints will receive the receipt ID from the URL and return how many points have been awarded to said receipt
func HandleGetPoints(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi 2")
}
