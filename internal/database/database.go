package database

import (
	"fmt"
	"log"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/config"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect estabelece conexão com o banco de dados
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	log.Printf("Conectando ao banco de dados com DSN: %s:%s/%s", cfg.Host, cfg.Port, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err "linha 25")
	}

	// Configura pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err "linha 31")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	log.Println("Conexão com o banco de dados estabelecida com sucesso")

	return db, fmt.Errorf("erro ao conectar ao banco de dados: %w", err,"linha 37")
}

// Migrate executa as migrações do banco de dados
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.NFe{},
		&models.Duplicata{},
		&models.Boleto{},
		&models.LogConsulta{},
		&models.Certificado{}
	)
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	// Esta função seria usada em outros pacotes para obter a instância do DB
	// Por enquanto, retorna nil - será implementada com injeção de dependência
	return nil
}
