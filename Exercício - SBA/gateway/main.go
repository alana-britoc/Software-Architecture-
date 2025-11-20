package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	usersServiceURL   = "http://localhost:8081"
	ordersServiceURL  = "http://localhost:8082"
	billingServiceURL = "http://localhost:8083"
)

// Gateway routes incoming requests to the appropriate microservice
func gatewayHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	query := r.URL.RawQuery

	var targetURL string
	var serviceName string

	// Route determination based on path
	switch {
	case strings.HasPrefix(path, "/api/users"):
		targetURL = usersServiceURL + strings.TrimPrefix(path, "/api")
		serviceName = "USERS"
	case strings.HasPrefix(path, "/api/user"):
		targetURL = usersServiceURL + strings.TrimPrefix(path, "/api")
		serviceName = "USERS"
	case strings.HasPrefix(path, "/api/orders"):
		targetURL = ordersServiceURL + strings.TrimPrefix(path, "/api")
		serviceName = "ORDERS"
	case strings.HasPrefix(path, "/api/order"):
		targetURL = ordersServiceURL + strings.TrimPrefix(path, "/api")
		serviceName = "ORDERS"
	case strings.HasPrefix(path, "/api/invoices"):
		targetURL = billingServiceURL + strings.TrimPrefix(path, "/api")
		serviceName = "BILLING"
	case strings.HasPrefix(path, "/api/invoice"):
		targetURL = billingServiceURL + strings.TrimPrefix(path, "/api")
		serviceName = "BILLING"
	default:
		http.Error(w, "Service not found", http.StatusNotFound)
		log.Printf("[GATEWAY] Unknown route: %s\n", path)
		return
	}

	if query != "" {
		targetURL += "?" + query
	}

	log.Printf("[GATEWAY] Routing %s %s -> %s Service (%s)\n", r.Method, path, serviceName, targetURL)

	// Forward the request to the appropriate service
	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("[GATEWAY] Error forwarding to %s: %v\n", serviceName, err)
		http.Error(w, fmt.Sprintf("Error contacting %s service", serviceName), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	io.Copy(w, resp.Body)

	log.Printf("[GATEWAY] Response from %s Service: %d\n", serviceName, resp.StatusCode)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("[GATEWAY] Health check")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy","gateway":"running"}`))
}

func main() {
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/api/", gatewayHandler)

	port := ":8090"
	log.Println("=================================================")
	log.Println("         API GATEWAY - Sistema SBA")
	log.Println("=================================================")
	log.Printf("[GATEWAY] Started on port %s\n", port)
	log.Println("[GATEWAY] Routing:")
	log.Println("  - /api/users, /api/user -> Users Service (8081)")
	log.Println("  - /api/orders, /api/order -> Orders Service (8082)")
	log.Println("  - /api/invoices, /api/invoice -> Billing Service (8083)")
	log.Println("=================================================")
	log.Fatal(http.ListenAndServe(port, nil))
}
