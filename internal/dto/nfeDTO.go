package dto

import "time"

type NFeDTO struct {
	ID                uint           `json:"id"`
	ChaveAcesso       string         `json:"chave_acesso"`
	Numero            string         `json:"numero"`
	Serie             string         `json:"serie"`
	DataEmissao       time.Time      `json:"data_emissao"`
	DataAutorizacao   *time.Time     `json:"data_autorizacao,omitempty"`
	Status            string         `json:"status"`
	Ambiente          string         `json:"ambiente"`
	UF                string         `json:"uf"`
	EmitenteCNPJ      string         `json:"emitente_cnpj"`
	EmitenteNome      string         `json:"emitente_nome"`
	EmitenteIE        string         `json:"emitente_ie"`
	DestinatarioCNPJ  string         `json:"destinatario_cnpj"`
	DestinatarioNome  string         `json:"destinatario_nome"`
	DestinatarioIE    string         `json:"destinatario_ie"`
	ValorTotal        float64        `json:"valor_total"`
	ValorProdutos     float64        `json:"valor_produtos"`
	ValorImpostos     float64        `json:"valor_impostos"`
	Duplicatas        []DuplicataDTO `json:"duplicatas,omitempty"`
	Boletos           []BoletoDTO    `json:"boletos,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}
