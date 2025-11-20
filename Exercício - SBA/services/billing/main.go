package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Invoice struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	OrderID     string    `json:"order_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	PaymentDate time.Time `json:"payment_date"`
}

var invoices = []Invoice{
	{ID: "INV-001", UserID: "1", OrderID: "1001", Amount: 3500.00, Status: "paid", PaymentDate: time.Now().AddDate(0, 0, -10)},
	{ID: "INV-002", UserID: "2", OrderID: "1002", Amount: 100.00, Status: "pending", PaymentDate: time.Time{}},
	{ID: "INV-003", UserID: "1", OrderID: "1003", Amount: 250.00, Status: "paid", PaymentDate: time.Now().AddDate(0, 0, -2)},
}

func getInvoices(w http.ResponseWriter, r *http.Request) {
	log.Println("[BILLING SERVICE] GET /invoices")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func getInvoiceByID(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.URL.Query().Get("id")
	log.Printf("[BILLING SERVICE] GET /invoice?id=%s\n", invoiceID)

	for _, invoice := range invoices {
		if invoice.ID == invoiceID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(invoice)
			return
		}
	}

	http.Error(w, "Invoice not found", http.StatusNotFound)
}

func getInvoicesByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	log.Printf("[BILLING SERVICE] GET /invoices/user?user_id=%s\n", userID)

	var userInvoices []Invoice
	for _, invoice := range invoices {
		if invoice.UserID == userID {
			userInvoices = append(userInvoices, invoice)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInvoices)
}

func main() {
	http.HandleFunc("/invoices", getInvoices)
	http.HandleFunc("/invoice", getInvoiceByID)
	http.HandleFunc("/invoices/user", getInvoicesByUser)

	port := ":8083"
	log.Printf("[BILLING SERVICE] Started on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
