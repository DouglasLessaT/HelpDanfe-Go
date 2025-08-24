package database

import (
	"github.com/Douglaslessat/HelpDanfe-Go/internal/config"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect estabelece conexão com o banco de dados
func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDatabaseDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		return nil, err
	}

	// Configura pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// Migrate executa as migrações do banco de dados
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.NFe{},
		&models.Duplicata{},
		&models.Boleto{},
	)
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	// Esta função seria usada em outros pacotes para obter a instância do DB
	// Por enquanto, retorna nil - será implementada com injeção de dependência
	return nil
}
