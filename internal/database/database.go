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
<<<<<<< HEAD
	log.Printf("Conectando ao banco de dados: %s:%s/%s", cfg.Host, cfg.Port, cfg.Name)
=======
	log.Printf("Conectando ao banco de dados com DSN: %s:%s/%s", cfg.Host, cfg.Port, cfg.DBName)
>>>>>>> 324acbc84a8f9e1d24b40ca72dc1246f9891ceb3
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
<<<<<<< HEAD
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
=======
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err "linha 25")
>>>>>>> 324acbc84a8f9e1d24b40ca72dc1246f9891ceb3
	}

	// Configura pool de conexões
	sqlDB, err := db.DB()
	if err != nil {
<<<<<<< HEAD
		return nil, fmt.Errorf("erro ao obter instância SQL: %w", err)
=======
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err "linha 31")
>>>>>>> 324acbc84a8f9e1d24b40ca72dc1246f9891ceb3
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	log.Println("Conexão com o banco de dados estabelecida com sucesso")

<<<<<<< HEAD
	log.Printf("Conexão com banco estabelecida com sucesso")
	return db, nil
=======
	return db, fmt.Errorf("erro ao conectar ao banco de dados: %w", err,"linha 37")
>>>>>>> 324acbc84a8f9e1d24b40ca72dc1246f9891ceb3
}

// Migrate executa as migrações do banco de dados
func Migrate(db *gorm.DB) error {
	log.Println("Iniciando migrações do banco de dados...")
	
	// Lista de modelos para migração
	modelsToMigrate := []interface{}{
		&models.NFe{},
		&models.Duplicata{},
		&models.Boleto{},
<<<<<<< HEAD
		&models.Certificado{},
		&models.LogConsulta{},
	}

	// Executa migração para cada modelo
	for _, model := range modelsToMigrate {
		log.Printf("Migrando modelo: %T", model)
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("erro ao migrar modelo %T: %w", model, err)
		}
		log.Printf("Modelo %T migrado com sucesso", model)
	}

	log.Println("Todas as migrações foram executadas com sucesso")
	return nil
=======
		&models.LogConsulta{},
		&models.Certificado{}
	)
>>>>>>> 324acbc84a8f9e1d24b40ca72dc1246f9891ceb3
}

// GetDB retorna a instância do banco de dados
func GetDB() *gorm.DB {
	// Esta função seria usada em outros pacotes para obter a instância do DB
	// Por enquanto, retorna nil - será implementada com injeção de dependência
	return nil
}
