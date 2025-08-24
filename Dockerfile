# Dockerfile para HelpDanfe-Go

# Estágio de build
FROM golang:1.21-alpine AS builder

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Copiar arquivos do frontend
COPY web/ ./web/

# Estágio de produção
FROM alpine:latest

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /root/

# Copiar binário do estágio de build
COPY --from=builder /app/main .

# Criar diretórios necessários
RUN mkdir -p /root/logs /root/certs

# Definir permissões
RUN chown -R appuser:appgroup /root/

# Mudar para usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]
