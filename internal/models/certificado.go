package models

import (
	"time"
	"gorm.io/gorm"
)

type Certificado struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Nome	  string         `json:"nome"` gorm:`"not null"`
	Cnpj      string         `json:"cnpj"` gorm:`"not null;uniqueIndex"`
	arquivo   string         `json:"arquivo"` gorm:`"not null"`
	Senha     string         `json:"senha"` gorm:`"not null"`
	DataValidade time.Time      `json:"data_validade"` gorm:`"not null"`
	Ativo     bool           `json:"ativo"` gorm:`"not null;default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	func (Certificado) TableName() string {
		return "certificados"
	}
}