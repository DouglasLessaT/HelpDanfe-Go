package models

type ConsultaBoletosRequest struct {
	Codigos []string `json:"codigos" binding:"required"`
}