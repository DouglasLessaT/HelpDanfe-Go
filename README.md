# HelpDanfe-Go

Sistema de consulta e geração de documentos fiscais eletrônicos (NFe) e boletos bancários desenvolvido em Go.

## 🚀 Funcionalidades

- **Consulta de NFe**: Busca e validação de notas fiscais eletrônicas
- **Geração de PDF**: Criação de DANFE (Documento Auxiliar da Nota Fiscal Eletrônica)
- **Consulta de Boletos**: Verificação de boletos bancários de múltiplos bancos
- **API REST**: Interface completa para integração com outros sistemas
- **Interface Web**: Frontend responsivo para consultas e visualização

## 🛠️ Tecnologias

- **Backend**: Go (Gin framework)
- **Banco de Dados**: PostgreSQL
- **Cache**: Redis
- **Proxy Reverso**: Nginx
- **Containerização**: Docker & Docker Compose
- **Frontend**: HTML, CSS, JavaScript

## 📋 Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.19+ (para desenvolvimento local)
- Git

## 🚀 Início Rápido

### 1. Clone o repositório
```bash
git clone <url-do-repositorio>
cd HelpDanfe-Go
```

### 2. Configure as variáveis de ambiente
```bash
# Copie o arquivo de exemplo
cp env.example .env

# Edite o arquivo .env com suas configurações
# (certificados, APIs bancárias, etc.)
```

### 3. Inicie os serviços
```bash
# Usando o script de gerenciamento (recomendado)
chmod +x scripts/manage.sh
./scripts/manage.sh start

# Ou usando docker-compose diretamente
docker-compose up -d
```

### 4. Acesse a aplicação
- **Frontend**: http://localhost:3000
- **API**: http://localhost:3000/api/v1/health
- **pgAdmin**: http://localhost:5050 (admin@helpdanfe.com / admin)

## 📚 Uso da API

### Endpoints Principais

#### NFe
- `POST /api/v1/nfe/consultar` - Consulta NFe por chave
- `GET /api/v1/nfe/{chave}/xml` - Download do XML da NFe
- `GET /api/v1/nfe/{chave}/pdf` - Geração do DANFE em PDF
- `GET /api/v1/nfe/{chave}/boletos` - Consulta boletos da NFe

#### Boletos
- `GET /api/v1/boletos/{codigo}` - Consulta boleto por código
- `POST /api/v1/boletos/consultar` - Consulta múltiplos boletos

#### Certificados
- `GET /api/v1/certificados/verificar` - Verifica certificados disponíveis
- `POST /api/v1/certificados/selecionar` - Seleciona certificado ativo

#### Health Check
- `GET /api/v1/health` - Status da aplicação

### Exemplo de Uso

```bash
# Consultar NFe
curl -X POST http://localhost:3000/api/v1/nfe/consultar \
  -H "Content-Type: application/json" \
  -d '{"chave": "12345678901234567890123456789012345678901234"}'

# Verificar saúde da API
curl http://localhost:3000/api/v1/health
```

## 🐳 Gerenciamento com Docker

### Script de Gerenciamento

O projeto inclui um script para facilitar o gerenciamento dos containers:

```bash
# Iniciar serviços
./scripts/manage.sh start

# Ver status
./scripts/manage.sh status

# Ver logs
./scripts/manage.sh logs app

# Parar serviços
./scripts/manage.sh stop

# Rebuild completo
./scripts/manage.sh rebuild

# Limpar tudo
./scripts/manage.sh clean
```

### Comandos Docker Compose

```bash
# Iniciar
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar
docker-compose down

# Rebuild
docker-compose build --no-cache
```

## 🔧 Configuração

### Variáveis de Ambiente

Principais variáveis que podem ser configuradas:

- `SERVER_PORT`: Porta do servidor Go (padrão: 8080)
- `DB_HOST`: Host do PostgreSQL
- `DB_PASSWORD`: Senha do PostgreSQL
- `CORS_ALLOWED_ORIGIN`: Origens permitidas para CORS
- `SEFAZ_AMBIENTE`: Ambiente da SEFAZ (homologacao/producao)
- `CERT_PATH`: Caminho para o certificado digital

### Banco de Dados

O sistema cria automaticamente as tabelas necessárias:

- `nfe`: Notas fiscais eletrônicas
- `boletos`: Boletos bancários
- `duplicatas`: Duplicatas das NFe
- `certificados`: Certificados digitais
- `logs_consulta`: Logs das consultas realizadas

## 🧪 Testes

```bash
# Executar testes no container
./scripts/manage.sh exec-app go test ./...

# Ou acessar o container e executar
docker-compose exec app go test ./...
```

## 📁 Estrutura do Projeto

```
HelpDanfe-Go/
├── cmd/server/          # Ponto de entrada da aplicação
├── internal/            # Código interno da aplicação
│   ├── config/         # Configurações
│   ├── database/       # Conexão com banco de dados
│   ├── handlers/       # Handlers HTTP
│   ├── middleware/     # Middlewares
│   ├── models/         # Modelos de dados
│   ├── services/       # Lógica de negócio
│   └── utils/          # Utilitários
├── web/                # Frontend
├── scripts/            # Scripts de automação
├── docker-compose.yml  # Configuração Docker
├── nginx.conf         # Configuração Nginx
└── README.md          # Este arquivo
```

## 🔒 Segurança

- **CORS configurável**: Permite definir origens específicas
- **Health checks**: Monitoramento de saúde dos serviços
- **Logs estruturados**: Rastreamento de todas as operações
- **Volumes read-only**: Arquivos estáticos e configurações protegidos

## 🐛 Solução de Problemas

### Problemas Comuns

1. **Porta já em uso**: Verifique se as portas 3000, 5050, 55432, 6379 estão livres
2. **Erro de permissão**: Execute `chmod +x scripts/manage.sh` no Linux/Mac
3. **Banco não conecta**: Aguarde o PostgreSQL inicializar (pode levar alguns segundos)

### Logs

```bash
# Ver logs da aplicação
./scripts/manage.sh logs app

# Ver logs do banco
./scripts/manage.sh logs postgres

# Ver logs do nginx
./scripts/manage.sh logs nginx
```

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença [MIT](LICENSE).

## 📞 Suporte

Para suporte, abra uma issue no repositório ou entre em contato com a equipe de desenvolvimento.

---

**Nota**: Este sistema é destinado para uso em ambiente de homologação. Para produção, configure adequadamente as variáveis de ambiente e certificados digitais.
