package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"
)

const (
	gatewayURL = "http://localhost:8090"
	webPort    = ":3000"
)

// Proxy handler to avoid CORS issues
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get the endpoint from query parameter
	endpoint := r.URL.Query().Get("endpoint")
	if endpoint == "" {
		http.Error(w, "Missing endpoint parameter", http.StatusBadRequest)
		return
	}

	// Make request to gateway
	targetURL := gatewayURL + endpoint
	log.Printf("[WEB CLIENT] Proxying request to: %s\n", targetURL)

	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("[WEB CLIENT] Error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Erro ao conectar ao Gateway: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// Forward response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		log.Printf("Não foi possível abrir o navegador automaticamente. Acesse: %s\n", url)
	}
}

func main() {
	// Serve static HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlContent))
	})

	// API proxy endpoint
	http.HandleFunc("/api/proxy", proxyHandler)

	// Server info
	serverURL := fmt.Sprintf("http://localhost%s", webPort)

	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         API Gateway - Dashboard Web                     ║")
	fmt.Println("║         Sistema SBA                                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Printf("🌐 Dashboard disponível em: %s\n", serverURL)
	fmt.Println("📡 Gateway esperado em: http://localhost:8090")
	fmt.Println("")
	fmt.Println("Abrindo navegador...")
	fmt.Println("")
	fmt.Println("Pressione Ctrl+C para encerrar")
	fmt.Println("══════════════════════════════════════════════════════════")

	// Open browser after short delay
	go func() {
		// Wait a bit for server to start
		time.Sleep(1 * time.Second)
		openBrowser(serverURL)
	}()

	log.Fatal(http.ListenAndServe(webPort, nil))
}

