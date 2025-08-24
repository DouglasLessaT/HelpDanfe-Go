package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/dto"
)

// NFe representa uma Nota Fiscal Eletrônica
type NFe struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	ChaveAcesso     string         `json:"chave_acesso" gorm:"uniqueIndex;not null"`
	Numero          string         `json:"numero"`
	Serie           string         `json:"serie"`
	DataEmissao     time.Time      `json:"data_emissao"`
	DataAutorizacao *time.Time     `json:"data_autorizacao"`
	Status          string         `json:"status"`
	Ambiente        string         `json:"ambiente"`
	UF              string         `json:"uf"`
	XML             string         `json:"xml" gorm:"type:text"`
	PDF             []byte         `json:"pdf" gorm:"type:bytea"`

	// Dados do emitente
	EmitenteCNPJ    string `json:"emitente_cnpj"`
	EmitenteNome    string `json:"emitente_nome"`
	EmitenteIE      string `json:"emitente_ie"`

	// Dados do destinatário
	DestinatarioCNPJ string `json:"destinatario_cnpj"`
	DestinatarioNome string `json:"destinatario_nome"`
	DestinatarioIE   string `json:"destinatario_ie"`

	// Valores
	ValorTotal      float64 `json:"valor_total"`
	ValorProdutos   float64 `json:"valor_produtos"`
	ValorImpostos   float64 `json:"valor_impostos"`

	// Relacionamentos
	Duplicatas      []Duplicata `json:"duplicatas" gorm:"foreignKey:NFeID"`
	Boletos         []Boleto    `json:"boletos" gorm:"foreignKey:NFeID"`

	// Metadados
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (n *NFe) ToDTO() dto.NFeDTO {
	duplicatas := make([]dto.DuplicataDTO, len(n.Duplicatas))
	for i, d := range n.Duplicatas {
		duplicatas[i] = d.ToDTO()
	}
	boletos := make([]dto.BoletoDTO, len(n.Boletos))
	for i, b := range n.Boletos {
		boletos[i] = b.ToDTO()
	}
	return dto.NFeDTO{
		ID:               n.ID,
		ChaveAcesso:      n.ChaveAcesso,
		Numero:           n.Numero,
		Serie:            n.Serie,
		DataEmissao:      n.DataEmissao,
		DataAutorizacao:  n.DataAutorizacao,
		Status:           n.Status,
		Ambiente:         n.Ambiente,
		UF:               n.UF,
		EmitenteCNPJ:     n.EmitenteCNPJ,
		EmitenteNome:     n.EmitenteNome,
		EmitenteIE:       n.EmitenteIE,
		DestinatarioCNPJ: n.DestinatarioCNPJ,
		DestinatarioNome: n.DestinatarioNome,
		DestinatarioIE:   n.DestinatarioIE,
		ValorTotal:       n.ValorTotal,
		ValorProdutos:    n.ValorProdutos,
		ValorImpostos:    n.ValorImpostos,
		Duplicatas:       duplicatas,
		Boletos:          boletos,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}
}
