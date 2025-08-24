package models

import (
	"time"

	"gorm.io/gorm"
)

// Certificado representa um certificado digital
type Certificado struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Nome        string         `json:"nome" gorm:"not null"`
	CNPJ        string         `json:"cnpj" gorm:"not null;uniqueIndex"`
	ArquivoPath string         `json:"arquivo_path" gorm:"not null"`
	Senha       string         `json:"senha"`
	DataValidade time.Time     `json:"data_validade" gorm:"not null"`
	Ativo       bool           `json:"ativo" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName garante o nome correto da tabela
func (Certificado) TableName() string {
	return "certificados"
} 