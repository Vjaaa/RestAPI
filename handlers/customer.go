package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"restapi/db"
	"restapi/models"

	"github.com/gorilla/mux"
)

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	json.NewDecoder(r.Body).Decode(&customer)

	sqlStatement := `INSERT INTO customer (customer_id, customer_name) VALUES ($1, $2) RETURNING customer_id`
	err := db.DB.QueryRow(sqlStatement, customer.ID, customer.Name).Scan(&customer.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT * FROM customer")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}
	json.NewEncoder(w).Encode(customers)
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["customer_id"])

	var customer models.Customer
	sqlStatement := `SELECT * FROM customer WHERE customer_id=$1`
	row := db.DB.QueryRow(sqlStatement, id)
	switch err := row.Scan(&customer.ID, &customer.Name); err {
	case sql.ErrNoRows:
		http.Error(w, "Customer not found", http.StatusNotFound)
	case nil:
		json.NewEncoder(w).Encode(customer)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["customer_id"])

	var customer models.Customer
	json.NewDecoder(r.Body).Decode(&customer)

	sqlStatement := `UPDATE customer SET customer_name=$1 WHERE customer_id=$2`
	res, err := db.DB.Exec(sqlStatement, customer.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(customer)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["customer_id"])

	sqlStatement := `DELETE FROM customer WHERE customer_id=$1`
	res, err := db.DB.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
