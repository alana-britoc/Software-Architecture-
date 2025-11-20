# Cliente GUI - API Gateway Dashboard

Interface gráfica para interagir com o sistema API Gateway usando Fyne.

## Pré-requisitos

1. **Go 1.16+** instalado
2. **Serviços rodando**:
   - API Gateway (porta 8080)
   - Users Service (porta 8081)
   - Orders Service (porta 8082)
   - Billing Service (porta 8083)

## Como Executar

### Método 1: Script (Recomendado)

**Windows:**
```bash
cd ..
.\run-gui.bat
```
ou clique duas vezes em `run-gui.bat`

**Linux/Mac:**
```bash
cd ..
chmod +x run-gui.sh
./run-gui.sh
```

### Método 2: Direto

```bash
go run main.go
```

### Método 3: Compilar e Executar

**Nota**: No Windows, a compilação pode falhar devido a dependências OpenGL. Use `go run` ao invés.

**Linux/Mac:**
```bash
go build -o gui-client main.go
./gui-client
```

## Primeira Execução

Na primeira vez, o Go irá:
1. Baixar todas as dependências do Fyne (~50MB)
2. Compilar a aplicação
3. Abrir a janela gráfica

Isso pode levar 1-2 minutos. Execuções seguintes serão instantâneas.

## Funcionalidades

### Painel de Controle (Esquerdo)

#### 👥 Usuários
- Listar todos os usuários
- Buscar usuário por ID (1, 2)

#### 📦 Pedidos
- Listar todos os pedidos
- Buscar pedido específico
- Listar pedidos por usuário

#### 💰 Faturamento
- Listar todas as faturas
- Buscar fatura específica
- Listar faturas por usuário

#### 🧪 Testes
- Health check do Gateway
- Teste de rota inválida (erro 404)

#### 🎬 Demonstração
- Executa automaticamente 6 operações principais

### Painel de Resultados (Direito)

- Status da requisição
- Endpoint sendo acessado
- Barra de progresso
- Resposta JSON formatada
- Tempo de resposta

## Troubleshooting

### Erro: "connection refused"
**Causa**: Serviços não estão rodando
**Solução**: Execute `run.bat` para iniciar todos os serviços primeiro

### Erro: "go: cannot find main module"
**Causa**: Não está no diretório correto
**Solução**: `cd client-gui` antes de executar

### Janela não abre
**Causa**: Dependências faltando
**Solução**: Execute `go mod tidy` e tente novamente

### Build falha no Windows
**Causa**: Dependências OpenGL complexas
**Solução**: Use `go run main.go` ao invés de compilar

## Dependências

- `fyne.io/fyne/v2` - Framework GUI multiplataforma

Todas as dependências são instaladas automaticamente via `go mod`.

## Estrutura do Código

```go
main.go (280 linhas)
├── Imports e constantes
├── Structs (User, Order, Invoice)
├── makeRequest() - Faz chamadas HTTP
└── main()
    ├── Configuração da janela
    ├── Área de status e output
    ├── Botões por categoria
    └── Layout com split panel
```

## Customização

Para modificar o endpoint do Gateway, edite a constante:

```go
const gatewayURL = "http://localhost:8080"
```

## Screenshots Conceituais

```
┌─────────────────────────────────────────────┐
│ Sistema SBA - API Gateway Dashboard         │
├──────────┬──────────────────────────────────┤
│ USUÁRIOS │ Status: ✅ Listando usuários     │
│ ━━━━━━━━ │ Endpoint: /api/users             │
│ List All │ ━━━━━━━━━                        │
│ User 1   │ [                                │
│ User 2   │   {"id":"1","name":"João"...}    │
├──────────┤   {"id":"2","name":"Maria"...}   │
│ PEDIDOS  │   {"id":"3","name":"Pedro"...}   │
│ ━━━━━━━━ │ ]                                │
│ List All │                                  │
│ Order 01 │ Tempo: 15ms                      │
├──────────┤                                  │
│ BILLING  │                                  │
│ ━━━━━━━━ │                                  │
│ List All │                                  │
│ Invoice  │                                  │
├──────────┴──────────────────────────────────┤
│ Gateway: 8080 | Users: 8081 | Orders: 8082 │
└─────────────────────────────────────────────┘
```

## Compatibilidade

- ✅ Windows 10/11
- ✅ macOS 10.13+
- ✅ Linux (com X11 ou Wayland)

## Licença

Parte do projeto educacional API Gateway - Sistema SBA
