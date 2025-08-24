#!/bin/bash

# Script para testar o frontend HelpDanfe-Go
# Uso: ./scripts/test-frontend.sh

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_message() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  Teste do Frontend HelpDanfe${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Função para testar se o nginx está rodando
test_nginx() {
    print_message "Testando se o container nginx está rodando..."
    
    if ! docker-compose ps nginx | grep -q "Up"; then
        print_error "Container nginx não está rodando. Execute 'docker-compose up -d' primeiro."
        return 1
    fi
    
    print_message "✅ Container nginx está rodando"
    return 0
}

# Função para testar se o frontend está acessível
test_frontend_access() {
    print_message "Testando acesso ao frontend..."
    
    # Aguardar um pouco para o nginx inicializar
    sleep 3
    
    # Testar se a página principal carrega
    if curl -s -f http://localhost:3000/ > /dev/null; then
        print_message "✅ Frontend acessível em http://localhost:3000/"
    else
        print_error "❌ Frontend não está acessível em http://localhost:3000/"
        return 1
    fi
    
    return 0
}

# Função para testar arquivos estáticos
test_static_files() {
    print_message "Testando arquivos estáticos..."
    
    local base_url="http://localhost:3000"
    local files=("css/style.css" "js/app.js" "index.html")
    local all_ok=true
    
    for file in "${files[@]}"; do
        if curl -s -f "$base_url/$file" > /dev/null; then
            print_message "✅ $file - OK"
        else
            print_error "❌ $file - FALHOU"
            all_ok=false
        fi
    done
    
    if $all_ok; then
        print_message "✅ Todos os arquivos estáticos estão acessíveis"
        return 0
    else
        print_error "❌ Alguns arquivos estáticos falharam"
        return 1
    fi
}

# Função para testar API
test_api() {
    print_message "Testando API..."
    
    if curl -s -f http://localhost:3000/api/v1/health > /dev/null; then
        print_message "✅ API está acessível em /api/v1/health"
        return 0
    else
        print_warning "⚠️  API não está acessível (pode ser normal se a aplicação Go não estiver rodando)"
        return 0
    fi
}

# Função para verificar logs do nginx
check_nginx_logs() {
    print_message "Verificando logs do nginx..."
    
    echo "📋 Últimas 10 linhas dos logs de acesso:"
    docker-compose exec nginx tail -n 10 /var/log/nginx/access.log 2>/dev/null || echo "Logs não disponíveis"
    
    echo ""
    echo "📋 Últimas 10 linhas dos logs de erro:"
    docker-compose exec nginx tail -n 10 /var/log/nginx/error.log 2>/dev/null || echo "Logs não disponíveis"
}

# Função para testar estrutura dos arquivos no container
test_container_files() {
    print_message "Verificando estrutura dos arquivos no container..."
    
    echo "📁 Estrutura do diretório /usr/share/nginx/html:"
    docker-compose exec nginx ls -la /usr/share/nginx/html/
    
    echo ""
    echo "📄 Conteúdo do index.html (primeiras 5 linhas):"
    docker-compose exec nginx head -n 5 /usr/share/nginx/html/index.html
}

# Função principal
main() {
    print_header
    
    # Verificar se o docker-compose está rodando
    if ! test_nginx; then
        exit 1
    fi
    
    # Testar acesso ao frontend
    if ! test_frontend_access; then
        print_error "❌ Teste de acesso falhou"
        check_nginx_logs
        exit 1
    fi
    
    # Testar arquivos estáticos
    if ! test_static_files; then
        print_error "❌ Teste de arquivos estáticos falhou"
        check_nginx_logs
        exit 1
    fi
    
    # Testar API
    test_api
    
    # Verificar logs
    check_nginx_logs
    
    # Verificar estrutura dos arquivos
    test_container_files
    
    print_message "🎉 Todos os testes do frontend foram concluídos com sucesso!"
    print_message "🌐 Frontend disponível em: http://localhost:3000"
}

# Executar função principal
main "$@"
