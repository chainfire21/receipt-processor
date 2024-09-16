package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Incoming JSON struct that needs to be processed to award points
type ReceiptToProcess struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float32 `json:"total"`
}

// Items within the JSON struct
type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float32 `json:"price"`
}

// Response struct of GetPoints endpoint
type GetPointsResponse struct {
	Points int32 `json:"points"`
}

// Response struct of ProcessReceipts endpoint
type ProcessReceiptsResponse struct {
	ID string `json:"id"`
}

// Receipt struct with id of receipt and points awarded
type Receipt struct {
	ID     string `json:"id"`
	Points int32  `json:"points"`
}

// Making a slice of receipts to keep within memory
var receipts []Receipt

// HandleProcessReceipts will receive a JSON and returns a JSON object with an ID and awards points
func HandleProcessReceipts(w http.ResponseWriter, r *http.Request) {
	var recToProc ReceiptToProcess
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
	err = json.Unmarshal(body, &recToProc)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	// gotta add point logic now

	newID := uuid.New().String()
	receipts = append(receipts, Receipt{ID: newID, Points: 55})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ProcessReceiptsResponse{ID: newID})
}

// HandleGetPoints will receive the receipt ID from the URL and returns a JSON object with how many points have been awarded to said receipt
func HandleGetPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reqID := vars["id"]
	found := false
	for _, id := range receipts {
		if id.ID == reqID {
			found = true
			resp := GetPointsResponse{Points: id.Points}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			break
		}
	}
	if !found {
		w.WriteHeader(404)
		w.Write([]byte("That ID was not found"))
	}
}
