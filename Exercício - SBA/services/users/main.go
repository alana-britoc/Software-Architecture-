package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: "1", Name: "João Silva", Email: "joao@example.com"},
	{ID: "2", Name: "Maria Santos", Email: "maria@example.com"},
	{ID: "3", Name: "Pedro Costa", Email: "pedro@example.com"},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("[USERS SERVICE] GET /users")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	log.Printf("[USERS SERVICE] GET /user?id=%s\n", userID)

	for _, user := range users {
		if user.ID == userID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/user", getUserByID)

	port := ":8081"
	log.Printf("[USERS SERVICE] Started on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
