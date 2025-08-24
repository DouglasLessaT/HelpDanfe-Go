# Guia de Deploy - HelpDanfe-Go

## Pré-requisitos

- Go 1.21+
- PostgreSQL 12+
- Redis (opcional, para cache)
- Certificado digital A1 ou A3
- APIs bancárias configuradas

## Configuração de Produção

### 1. Variáveis de Ambiente

Crie um arquivo `.env` com as seguintes configurações:

```bash
# Configurações do Servidor
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ENVIRONMENT=production

# Configurações do Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=helpdanfe_user
DB_PASSWORD=sua_senha_segura
DB_NAME=helpdanfe_prod
DB_SSL_MODE=require

# Configurações do Certificado Digital
CERT_PATH=/path/to/your/certificate.p12
CERT_PASSWORD=sua_senha_certificado

# Configurações da SEFAZ
SEFAZ_AMBIENTE=producao
SEFAZ_UF=SP
SEFAZ_TIMEOUT=30

# Configurações das APIs Bancárias
ITAÚ_API_URL=https://api.itau.com.br
ITAÚ_CLIENT_ID=seu_client_id
ITAÚ_CLIENT_SECRET=seu_client_secret

BRADESCO_API_URL=https://api.bradesco.com.br
BRADESCO_CLIENT_ID=seu_client_id
BRADESCO_CLIENT_SECRET=seu_client_secret

OPEN_BANKING_URL=https://api.openbanking.com.br
OPEN_BANKING_CLIENT_ID=seu_client_id
OPEN_BANKING_CLIENT_SECRET=seu_client_secret

# Configurações de Log
LOG_LEVEL=info
LOG_FILE=/var/log/helpdanfe/app.log

# Configurações de Cache
CACHE_TTL=3600
REDIS_URL=redis://localhost:6379
```

### 2. Configuração do Banco de Dados

```sql
-- Criar usuário e banco
CREATE USER helpdanfe_user WITH PASSWORD 'sua_senha_segura';
CREATE DATABASE helpdanfe_prod OWNER helpdanfe_user;

-- Conceder permissões
GRANT ALL PRIVILEGES ON DATABASE helpdanfe_prod TO helpdanfe_user;
```

### 3. Configuração do Sistema

```bash
# Criar usuário do sistema
sudo useradd -r -s /bin/false helpdanfe

# Criar diretórios
sudo mkdir -p /opt/helpdanfe
sudo mkdir -p /var/log/helpdanfe
sudo mkdir -p /etc/helpdanfe/certs

# Definir permissões
sudo chown -R helpdanfe:helpdanfe /opt/helpdanfe
sudo chown -R helpdanfe:helpdanfe /var/log/helpdanfe
sudo chown -R helpdanfe:helpdanfe /etc/helpdanfe
```

## Deploy com Docker

### 1. Usando Docker Compose

```bash
# Clonar o repositório
git clone https://github.com/seu-usuario/helpdanfe-go.git
cd helpdanfe-go

# Configurar variáveis de ambiente
cp env.example .env
# Editar .env com suas configurações

# Executar com Docker Compose
docker-compose up -d
```

### 2. Usando Docker diretamente

```bash
# Construir imagem
docker build -t helpdanfe-go .

# Executar container
docker run -d \
  --name helpdanfe-go \
  -p 8080:8080 \
  -v /path/to/certs:/app/certs \
  -v /path/to/logs:/app/logs \
  --env-file .env \
  helpdanfe-go
```

## Deploy Manual

### 1. Compilar a aplicação

```bash
# Instalar dependências
go mod tidy

# Compilar para produção
GOOS=linux GOARCH=amd64 go build -o helpdanfe-go cmd/server/main.go
```

### 2. Configurar systemd

Crie o arquivo `/etc/systemd/system/helpdanfe.service`:

```ini
[Unit]
Description=HelpDanfe-Go API
After=network.target postgresql.service

[Service]
Type=simple
User=helpdanfe
Group=helpdanfe
WorkingDirectory=/opt/helpdanfe
ExecStart=/opt/helpdanfe/helpdanfe-go
Restart=always
RestartSec=5
Environment=ENVIRONMENT=production

[Install]
WantedBy=multi-user.target
```

### 3. Iniciar o serviço

```bash
# Recarregar systemd
sudo systemctl daemon-reload

# Habilitar serviço
sudo systemctl enable helpdanfe

# Iniciar serviço
sudo systemctl start helpdanfe

# Verificar status
sudo systemctl status helpdanfe
```

## Configuração de Proxy Reverso (Nginx)

### 1. Instalar Nginx

```bash
sudo apt update
sudo apt install nginx
```

### 2. Configurar site

Crie o arquivo `/etc/nginx/sites-available/helpdanfe`:

```nginx
server {
    listen 80;
    server_name api.helpdanfe.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3. Habilitar site

```bash
sudo ln -s /etc/nginx/sites-available/helpdanfe /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## Configuração de SSL (Let's Encrypt)

### 1. Instalar Certbot

```bash
sudo apt install certbot python3-certbot-nginx
```

### 2. Obter certificado

```bash
sudo certbot --nginx -d api.helpdanfe.com
```

## Monitoramento

### 1. Logs

```bash
# Ver logs da aplicação
sudo journalctl -u helpdanfe -f

# Ver logs do Nginx
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### 2. Métricas

Configure o endpoint `/api/v1/health` para monitoramento:

```bash
# Verificar saúde da aplicação
curl -f http://localhost:8080/api/v1/health
```

### 3. Backup

```bash
# Backup do banco de dados
pg_dump -h localhost -U helpdanfe_user helpdanfe_prod > backup_$(date +%Y%m%d_%H%M%S).sql

# Backup dos certificados
sudo tar -czf certs_backup_$(date +%Y%m%d_%H%M%S).tar.gz /etc/helpdanfe/certs/
```

## Segurança

### 1. Firewall

```bash
# Configurar UFW
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
```

### 2. Certificados Digitais

- Mantenha os certificados em local seguro
- Use permissões restritas (600)
- Faça backup regular dos certificados
- Monitore a validade dos certificados

### 3. APIs Bancárias

- Use HTTPS para todas as comunicações
- Implemente rate limiting
- Monitore logs de acesso
- Use credenciais seguras

## Troubleshooting

### Problemas Comuns

1. **Erro de conexão com banco**
   - Verificar se PostgreSQL está rodando
   - Verificar credenciais no .env
   - Verificar se o banco existe

2. **Erro de certificado digital**
   - Verificar se o arquivo existe
   - Verificar permissões
   - Verificar senha do certificado

3. **Erro de APIs bancárias**
   - Verificar credenciais
   - Verificar conectividade de rede
   - Verificar logs da aplicação

### Logs de Debug

Para habilitar logs detalhados, configure:

```bash
LOG_LEVEL=debug
```

## Suporte

Para suporte técnico:
- Email: suporte@helpdanfe.com
- Documentação: https://docs.helpdanfe.com
- Issues: https://github.com/seu-usuario/helpdanfe-go/issues
