package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/dto"
)

type Boleto struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	NFeID           uint           `json:"nfe_id"`
	DuplicataID     *uint          `json:"duplicata_id"`
	Banco           string         `json:"banco"`
	Numero          string         `json:"numero"`
	CodigoBarras    string         `json:"codigo_barras"`
	LinhaDigitavel  string         `json:"linha_digitavel"`
	Valor           float64        `json:"valor"`
	Vencimento      time.Time      `json:"vencimento"`
	Status          string         `json:"status"`
	DataPagamento   *time.Time     `json:"data_pagamento"`
	ValorPago       *float64       `json:"valor_pago"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (b *Boleto) ToDTO() dto.BoletoDTO {
	return dto.BoletoDTO{
		ID:             b.ID,
		DuplicataID:    b.DuplicataID,
		Banco:          b.Banco,
		Numero:         b.Numero,
		CodigoBarras:   b.CodigoBarras,
		LinhaDigitavel: b.LinhaDigitavel,
		Valor:          b.Valor,
		Vencimento:     b.Vencimento,
		Status:         b.Status,
		DataPagamento:  b.DataPagamento,
		ValorPago:      b.ValorPago,
		CreatedAt:      b.CreatedAt,
		UpdatedAt:      b.UpdatedAt,
	}
}
// TableName garante o nome correto da tabela para Boleto
func (Boleto) TableName() string {
	return "boletos"
}
