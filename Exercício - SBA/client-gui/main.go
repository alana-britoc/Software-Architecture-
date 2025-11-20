package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const gatewayURL = "http://localhost:8090"

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Total    float64 `json:"total"`
	Status   string  `json:"status"`
}

type Invoice struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	OrderID     string  `json:"order_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	PaymentDate string  `json:"payment_date"`
}

func makeRequest(endpoint string) (string, error) {
	url := gatewayURL + endpoint
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("Error: %d %s\n%s", resp.StatusCode, http.StatusText(resp.StatusCode), string(body)), nil
	}

	// Pretty print JSON
	var prettyJSON interface{}
	if err := json.Unmarshal(body, &prettyJSON); err == nil {
		formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
		return string(formatted), nil
	}

	return string(body), nil
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Sistema SBA - API Gateway Dashboard")
	myWindow.Resize(fyne.NewSize(900, 650))

	// Status label
	statusLabel := widget.NewLabel("Sistema pronto. Selecione uma operação.")
	statusLabel.Wrapping = fyne.TextWrapWord

	// Output area
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Resposta do servidor aparecerá aqui...")
	output.Wrapping = fyne.TextWrapWord
	output.Disable()

	// Request info
	requestInfo := widget.NewLabel("")
	requestInfo.Wrapping = fyne.TextWrapWord

	// Loading indicator
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	// Function to make API calls
	makeAPICall := func(endpoint, description string) {
		statusLabel.SetText(fmt.Sprintf("🔄 %s...", description))
		requestInfo.SetText(fmt.Sprintf("Endpoint: %s%s", gatewayURL, endpoint))
		progressBar.Show()
		output.SetText("Carregando...")

		go func() {
			start := time.Now()
			result, err := makeRequest(endpoint)
			duration := time.Since(start)

			if err != nil {
				output.SetText(fmt.Sprintf("❌ Erro: %v\n\nVerifique se o Gateway e os serviços estão rodando.", err))
				statusLabel.SetText("❌ Erro na requisição")
			} else {
				output.SetText(result)
				statusLabel.SetText(fmt.Sprintf("✅ %s concluído em %v", description, duration.Round(time.Millisecond)))
			}
			progressBar.Hide()
		}()
	}

	// === USERS SECTION ===
	usersLabel := widget.NewLabelWithStyle("👥 USUÁRIOS", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	btnListUsers := widget.NewButton("📋 Listar Todos os Usuários", func() {
		makeAPICall("/api/users", "Listando usuários")
	})
	btnListUsers.Importance = widget.HighImportance

	btnUser1 := widget.NewButton("👤 Usuário ID: 1", func() {
		makeAPICall("/api/user?id=1", "Buscando usuário 1")
	})

	btnUser2 := widget.NewButton("👤 Usuário ID: 2", func() {
		makeAPICall("/api/user?id=2", "Buscando usuário 2")
	})

	usersBox := container.NewVBox(
		usersLabel,
		btnListUsers,
		container.NewGridWithColumns(2, btnUser1, btnUser2),
	)

	// === ORDERS SECTION ===
	ordersLabel := widget.NewLabelWithStyle("📦 PEDIDOS", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	btnListOrders := widget.NewButton("📋 Listar Todos os Pedidos", func() {
		makeAPICall("/api/orders", "Listando pedidos")
	})
	btnListOrders.Importance = widget.HighImportance

	btnOrder1001 := widget.NewButton("📦 Pedido ID: 1001", func() {
		makeAPICall("/api/order?id=1001", "Buscando pedido 1001")
	})

	btnOrdersUser1 := widget.NewButton("👤 Pedidos do Usuário 1", func() {
		makeAPICall("/api/orders/user?user_id=1", "Buscando pedidos do usuário 1")
	})

	ordersBox := container.NewVBox(
		ordersLabel,
		btnListOrders,
		container.NewGridWithColumns(2, btnOrder1001, btnOrdersUser1),
	)

	// === BILLING SECTION ===
	billingLabel := widget.NewLabelWithStyle("💰 FATURAMENTO", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	btnListInvoices := widget.NewButton("📋 Listar Todas as Faturas", func() {
		makeAPICall("/api/invoices", "Listando faturas")
	})
	btnListInvoices.Importance = widget.HighImportance

	btnInvoice001 := widget.NewButton("💳 Fatura INV-001", func() {
		makeAPICall("/api/invoice?id=INV-001", "Buscando fatura INV-001")
	})

	btnInvoicesUser1 := widget.NewButton("👤 Faturas do Usuário 1", func() {
		makeAPICall("/api/invoices/user?user_id=1", "Buscando faturas do usuário 1")
	})

	billingBox := container.NewVBox(
		billingLabel,
		btnListInvoices,
		container.NewGridWithColumns(2, btnInvoice001, btnInvoicesUser1),
	)

	// === TEST SECTION ===
	testLabel := widget.NewLabelWithStyle("🧪 TESTES", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	btnHealth := widget.NewButton("💚 Health Check", func() {
		makeAPICall("/health", "Verificando saúde do Gateway")
	})

	btnInvalid := widget.NewButton("❌ Rota Inválida (Teste de Erro)", func() {
		makeAPICall("/api/invalid", "Testando rota inválida")
	})

	testBox := container.NewVBox(
		testLabel,
		container.NewGridWithColumns(2, btnHealth, btnInvalid),
	)

	// === DEMO SECTION ===
	demoLabel := widget.NewLabelWithStyle("🎬 DEMONSTRAÇÃO", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	btnRunAll := widget.NewButton("▶️  Executar Demonstração Completa", func() {
		statusLabel.SetText("🎬 Executando demonstração...")

		endpoints := []struct {
			path string
			desc string
		}{
			{"/api/users", "Usuários"},
			{"/api/orders", "Pedidos"},
			{"/api/invoices", "Faturas"},
			{"/api/user?id=1", "Detalhes do Usuário 1"},
			{"/api/orders/user?user_id=1", "Pedidos do Usuário 1"},
			{"/api/invoices/user?user_id=1", "Faturas do Usuário 1"},
		}

		go func() {
			for i, ep := range endpoints {
				statusLabel.SetText(fmt.Sprintf("🎬 Demo [%d/%d]: %s", i+1, len(endpoints), ep.desc))
				makeAPICall(ep.path, ep.desc)
				time.Sleep(2 * time.Second)
			}
			statusLabel.SetText("✅ Demonstração completa!")
		}()
	})
	btnRunAll.Importance = widget.WarningImportance

	demoBox := container.NewVBox(
		demoLabel,
		btnRunAll,
	)

	// Left panel with all buttons
	leftPanel := container.NewVBox(
		widget.NewLabelWithStyle("🌐 API Gateway Dashboard", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		usersBox,
		widget.NewSeparator(),
		ordersBox,
		widget.NewSeparator(),
		billingBox,
		widget.NewSeparator(),
		testBox,
		widget.NewSeparator(),
		demoBox,
		layout.NewSpacer(),
	)

	// Right panel with output
	rightPanel := container.NewBorder(
		container.NewVBox(
			statusLabel,
			requestInfo,
			progressBar,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewScroll(output),
	)

	// Main split layout
	split := container.NewHSplit(
		container.NewScroll(leftPanel),
		rightPanel,
	)
	split.Offset = 0.35

	// Info footer
	footer := widget.NewLabel("💡 Gateway: http://localhost:8080 | Users: 8081 | Orders: 8082 | Billing: 8083")
	footer.Alignment = fyne.TextAlignCenter

	// Main container
	content := container.NewBorder(
		nil,
		footer,
		nil,
		nil,
		split,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
