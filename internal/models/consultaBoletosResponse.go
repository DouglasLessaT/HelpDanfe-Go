package models

type ConsultaBoletosResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []Boleto `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
}