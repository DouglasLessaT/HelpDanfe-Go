#!/bin/bash

# Script de gerenciamento para o projeto HelpDanfe-Go
# Uso: ./scripts/manage.sh [comando]

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
    echo -e "${BLUE}  HelpDanfe-Go - Gerenciador${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Função para verificar se o Docker está rodando
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker não está rodando. Inicie o Docker e tente novamente."
        exit 1
    fi
}

# Função para criar diretórios necessários
create_directories() {
    print_message "Criando diretórios necessários..."
    mkdir -p logs/nginx
    mkdir -p certs
    mkdir -p web/css web/js
    print_message "Diretórios criados com sucesso!"
}

# Função para iniciar os serviços
start() {
    print_header
    check_docker
    create_directories
    
    print_message "Iniciando serviços..."
    docker-compose up -d
    
    print_message "Aguardando serviços ficarem prontos..."
    sleep 10
    
    # Verificar status dos serviços
    docker-compose ps
    
    print_message "Serviços iniciados!"
    print_message "Frontend: http://localhost:3000"
    print_message "API: http://localhost:3000/api/v1/health"
    print_message "pgAdmin: http://localhost:5050 (admin@helpdanfe.com / admin)"
    print_message "PostgreSQL: localhost:55432"
    print_message "Redis: localhost:6379"
}

# Função para parar os serviços
stop() {
    print_header
    print_message "Parando serviços..."
    docker-compose down
    print_message "Serviços parados!"
}

# Função para reiniciar os serviços
restart() {
    print_header
    print_message "Reiniciando serviços..."
    docker-compose restart
    print_message "Serviços reiniciados!"
}

# Função para ver logs
logs() {
    print_header
    if [ -z "$2" ]; then
        print_message "Mostrando logs de todos os serviços..."
        docker-compose logs -f
    else
        print_message "Mostrando logs do serviço: $2"
        docker-compose logs -f "$2"
    fi
}

# Função para ver status
status() {
    print_header
    print_message "Status dos serviços:"
    docker-compose ps
}

# Função para limpar tudo
clean() {
    print_header
    print_warning "Esta operação irá remover todos os containers, volumes e imagens!"
    read -p "Tem certeza? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_message "Limpando tudo..."
        docker-compose down -v --rmi all
        docker system prune -f
        print_message "Limpeza concluída!"
    else
        print_message "Operação cancelada."
    fi
}

# Função para rebuild
rebuild() {
    print_header
    print_message "Fazendo rebuild dos serviços..."
    docker-compose down
    docker-compose build --no-cache
    docker-compose up -d
    print_message "Rebuild concluído!"
}

# Função para executar comandos no container da aplicação
exec_app() {
    print_message "Executando comando no container da aplicação..."
    docker-compose exec app "$@"
}

# Função para executar comandos no banco de dados
exec_db() {
    print_message "Executando comando no banco de dados..."
    docker-compose exec postgres "$@"
}

# Função para testar o banco de dados
test_db() {
    print_header
    print_message "Testando banco de dados..."
    if [ -f "./scripts/test-db.sh" ]; then
        chmod +x ./scripts/test-db.sh
        ./scripts/test-db.sh
    else
        print_error "Script de teste do banco não encontrado"
    fi
}

# Função para mostrar ajuda
show_help() {
    print_header
    echo "Uso: $0 [comando]"
    echo ""
    echo "Comandos disponíveis:"
    echo "  start       - Inicia todos os serviços"
    echo "  stop        - Para todos os serviços"
    echo "  restart     - Reinicia todos os serviços"
    echo "  status      - Mostra o status dos serviços"
    echo "  logs [svc]  - Mostra logs (todos ou de um serviço específico)"
    echo "  rebuild     - Faz rebuild de todos os serviços"
    echo "  clean       - Remove todos os containers, volumes e imagens"
    echo "  test-db     - Testa a conexão com o banco de dados"
    echo "  exec-app    - Executa comando no container da aplicação"
    echo "  exec-db     - Executa comando no container do banco"
    echo "  help        - Mostra esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0 start"
    echo "  $0 logs app"
    echo "  $0 test-db"
    echo "  $0 exec-app go test ./..."
    echo "  $0 exec-db psql -U postgres -d helpdanfe"
}

# Verificar se o comando foi fornecido
if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

# Processar comandos
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status)
        status
        ;;
    logs)
        logs "$@"
        ;;
    rebuild)
        rebuild
        ;;
    clean)
        clean
        ;;
    test-db)
        test_db
        ;;
    exec-app)
        shift
        exec_app "$@"
        ;;
    exec-db)
        shift
        exec_db "$@"
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Comando desconhecido: $1"
        echo ""
        show_help
        exit 1
        ;;
esac 