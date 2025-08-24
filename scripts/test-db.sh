#!/bin/bash

# Script para testar a conexão com o banco de dados
# Uso: ./scripts/test-db.sh

set -e

echo "🧪 Testando conexão com o banco de dados..."

# Verificar se o container está rodando
if ! docker-compose ps postgres | grep -q "Up"; then
    echo "❌ Container PostgreSQL não está rodando. Execute 'docker-compose up -d' primeiro."
    exit 1
fi

echo "✅ Container PostgreSQL está rodando"

# Testar conexão
echo "🔌 Testando conexão com o banco..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "SELECT version();"

echo "📊 Verificando tabelas criadas..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "\dt"

echo "🔍 Verificando estrutura da tabela NFe..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "\d nfe"

echo "🔍 Verificando estrutura da tabela Boletos..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "\d boletos"

echo "🔍 Verificando estrutura da tabela Duplicatas..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "\d duplicatas"

echo "🔍 Verificando estrutura da tabela Certificados..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "\d certificados"

echo "🔍 Verificando estrutura da tabela Logs..."
docker-compose exec postgres psql -U postgres -d helpdanfe -c "\d logs_consulta"

echo "✅ Teste do banco de dados concluído com sucesso!" 