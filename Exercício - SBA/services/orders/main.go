package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Order struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Total    float64 `json:"total"`
	Status   string  `json:"status"`
}

var orders = []Order{
	{ID: "1001", UserID: "1", Product: "Notebook", Quantity: 1, Total: 3500.00, Status: "delivered"},
	{ID: "1002", UserID: "2", Product: "Mouse", Quantity: 2, Total: 100.00, Status: "processing"},
	{ID: "1003", UserID: "1", Product: "Keyboard", Quantity: 1, Total: 250.00, Status: "shipped"},
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	log.Println("[ORDERS SERVICE] GET /orders")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	log.Printf("[ORDERS SERVICE] GET /order?id=%s\n", orderID)

	for _, order := range orders {
		if order.ID == orderID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(order)
			return
		}
	}

	http.Error(w, "Order not found", http.StatusNotFound)
}

func getOrdersByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	log.Printf("[ORDERS SERVICE] GET /orders/user?user_id=%s\n", userID)

	var userOrders []Order
	for _, order := range orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userOrders)
}

func main() {
	http.HandleFunc("/orders", getOrders)
	http.HandleFunc("/order", getOrderByID)
	http.HandleFunc("/orders/user", getOrdersByUser)

	port := ":8082"
	log.Printf("[ORDERS SERVICE] Started on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
