package service

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Incoming JSON struct that needs to be processed to award points
type ReceiptToProcess struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Items within the JSON struct
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
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

	points, err := assignPoints(recToProc)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		newID := uuid.New().String()
		receipts = append(receipts, Receipt{ID: newID, Points: points})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ProcessReceiptsResponse{ID: newID})
	}
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
		w.Write([]byte("\nThat ID was not found\n"))
	}
}

func assignPoints(receipt ReceiptToProcess) (int32, error) {
	var pointTotal int32
	pointTotal = 0
	// Points for retailer alphanumeric amt
	for _, byte := range []byte(receipt.Retailer) {
		if isAlphaNumeric(byte) {
			pointTotal++
		}
	}

	// Points for if total is a round number
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		return 0, err
	}

	if totalFloat == math.Trunc(totalFloat) {
		pointTotal += 50
	}

	// Points if total is multiple of 0.25
	if math.Remainder(totalFloat, 0.25) == 0 {
		pointTotal += 25
	}

	// Points for every pair of items on receipt
	pointTotal += int32(len(receipt.Items)/2) * 5

	// Points for multiple of 3 item description
	for _, item := range receipt.Items {
		trimmedStr := strings.Trim(item.ShortDescription, " ")
		if math.Remainder(float64(len(trimmedStr)), 3) == 0 {
			// Multiply price by 0.2 and round up to nearest integer
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			pointTotal += int32(math.Trunc(priceFloat*0.2)) + 1
		}
	}

	// Points for odd purchase day
	date, err := time.Parse(time.DateOnly, receipt.PurchaseDate)
	if err != nil {
		return 0, err
	}

	if math.Remainder(float64(date.Day()), 2) != 0 {
		pointTotal += 6
	}

	// Points for time of purchase between 2 and 4 pm
	receiptTime, err := time.Parse(time.TimeOnly, receipt.PurchaseTime+":00")
	if err != nil {
		return 0, err
	}
	afterTime, _ := time.Parse(time.TimeOnly, "14:00:00")
	beforeTime, _ := time.Parse(time.TimeOnly, "16:00:00")

	if receiptTime.After(afterTime) && receiptTime.Before(beforeTime) {
		pointTotal += 10
	}

	return pointTotal, nil
}

func isAlphaNumeric(c byte) bool {
	// Check if the byte value falls within the range of alphanumeric characters
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