const htmlContent = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Gateway Dashboard - Sistema SBA</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        .header {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
            margin-bottom: 30px;
            text-align: center;
        }

        .header h1 {
            color: #667eea;
            font-size: 32px;
            margin-bottom: 10px;
        }

        .header p {
            color: #666;
            font-size: 16px;
        }

        .main-grid {
            display: grid;
            grid-template-columns: 350px 1fr;
            gap: 20px;
        }

        .sidebar {
            background: white;
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
            height: fit-content;
        }

        .content {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
            min-height: 600px;
        }

        .section {
            margin-bottom: 25px;
        }

        .section-title {
            font-size: 14px;
            font-weight: 600;
            color: #667eea;
            margin-bottom: 12px;
            padding-bottom: 8px;
            border-bottom: 2px solid #f0f0f0;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .btn {
            width: 100%;
            padding: 12px 16px;
            margin-bottom: 8px;
            border: none;
            border-radius: 8px;
            background: #f8f9fa;
            color: #333;
            cursor: pointer;
            font-size: 14px;
            text-align: left;
            transition: all 0.3s ease;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .btn:hover {
            background: #667eea;
            color: white;
            transform: translateX(5px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }

        .btn-primary {
            background: #667eea;
            color: white;
            font-weight: 600;
        }

        .btn-primary:hover {
            background: #5568d3;
        }

        .btn-demo {
            background: #ff6b6b;
            color: white;
            font-weight: 600;
        }

        .btn-demo:hover {
            background: #ee5a52;
        }

        .status-bar {
            background: #f8f9fa;
            padding: 15px 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            border-left: 4px solid #667eea;
        }

        .status-text {
            font-size: 14px;
            color: #666;
            margin-bottom: 5px;
        }

        .endpoint-text {
            font-size: 12px;
            color: #999;
            font-family: 'Courier New', monospace;
        }

        .response-area {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            min-height: 400px;
            font-family: 'Courier New', monospace;
            font-size: 13px;
            overflow-x: auto;
            border: 1px solid #e0e0e0;
        }

        .loading {
            text-align: center;
            padding: 40px;
            color: #667eea;
        }

        .spinner {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #667eea;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 15px;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .footer {
            background: white;
            padding: 15px;
            border-radius: 15px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.1);
            margin-top: 20px;
            text-align: center;
            font-size: 13px;
            color: #666;
        }

        .success { color: #28a745; }
        .error { color: #dc3545; }

        pre {
            white-space: pre-wrap;
            word-wrap: break-word;
        }

        .icon {
            font-size: 16px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌐 API Gateway Dashboard</h1>
            <p>Sistema SBA - Service-Based Architecture</p>
        </div>

        <div class="main-grid">
            <div class="sidebar">
                <div class="section">
                    <div class="section-title"><span class="icon">👥</span> USUÁRIOS</div>
                    <button class="btn btn-primary" onclick="makeRequest('/api/users', 'Listando todos os usuários')">
                        📋 Listar Todos os Usuários
                    </button>
                    <button class="btn" onclick="makeRequest('/api/user?id=1', 'Buscando usuário ID 1')">
                        👤 Usuário ID: 1
                    </button>
                    <button class="btn" onclick="makeRequest('/api/user?id=2', 'Buscando usuário ID 2')">
                        👤 Usuário ID: 2
                    </button>
                </div>

                <div class="section">
                    <div class="section-title"><span class="icon">📦</span> PEDIDOS</div>
                    <button class="btn btn-primary" onclick="makeRequest('/api/orders', 'Listando todos os pedidos')">
                        📋 Listar Todos os Pedidos
                    </button>
                    <button class="btn" onclick="makeRequest('/api/order?id=1001', 'Buscando pedido 1001')">
                        📦 Pedido ID: 1001
                    </button>
                    <button class="btn" onclick="makeRequest('/api/orders/user?user_id=1', 'Pedidos do usuário 1')">
                        👤 Pedidos do Usuário 1
                    </button>
                </div>

                <div class="section">
                    <div class="section-title"><span class="icon">💰</span> FATURAMENTO</div>
                    <button class="btn btn-primary" onclick="makeRequest('/api/invoices', 'Listando todas as faturas')">
                        📋 Listar Todas as Faturas
                    </button>
                    <button class="btn" onclick="makeRequest('/api/invoice?id=INV-001', 'Buscando fatura INV-001')">
                        💳 Fatura INV-001
                    </button>
                    <button class="btn" onclick="makeRequest('/api/invoices/user?user_id=1', 'Faturas do usuário 1')">
                        👤 Faturas do Usuário 1
                    </button>
                </div>

                <div class="section">
                    <div class="section-title"><span class="icon">🧪</span> TESTES</div>
                    <button class="btn" onclick="makeRequest('/health', 'Health check do Gateway')">
                        💚 Health Check
                    </button>
                    <button class="btn" onclick="makeRequest('/api/invalid', 'Testando rota inválida')">
                        ❌ Rota Inválida (Erro)
                    </button>
                </div>

                <div class="section">
                    <div class="section-title"><span class="icon">🎬</span> DEMONSTRAÇÃO</div>
                    <button class="btn btn-demo" onclick="runDemo()">
                        ▶️  Executar Demo Completa
                    </button>
                </div>
            </div>

            <div class="content">
                <div class="status-bar">
                    <div class="status-text" id="status">✅ Sistema pronto. Selecione uma operação.</div>
                    <div class="endpoint-text" id="endpoint"></div>
                </div>

                <div class="response-area" id="response">
                    <div style="text-align: center; padding: 60px 20px; color: #999;">
                        <div style="font-size: 48px; margin-bottom: 20px;">🚀</div>
                        <div style="font-size: 18px; margin-bottom: 10px;">Bem-vindo ao Dashboard!</div>
                        <div style="font-size: 14px;">Clique em qualquer botão à esquerda para começar</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="footer">
            💡 Gateway: http://localhost:8090 | Users: 8081 | Orders: 8082 | Billing: 8083
        </div>
    </div>

    <script>
        const proxyURL = '/api/proxy';

        async function makeRequest(endpoint, description) {
            const statusEl = document.getElementById('status');
            const endpointEl = document.getElementById('endpoint');
            const responseEl = document.getElementById('response');

            // Update status
            statusEl.textContent = '\uD83D\uDD04 ' + description + '...';
            endpointEl.textContent = 'Endpoint: http://localhost:8090' + endpoint;

            // Show loading
            responseEl.innerHTML = '<div class="loading"><div class="spinner"></div>Carregando...</div>';

            try {
                const startTime = Date.now();
                const response = await fetch(proxyURL + '?endpoint=' + encodeURIComponent(endpoint));
                const duration = Date.now() - startTime;

                const text = await response.text();
                let data;

                try {
                    data = JSON.parse(text);
                } catch (e) {
                    data = { error: text };
                }

                if (response.ok) {
                    statusEl.textContent = '\u2705 ' + description + ' concluido em ' + duration + 'ms';
                    responseEl.innerHTML = '<pre>' + JSON.stringify(data, null, 2) + '</pre>';
                } else {
                    statusEl.textContent = '\u274C Erro na requisicao';
                    responseEl.innerHTML = '<pre class="error">' + JSON.stringify(data, null, 2) + '</pre>';
                }
            } catch (error) {
                statusEl.textContent = '\u274C Erro ao conectar ao servidor';
                responseEl.innerHTML = '<pre class="error">Erro: ' + error.message + '\n\nVerifique se:\n1. O Gateway esta rodando na porta 8090\n2. Os microsservicos estao ativos\n3. Nao ha bloqueio de firewall</pre>';
            }
        }

        async function runDemo() {
            const demos = [
                ['/api/users', 'Listando usuários'],
                ['/api/orders', 'Listando pedidos'],
                ['/api/invoices', 'Listando faturas'],
                ['/api/user?id=1', 'Detalhes do usuário 1'],
                ['/api/orders/user?user_id=1', 'Pedidos do usuário 1'],
                ['/api/invoices/user?user_id=1', 'Faturas do usuário 1']
            ];

            for (let i = 0; i < demos.length; i++) {
                const [endpoint, desc] = demos[i];
                document.getElementById('status').textContent = '\u{1F3AC} Demo [' + (i+1) + '/' + demos.length + ']: ' + desc;
                await makeRequest(endpoint, desc);
                await new Promise(resolve => setTimeout(resolve, 2000));
            }

            document.getElementById('status').textContent = '\u2705 Demonstracao completa!';
        }
    </script>
</body>
</html>
`
