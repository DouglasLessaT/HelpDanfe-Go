# Makefile para HelpDanfe-Go

# Variáveis
BINARY_NAME=helpdanfe-go
BUILD_DIR=build
MAIN_FILE=cmd/server/main.go

# Comandos principais
.PHONY: all build clean test run dev docker-build docker-run help

all: clean build

# Compilar a aplicação
build:
	@echo "Compilando aplicação..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Aplicação compilada em $(BUILD_DIR)/$(BINARY_NAME)"

# Limpar arquivos de build
clean:
	@echo "Limpando arquivos de build..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Executar testes
test:
	@echo "Executando testes..."
	@go test -v ./...

# Executar testes com coverage
test-coverage:
	@echo "Executando testes com coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de coverage gerado em coverage.html"

# Executar a aplicação
run: build
	@echo "Executando aplicação..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Executar em modo desenvolvimento
dev:
	@echo "Executando em modo desenvolvimento..."
	@go run $(MAIN_FILE)

# Instalar dependências
deps:
	@echo "Instalando dependências..."
	@go mod tidy
	@go mod download

# Verificar dependências
deps-check:
	@echo "Verificando dependências..."
	@go mod verify

# Executar linter
lint:
	@echo "Executando linter..."
	@golangci-lint run

# Formatar código
fmt:
	@echo "Formatando código..."
	@go fmt ./...

# Gerar documentação
docs:
	@echo "Gerando documentação..."
	@swag init -g $(MAIN_FILE) -o docs

# Criar diretórios necessários
setup:
	@echo "Criando diretórios necessários..."
	@mkdir -p logs
	@mkdir -p certs
	@mkdir -p build

# Instalar ferramentas de desenvolvimento
install-tools:
	@echo "Instalando ferramentas de desenvolvimento..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# Docker
docker-build:
	@echo "Construindo imagem Docker..."
	@docker build -t $(BINARY_NAME) .

docker-run:
	@echo "Executando container Docker..."
	@docker run -p 8080:8080 $(BINARY_NAME)

# Migração do banco de dados
migrate:
	@echo "Executando migrações..."
	@go run $(MAIN_FILE) migrate

# Seed do banco de dados
seed:
	@echo "Executando seed do banco..."
	@go run $(MAIN_FILE) seed

# Backup do banco de dados
backup:
	@echo "Fazendo backup do banco..."
	@pg_dump -h localhost -U postgres helpdanfe > backup_$(shell date +%Y%m%d_%H%M%S).sql

# Restaurar backup do banco de dados
restore:
	@echo "Restaurando backup do banco..."
	@psql -h localhost -U postgres helpdanfe < $(BACKUP_FILE)

# Verificar saúde da aplicação
health:
	@echo "Verificando saúde da aplicação..."
	@curl -f http://localhost:8080/api/v1/health || echo "Aplicação não está respondendo"

# Testar frontend
frontend:
	@echo "Testando frontend..."
	@./scripts/test-frontend.sh

# Monitorar logs
logs:
	@echo "Monitorando logs..."
	@tail -f logs/app.log

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  build        - Compilar a aplicação"
	@echo "  clean        - Limpar arquivos de build"
	@echo "  test         - Executar testes"
	@echo "  run          - Executar a aplicação"
	@echo "  dev          - Executar em modo desenvolvimento"
	@echo "  deps         - Instalar dependências"
	@echo "  lint         - Executar linter"
	@echo "  fmt          - Formatar código"
	@echo "  setup        - Criar diretórios necessários"
	@echo "  docker-build - Construir imagem Docker"
	@echo "  docker-run   - Executar container Docker"
	@echo "  migrate      - Executar migrações"
	@echo "  health       - Verificar saúde da aplicação"
	@echo "  frontend     - Testar frontend"
	@echo "  help         - Mostrar esta ajuda"
