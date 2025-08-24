-- Script de inicialização do banco de dados HelpDanfe
-- Este script é executado automaticamente quando o container PostgreSQL é criado

-- Criação das tabelas principais

-- Tabela de NFe
CREATE TABLE IF NOT EXISTS nfe (
    id SERIAL PRIMARY KEY,
    chave VARCHAR(44) UNIQUE NOT NULL,
    numero VARCHAR(9) NOT NULL,
    serie VARCHAR(3) NOT NULL,
    data_emissao TIMESTAMP NOT NULL,
    valor_total DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pendente',
    xml_content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Boletos
CREATE TABLE IF NOT EXISTS boletos (
    id SERIAL PRIMARY KEY,
    codigo_barras VARCHAR(44) UNIQUE NOT NULL,
    linha_digitavel VARCHAR(47) NOT NULL,
    valor DECIMAL(15,2) NOT NULL,
    data_vencimento DATE NOT NULL,
    banco VARCHAR(10) NOT NULL,
    agencia VARCHAR(10),
    conta VARCHAR(20),
    nosso_numero VARCHAR(20),
    status VARCHAR(20) DEFAULT 'pendente',
    nfe_id INTEGER REFERENCES nfe(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Duplicatas
CREATE TABLE IF NOT EXISTS duplicatas (
    id SERIAL PRIMARY KEY,
    numero VARCHAR(20) NOT NULL,
    data_vencimento DATE NOT NULL,
    valor DECIMAL(15,2) NOT NULL,
    nfe_id INTEGER REFERENCES nfe(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Certificados
CREATE TABLE IF NOT EXISTS certificados (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) NOT NULL,
    arquivo_path VARCHAR(500) NOT NULL,
    senha VARCHAR(255),
    data_validade DATE NOT NULL,
    ativo BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Logs de Consulta
CREATE TABLE IF NOT EXISTS logs_consulta (
    id SERIAL PRIMARY KEY,
    tipo_consulta VARCHAR(50) NOT NULL,
    parametros JSONB,
    resultado JSONB,
    tempo_execucao INTEGER, -- em milissegundos
    status VARCHAR(20) NOT NULL,
    erro TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para melhorar performance
CREATE INDEX IF NOT EXISTS idx_nfe_chave ON nfe(chave);
CREATE INDEX IF NOT EXISTS idx_nfe_numero_serie ON nfe(numero, serie);
CREATE INDEX IF NOT EXISTS idx_boletos_codigo ON boletos(codigo_barras);
CREATE INDEX IF NOT EXISTS idx_boletos_nfe ON boletos(nfe_id);
CREATE INDEX IF NOT EXISTS idx_duplicatas_nfe ON duplicatas(nfe_id);
CREATE INDEX IF NOT EXISTS idx_certificados_cnpj ON certificados(cnpj);
CREATE INDEX IF NOT EXISTS idx_logs_tipo ON logs_consulta(tipo_consulta);
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON logs_consulta(created_at);

-- Função para atualizar o timestamp de updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers para atualizar updated_at automaticamente
CREATE TRIGGER update_nfe_updated_at BEFORE UPDATE ON nfe
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_boletos_updated_at BEFORE UPDATE ON boletos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certificados_updated_at BEFORE UPDATE ON certificados
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Inserir dados de exemplo (opcional para desenvolvimento)
INSERT INTO certificados (nome, cnpj, arquivo_path, data_validade) 
VALUES ('Certificado Exemplo', '12345678000199', './certs/certificado.p12', '2025-12-31')
ON CONFLICT DO NOTHING;

-- Comentários sobre as tabelas
COMMENT ON TABLE nfe IS 'Tabela para armazenar informações das Notas Fiscais Eletrônicas';
COMMENT ON TABLE boletos IS 'Tabela para armazenar informações dos boletos bancários';
COMMENT ON TABLE duplicatas IS 'Tabela para armazenar informações das duplicatas';
COMMENT ON TABLE certificados IS 'Tabela para armazenar informações dos certificados digitais';
COMMENT ON TABLE logs_consulta IS 'Tabela para armazenar logs das consultas realizadas'; 