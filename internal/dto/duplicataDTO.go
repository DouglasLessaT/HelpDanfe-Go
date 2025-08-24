package dto

import "time"

type DuplicataDTO struct {
	ID         uint      `json:"id"`
	Numero     string    `json:"numero"`
	Vencimento time.Time `json:"vencimento"`
	Valor      float64   `json:"valor"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
