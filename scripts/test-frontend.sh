#!/bin/bash

# Script para testar o frontend HelpDanfe

echo "🚀 Testando o Frontend HelpDanfe"
echo "=================================="

# Verificar se os arquivos do frontend existem
echo "📁 Verificando arquivos do frontend..."

if [ ! -f "web/index.html" ]; then
    echo "❌ Erro: web/index.html não encontrado"
    exit 1
fi

if [ ! -f "web/css/style.css" ]; then
    echo "❌ Erro: web/css/style.css não encontrado"
    exit 1
fi

if [ ! -f "web/js/app.js" ]; then
    echo "❌ Erro: web/js/app.js não encontrado"
    exit 1
fi

echo "✅ Todos os arquivos do frontend encontrados"

# Verificar se a API está rodando
echo "🔍 Verificando se a API está rodando..."

if curl -s http://localhost:8080/api/v1/health > /dev/null; then
    echo "✅ API está rodando em http://localhost:8080"
else
    echo "⚠️  API não está rodando. Iniciando..."
    ./build/helpdanfe-go &
    API_PID=$!
    sleep 3
    
    if curl -s http://localhost:8080/api/v1/health > /dev/null; then
        echo "✅ API iniciada com sucesso"
    else
        echo "❌ Erro ao iniciar a API"
        exit 1
    fi
fi

# Testar endpoints da API
echo "🧪 Testando endpoints da API..."

# Health check
if curl -s http://localhost:8080/api/v1/health | grep -q "success"; then
    echo "✅ Health check funcionando"
else
    echo "❌ Health check falhou"
fi

# Teste de consulta de NFe (mock)
echo "📄 Testando consulta de NFe..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/nfe/consultar \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "chave_acesso=12345678901234567890123456789012345678901234")

if echo "$RESPONSE" | grep -q "success"; then
    echo "✅ Consulta de NFe funcionando"
else
    echo "⚠️  Consulta de NFe retornou erro (esperado para chave inválida)"
fi

# Testar acesso ao frontend
echo "🌐 Testando acesso ao frontend..."

if curl -s http://localhost:8080/ | grep -q "HelpDanfe"; then
    echo "✅ Frontend acessível"
else
    echo "❌ Erro ao acessar frontend"
fi

# Verificar arquivos estáticos
echo "📂 Verificando arquivos estáticos..."

if curl -s http://localhost:8080/css/style.css | grep -q "body"; then
    echo "✅ CSS carregando"
else
    echo "❌ Erro ao carregar CSS"
fi

if curl -s http://localhost:8080/js/app.js | grep -q "function"; then
    echo "✅ JavaScript carregando"
else
    echo "❌ Erro ao carregar JavaScript"
fi

echo ""
echo "🎉 Testes concluídos!"
echo ""
echo "📋 Resumo:"
echo "   - Frontend: http://localhost:8080"
echo "   - API: http://localhost:8080/api/v1"
echo "   - Health: http://localhost:8080/api/v1/health"
echo ""
echo "💡 Para testar manualmente:"
echo "   1. Abra http://localhost:8080 no navegador"
echo "   2. Digite uma chave de acesso válida"
echo "   3. Teste as funcionalidades do frontend"
echo ""

# Se a API foi iniciada pelo script, perguntar se quer parar
if [ ! -z "$API_PID" ]; then
    echo "🤔 Deseja parar a API? (s/n)"
    read -r response
    if [[ "$response" =~ ^[Ss]$ ]]; then
        kill $API_PID
        echo "🛑 API parada"
    else
        echo "▶️  API continua rodando (PID: $API_PID)"
    fi
fi
