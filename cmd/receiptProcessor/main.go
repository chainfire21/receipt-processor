package main

import (
	"log"
	"net/http"
	service "receipt-processor/service"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", service.HandleProcessReceipts).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", service.HandleGetPoints).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())
}
