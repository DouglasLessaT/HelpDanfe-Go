package dto

import "time"

type BoletoDTO struct {
	ID             uint       `json:"id"`
	DuplicataID    *uint      `json:"duplicata_id,omitempty"`
	Banco          string     `json:"banco"`
	Numero         string     `json:"numero"`
	CodigoBarras   string     `json:"codigo_barras"`
	LinhaDigitavel string     `json:"linha_digitavel"`
	Valor          float64    `json:"valor"`
	Vencimento     time.Time  `json:"vencimento"`
	Status         string     `json:"status"`
	DataPagamento  *time.Time `json:"data_pagamento,omitempty"`
	ValorPago      *float64   `json:"valor_pago,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
