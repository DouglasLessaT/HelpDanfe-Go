package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// JSONB é um tipo personalizado para campos JSONB do PostgreSQL
type JSONB map[string]interface{}

// Value implementa driver.Valuer para JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implementa sql.Scanner para JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}
	
	return json.Unmarshal(bytes, j)
}

// LogConsulta representa um log de consulta realizada
type LogConsulta struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TipoConsulta   string         `json:"tipo_consulta" gorm:"not null"`
	Parametros     JSONB          `json:"parametros"`
	Resultado      JSONB          `json:"resultado"`
	TempoExecucao  int            `json:"tempo_execucao"` // em milissegundos
	Status         string         `json:"status" gorm:"not null"`
	Erro           string         `json:"erro"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName garante o nome correto da tabela
func (LogConsulta) TableName() string {
	return "logs_consulta"
} 