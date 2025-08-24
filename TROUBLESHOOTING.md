# 🔧 Guia de Solução de Problemas - HelpDanfe-Go

## 🚨 Erro: "insufficient arguments" no banco de dados

### **Problema Identificado:**
O erro "insufficient arguments" está ocorrendo porque:
1. **Incompatibilidade entre nomes de tabelas**: Script SQL cria `nfe` mas GORM tenta acessar `nfes`
2. **Modelos faltando**: Faltavam modelos para `Certificado` e `LogConsulta`
3. **Estrutura de tabelas incompatível**: Campos não correspondiam entre SQL e Go

### **✅ Soluções Implementadas:**

#### 1. **Corrigidos os nomes das tabelas:**
- `nfes` → `nfe` (singular)
- `boletos` → `boletos` (já estava correto)
- `duplicatas` → `duplicatas` (já estava correto)

#### 2. **Criados modelos faltantes:**
- `internal/models/certificado.go`
- `internal/models/log_consulta.go`

#### 3. **Atualizada função de migração:**
- Adicionados logs detalhados
- Migração individual para cada modelo
- Tratamento de erros melhorado

#### 4. **Corrigido script SQL:**
- Campos alinhados com os modelos Go
- Adicionados campos `deleted_at` para soft delete
- Índices para soft delete

### **🔄 Como Aplicar as Correções:**

#### **Opção 1: Rebuild completo (Recomendado)**
```bash
# Parar serviços
./scripts/manage.sh stop

# Limpar tudo
./scripts/manage.sh clean

# Rebuild completo
./scripts/manage.sh rebuild

# Iniciar serviços
./scripts/manage.sh start
```

#### **Opção 2: Apenas rebuild da aplicação**
```bash
# Parar apenas a aplicação
docker-compose stop app

# Rebuild da aplicação
docker-compose build app

# Iniciar tudo
docker-compose up -d
```

### **🧪 Testar se Funcionou:**

#### **1. Verificar logs da aplicação:**
```bash
./scripts/manage.sh logs app
```

#### **2. Testar conexão com banco:**
```bash
./scripts/manage.sh test-db
```

#### **3. Verificar health check:**
```bash
curl http://localhost:3000/api/v1/health
```

### **📋 Verificações Adicionais:**

#### **1. Verificar se as tabelas foram criadas:**
```bash
./scripts/manage.sh exec-db psql -U postgres -d helpdanfe -c "\dt"
```

#### **2. Verificar estrutura das tabelas:**
```bash
./scripts/manage.sh exec-db psql -U postgres -d helpdanfe -c "\d nfe"
```

#### **3. Verificar logs do PostgreSQL:**
```bash
./scripts/manage.sh logs postgres
```

### **🐛 Outros Problemas Comuns:**

#### **Problema: Porta já em uso**
```bash
# Verificar portas em uso
netstat -tulpn | grep :3000
netstat -tulpn | grep :55432

# Parar processos usando as portas
sudo kill -9 <PID>
```

#### **Problema: Erro de permissão no script**
```bash
# Dar permissão de execução
chmod +x scripts/manage.sh
chmod +x scripts/test-db.sh
```

#### **Problema: Container não inicia**
```bash
# Ver logs detalhados
docker-compose logs -f

# Verificar recursos do sistema
docker system df
docker system prune -f
```

### **📞 Se o Problema Persistir:**

1. **Verificar logs completos:**
   ```bash
   ./scripts/manage.sh logs
   ```

2. **Verificar status dos serviços:**
   ```bash
   ./scripts/manage.sh status
   ```

3. **Verificar configuração do banco:**
   ```bash
   ./scripts/manage.sh test-db
   ```

4. **Criar issue no repositório** com:
   - Logs de erro completos
   - Versão do Docker e Docker Compose
   - Sistema operacional
   - Comandos executados

### **🔍 Comandos Úteis para Debug:**

```bash
# Ver logs em tempo real
./scripts/manage.sh logs -f

# Acessar shell do container da aplicação
./scripts/manage.sh exec-app sh

# Acessar shell do banco
./scripts/manage.sh exec-db psql -U postgres -d helpdanfe

# Ver estatísticas dos containers
docker stats

# Ver informações da rede Docker
docker network ls
docker network inspect helpdanfe-Go_helpdanfe-network
```

---

**💡 Dica:** Sempre use `./scripts/manage.sh` em vez de comandos Docker diretos para melhor controle e logs. 