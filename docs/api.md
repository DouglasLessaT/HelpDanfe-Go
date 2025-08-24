# Documentação da API HelpDanfe-Go

## Visão Geral

A API HelpDanfe-Go é uma solução completa para consulta de NFes e localização de boletos bancários. Ela permite consultar NFes na SEFAZ, baixar XMLs, gerar DANFEs em PDF e localizar boletos vinculados.

## Base URL

```
http://localhost:8080/api/v1
```

## Autenticação

Atualmente a API não requer autenticação, mas é recomendado implementar JWT ou API Key para produção.

## Endpoints

### 1. Health Check

**GET** `/health`

Verifica o status da API.

**Resposta:**
```json
{
  "status": "ok",
  "message": "API funcionando normalmente",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0"
}
```

### 2. Consultar NFe

**POST** `/nfe/consultar`

Consulta uma NFe na SEFAZ usando a chave de acesso.

**Corpo da requisição:**
```json
{
  "chave_acesso": "12345678901234567890123456789012345678901234",
  "certificado": "caminho/para/certificado.p12",
  "senha": "senha_certificado"
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "NFe consultada com sucesso",
  "data": {
    "id": 1,
    "chave_acesso": "12345678901234567890123456789012345678901234",
    "numero": "123456",
    "serie": "1",
    "data_emissao": "2024-01-01T10:00:00Z",
    "status": "AUTORIZADA",
    "emitente_cnpj": "12345678000123",
    "emitente_nome": "EMPRESA EXEMPLO LTDA",
    "destinatario_cnpj": "98765432000198",
    "destinatario_nome": "CLIENTE EXEMPLO LTDA",
    "valor_total": 1000.00,
    "duplicatas": [
      {
        "numero": "001",
        "vencimento": "2024-02-01T00:00:00Z",
        "valor": 1000.00
      }
    ]
  }
}
```

### 3. Baixar XML da NFe

**GET** `/nfe/{chave}/xml`

Baixa o XML completo da NFe.

**Parâmetros:**
- `chave` (string, obrigatório): Chave de acesso da NFe (44 dígitos)

**Resposta:** Arquivo XML da NFe

### 4. Gerar PDF da NFe

**GET** `/nfe/{chave}/pdf`

Gera o DANFE em PDF.

**Parâmetros:**
- `chave` (string, obrigatório): Chave de acesso da NFe (44 dígitos)

**Resposta:** Arquivo PDF do DANFE

### 5. Consultar Boletos da NFe

**GET** `/nfe/{chave}/boletos`

Localiza boletos vinculados a uma NFe.

**Parâmetros:**
- `chave` (string, obrigatório): Chave de acesso da NFe (44 dígitos)

**Resposta:**
```json
{
  "success": true,
  "message": "Boletos consultados com sucesso",
  "data": {
    "nfe": {
      "chave_acesso": "12345678901234567890123456789012345678901234",
      "numero": "123456"
    },
    "boletos": [
      {
        "id": 1,
        "banco": "001",
        "numero": "BOL001",
        "codigo_barras": "00193373700000001000500940144816060680935031",
        "linha_digitavel": "00190.00009 04441.601448 60606.809350 3 37370000000100",
        "valor": 1000.00,
        "vencimento": "2024-02-01T00:00:00Z",
        "status": "ABERTO"
      }
    ]
  }
}
```

### 6. Consultar Boleto Específico

**GET** `/boletos/{codigo}`

Consulta um boleto específico.

**Parâmetros:**
- `codigo` (string, obrigatório): Código ou número do boleto

**Resposta:**
```json
{
  "success": true,
  "message": "Boleto consultado com sucesso",
  "data": {
    "id": 1,
    "banco": "001",
    "numero": "BOL001",
    "codigo_barras": "00193373700000001000500940144816060680935031",
    "linha_digitavel": "00190.00009 04441.601448 60606.809350 3 37370000000100",
    "valor": 1000.00,
    "vencimento": "2024-02-01T00:00:00Z",
    "status": "ABERTO"
  }
}
```

### 7. Consultar Múltiplos Boletos

**POST** `/boletos/consultar`

Consulta múltiplos boletos de uma vez.

**Corpo da requisição:**
```json
{
  "codigos": ["BOL001", "BOL002", "BOL003"]
}
```

**Resposta:**
```json
{
  "success": true,
  "message": "Boletos consultados com sucesso",
  "data": [
    {
      "id": 1,
      "banco": "001",
      "numero": "BOL001",
      "valor": 1000.00,
      "status": "ABERTO"
    },
    {
      "id": 2,
      "banco": "341",
      "numero": "BOL002",
      "valor": 500.00,
      "status": "PAGO"
    }
  ]
}
```

## Códigos de Status HTTP

- `200` - Sucesso
- `400` - Requisição inválida
- `404` - Recurso não encontrado
- `500` - Erro interno do servidor

## Exemplos de Uso

### cURL

```bash
# Consultar NFe
curl -X POST http://localhost:8080/api/v1/nfe/consultar \
  -H "Content-Type: application/json" \
  -d '{"chave_acesso": "12345678901234567890123456789012345678901234"}'

# Baixar XML
curl -O http://localhost:8080/api/v1/nfe/12345678901234567890123456789012345678901234/xml

# Gerar PDF
curl -O http://localhost:8080/api/v1/nfe/12345678901234567890123456789012345678901234/pdf

# Consultar boletos
curl http://localhost:8080/api/v1/nfe/12345678901234567890123456789012345678901234/boletos
```

### JavaScript (Fetch)

```javascript
// Consultar NFe
const response = await fetch('/api/v1/nfe/consultar', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    chave_acesso: '12345678901234567890123456789012345678901234'
  })
});

const data = await response.json();
console.log(data);
```

## Limitações

- Chave de acesso deve ter exatamente 44 dígitos
- Certificado digital deve estar no formato .p12 ou .pfx
- APIs bancárias requerem configuração prévia
- Rate limiting pode ser aplicado em produção

## Suporte

Para suporte técnico, entre em contato através do email: suporte@helpdanfe.com
