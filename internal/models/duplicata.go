package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/dto"
)

type Duplicata struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	NFeID       uint           `json:"nfe_id"`
	Numero      string         `json:"numero"`
	Vencimento  time.Time      `json:"vencimento"`
	Valor       float64        `json:"valor"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (d *Duplicata) ToDTO() dto.DuplicataDTO {
	return dto.DuplicataDTO{
		ID:         d.ID,
		Numero:     d.Numero,
		Vencimento: d.Vencimento,
		Valor:      d.Valor,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}
}