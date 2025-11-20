package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const gatewayURL = "http://localhost:8090"

func makeRequest(endpoint string, description string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("UI REQUEST: %s\n", description)
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Println(strings.Repeat("=", 60))

	url := gatewayURL + endpoint
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	fmt.Printf("Response:\n%s\n", string(body))
	fmt.Println(strings.Repeat("=", 60))

	time.Sleep(1 * time.Second) // Pause between requests
}

func main() {
	fmt.Println("\n")
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         UI CLIENT - Simulador de Interface              ║")
	fmt.Println("║         Padrão API Gateway - Arquitetura SBA            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println("\nAguardando serviços iniciarem...")
	time.Sleep(2 * time.Second)

	// Simulate UI requests through the Gateway

	// 1. Get all users
	makeRequest("/api/users", "Listar todos os usuários")

	// 2. Get specific user
	makeRequest("/api/user?id=1", "Obter detalhes do usuário ID 1")

	// 3. Get all orders
	makeRequest("/api/orders", "Listar todos os pedidos")

	// 4. Get specific order
	makeRequest("/api/order?id=1001", "Obter detalhes do pedido ID 1001")

	// 5. Get orders by user
	makeRequest("/api/orders/user?user_id=1", "Listar pedidos do usuário ID 1")

	// 6. Get all invoices
	makeRequest("/api/invoices", "Listar todas as faturas")

	// 7. Get specific invoice
	makeRequest("/api/invoice?id=INV-001", "Obter detalhes da fatura INV-001")

	// 8. Get invoices by user
	makeRequest("/api/invoices/user?user_id=1", "Listar faturas do usuário ID 1")

	// 9. Test invalid route
	makeRequest("/api/invalid", "Testar rota inválida (erro esperado)")

	fmt.Println("\n")
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║            Simulação Concluída!                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println("")
}
