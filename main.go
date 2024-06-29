package main

import (
	"log"
	"net/http"

	"restapi/db"
	"restapi/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.Init()

	r := mux.NewRouter()

	r.HandleFunc("/customer", handlers.CreateCustomer).Methods("POST")
	r.HandleFunc("/customer", handlers.GetCustomers).Methods("GET")
	r.HandleFunc("/customer/{customer_id}", handlers.GetCustomer).Methods("GET")
	r.HandleFunc("/customer/{customer_id}", handlers.UpdateCustomer).Methods("PUT")
	r.HandleFunc("/customer/{customer_id}", handlers.DeleteCustomer).Methods("DELETE")

	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
