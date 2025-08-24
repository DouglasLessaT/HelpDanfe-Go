# Makefile para HelpDanfe-Go

.PHONY: help build run test clean docker-build docker-run docker-stop docker-logs docker-clean

# Variáveis
APP_NAME=helpdanfe-go
BINARY_NAME=main
DOCKER_IMAGE=helpdanfe-go
DOCKER_TAG=latest

# Comandos padrão
help: ## Mostra esta ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Desenvolvimento local
build: ## Compila a aplicação Go
	@echo "Compilando aplicação..."
	@go build -o bin/$(BINARY_NAME) cmd/server/main.go

run: build ## Compila e executa a aplicação localmente
	@echo "Executando aplicação..."
	@./bin/$(BINARY_NAME)

test: ## Executa os testes
	@echo "Executando testes..."
	@go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "Executando testes com cobertura..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado: coverage.html"

clean: ## Remove arquivos compilados
	@echo "Limpando arquivos..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Docker
docker-build: ## Constrói a imagem Docker
	@echo "Construindo imagem Docker..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: docker-build ## Constrói e executa com Docker Compose
	@echo "Iniciando serviços com Docker Compose..."
	@docker-compose up -d

docker-stop: ## Para os serviços Docker
	@echo "Parando serviços Docker..."
	@docker-compose down

docker-logs: ## Mostra logs dos serviços
	@echo "Mostrando logs..."
	@docker-compose logs -f

docker-clean: ## Remove containers, volumes e imagens
	@echo "Limpando Docker..."
	@docker-compose down -v --rmi all
	@docker system prune -f

# Desenvolvimento
dev-setup: ## Configura ambiente de desenvolvimento
	@echo "Configurando ambiente de desenvolvimento..."
	@mkdir -p logs/nginx certs web/css web/js
	@if [ ! -f .env ]; then cp env.example .env; echo "Arquivo .env criado. Configure as variáveis necessárias."; fi

dev-start: dev-setup ## Inicia ambiente de desenvolvimento
	@echo "Iniciando ambiente de desenvolvimento..."
	@./scripts/manage.sh start

dev-stop: ## Para ambiente de desenvolvimento
	@echo "Parando ambiente de desenvolvimento..."
	@./scripts/manage.sh stop

dev-logs: ## Mostra logs do ambiente de desenvolvimento
	@echo "Mostrando logs..."
	@./scripts/manage.sh logs

# Dependências
deps: ## Instala/atualiza dependências Go
	@echo "Instalando dependências..."
	@go mod tidy
	@go mod download

deps-update: ## Atualiza dependências para versões mais recentes
	@echo "Atualizando dependências..."
	@go get -u ./...
	@go mod tidy

# Linting e formatação
lint: ## Executa linter
	@echo "Executando linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint não encontrado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Formata código Go
	@echo "Formatando código..."
	@go fmt ./...

fmt-check: ## Verifica se o código está formatado
	@echo "Verificando formatação..."
	@if [ "$$(gofmt -l . | wc -l)" -gt 0 ]; then \
		echo "Código não está formatado. Execute 'make fmt'"; \
		exit 1; \
	fi

# Deploy
deploy-prepare: test lint fmt-check ## Prepara para deploy (testa, lint e formata)
	@echo "Preparação para deploy concluída!"

# Utilitários
status: ## Mostra status dos serviços
	@echo "Status dos serviços:"
	@./scripts/manage.sh status

shell: ## Acessa shell do container da aplicação
	@echo "Acessando shell do container..."
	@docker-compose exec app sh

db-shell: ## Acessa shell do container do banco
	@echo "Acessando shell do banco..."
	@docker-compose exec postgres psql -U postgres -d helpdanfe

# Comando padrão
.DEFAULT_GOAL := help
