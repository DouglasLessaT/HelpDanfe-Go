package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/config"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/database"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/handlers"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/middleware"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	// Configura logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	// Carrega configurações
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Conecta ao banco de dados
	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Executa migrações
	if err := database.Migrate(db); err != nil {
		logger.Fatalf("Erro ao executar migrações: %v", err)
	}

	// Inicializa serviços
	nfeService := services.NewNFEService(cfg, db, logger)
	bankService := services.NewBankService(cfg, db, logger)
	pdfService := services.NewPDFService(logger)

	// Configura router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())

	// Configura rotas da API
	api := router.Group("/api/v1")
	{
		// Rotas de NFe
		nfeGroup := api.Group("/nfe")
		{
			nfeGroup.POST("/consultar", handlers.ConsultarNFe(nfeService))
			nfeGroup.GET("/:chave/xml", handlers.BaixarXMLNFe(nfeService))
			nfeGroup.GET("/:chave/pdf", handlers.GerarPDFNFe(nfeService, pdfService))
			nfeGroup.GET("/:chave/boletos", handlers.ConsultarBoletosNFe(nfeService, bankService))
		}

		// Rotas de Boletos
		boletosGroup := api.Group("/boletos")
		{
			boletosGroup.GET("/:codigo", handlers.ConsultarBoleto(bankService))
			boletosGroup.POST("/consultar", handlers.ConsultarMultiplosBoletos(bankService))
		}

		// Rotas de Certificados
		certificadosGroup := api.Group("/certificados")
		{
			certificadosGroup.GET("/verificar", handlers.VerificarCertificados())
			certificadosGroup.POST("/selecionar", handlers.SelecionarCertificado())
		}

		// Rota de health check
		api.GET("/health", handlers.HealthCheck())
	}

	// Configura rotas para arquivos estáticos (frontend)
	router.Static("/css", "./web/css")
	router.Static("/js", "./web/js")
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/index.html", "./web/index.html")

	// Configura servidor HTTP
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Canal para receber sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Inicia servidor em goroutine
	go func() {
		logger.Infof("Servidor iniciado na porta %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguarda sinal de interrupção
	<-quit
	logger.Info("Desligando servidor...")

	// Contexto com timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown graceful
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Erro durante shutdown: %v", err)
	}

	logger.Info("Servidor desligado com sucesso")
}
