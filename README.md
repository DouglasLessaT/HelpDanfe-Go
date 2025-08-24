# HelpDanfe-Go

API web em Go para consultar NFes e localizar boletos/duplicatas vinculados.

## 🚀 Funcionalidades

- **Consulta de NFe**: Consulta NFe na SEFAZ usando chave de acesso e certificado digital
- **Download de XML**: Baixa o XML completo da NFe autorizada
- **Geração de PDF**: Monta o DANFE em PDF a partir do XML
- **Localização de Boletos**: Consulta APIs bancárias para localizar duplicatas vinculadas à NFe
- **Código de Barras**: Exibe código de barras e valores dos boletos encontrados

## 📋 Pré-requisitos

- Go 1.21+
- PostgreSQL
- Certificado digital A1 ou A3 da empresa
- APIs bancárias configuradas

## 🏗️ Arquitetura

```
├── cmd/
│   └── server/          # Ponto de entrada da aplicação
├── internal/
│   ├── config/          # Configurações da aplicação
│   ├── database/        # Configuração do banco de dados
│   ├── handlers/        # Handlers HTTP
│   ├── middleware/      # Middlewares
│   ├── models/          # Modelos de dados
│   ├── services/        # Lógica de negócio
│   │   ├── nfe/         # Serviços relacionados à NFe
│   │   ├── bank/        # Serviços bancários
│   │   └── pdf/         # Geração de PDF
│   └── utils/           # Utilitários
├── pkg/                 # Pacotes públicos
└── docs/               # Documentação da API
```

## 🔧 Instalação

1. Clone o repositório:
```bash
git clone https://github.com/DouglasLessaT/helpdanfe-go.git
cd helpdanfe-go
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Configure as variáveis de ambiente:
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

4. Execute a aplicação:
```bash
go run cmd/server/main.go
```

## 📡 Endpoints da API

### NFe
- `POST /api/v1/nfe/consultar` - Consulta NFe na SEFAZ
- `GET /api/v1/nfe/{chave}/xml` - Baixa XML da NFe
- `GET /api/v1/nfe/{chave}/pdf` - Gera DANFE em PDF
- `GET /api/v1/nfe/{chave}/boletos` - Localiza boletos vinculados

### Boletos
- `GET /api/v1/boletos/{codigo}` - Consulta boleto específico
- `POST /api/v1/boletos/consultar` - Consulta múltiplos boletos

## 🔐 Configuração do Certificado Digital

1. Coloque seu certificado digital no diretório `certs/`
2. Configure o caminho no arquivo `.env`
3. Certifique-se de que o certificado está no formato correto (.p12 ou .pfx)

## 🏦 Integração Bancária

O sistema suporta integração com:
- Itaú API
- Bradesco API
- Open Banking Brasil
- APIs customizadas

## 📊 Banco de Dados

O sistema utiliza PostgreSQL para armazenar:
- Logs de consultas
- Cache de NFes
- Configurações de boletos
- Histórico de transações

## 🧪 Testes

```bash
go test ./...
```

## 📝 Licença

MIT License
