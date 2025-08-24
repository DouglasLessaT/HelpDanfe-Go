package models

type ConsultaNFeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *NFe   `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ConsultaNFeRequest struct {
	ChaveAcesso string `json:"chave_acesso" binding:"required,len=44"`
	Certificado  string `json:"certificado,omitempty"`
	Senha        string `json:"senha,omitempty"`
}