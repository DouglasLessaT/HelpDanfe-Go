package services

import (
	"fmt"
	"time"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/config"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// BankService representa o serviço bancário
type BankService struct {
	config *config.Config
	db     *gorm.DB
	logger *logrus.Logger
}

// NewBankService cria uma nova instância do serviço bancário
func NewBankService(cfg *config.Config, db *gorm.DB, logger *logrus.Logger) *BankService {
	return &BankService{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

// ConsultarBoleto consulta um boleto específico
func (s *BankService) ConsultarBoleto(codigo string) (*models.Boleto, error) {
	s.logger.WithField("codigo", codigo).Info("Consultando boleto")

	// Verifica se existe no banco
	var boleto models.Boleto
	if err := s.db.Where("numero = ? OR codigo_barras = ?", codigo, codigo).First(&boleto).Error; err == nil {
		return &boleto, nil
	}

	// Consulta nas APIs bancárias
	boleto, err := s.consultarAPIsBancarias(codigo)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar APIs bancárias: %w", err)
	}

	// Salva no banco
	if err := s.db.Create(&boleto).Error; err != nil {
		s.logger.WithError(err).Error("Erro ao salvar boleto no banco")
	}

	return &boleto, nil
}

// ConsultarMultiplosBoletos consulta múltiplos boletos
func (s *BankService) ConsultarMultiplosBoletos(codigos []string) ([]models.Boleto, error) {
	s.logger.WithField("codigos", codigos).Info("Consultando múltiplos boletos")

	var boletos []models.Boleto

	for _, codigo := range codigos {
		boleto, err := s.ConsultarBoleto(codigo)
		if err != nil {
			s.logger.WithError(err).WithField("codigo", codigo).Error("Erro ao consultar boleto")
			continue
		}
		boletos = append(boletos, *boleto)
	}

	return boletos, nil
}

// ConsultarBoletosPorNFe consulta boletos vinculados a uma NFe
func (s *BankService) ConsultarBoletosPorNFe(nfeID uint) ([]models.Boleto, error) {
	s.logger.WithField("nfe_id", nfeID).Info("Consultando boletos por NFe")

	var boletos []models.Boleto
	if err := s.db.Where("nfe_id = ?", nfeID).Find(&boletos).Error; err != nil {
		return nil, fmt.Errorf("erro ao consultar boletos no banco: %w", err)
	}

	// Se não encontrou boletos, tenta localizar pelas duplicatas
	if len(boletos) == 0 {
		var localizados []models.Boleto
		localizados, localizarErr := s.localizarBoletosPorDuplicatas(nfeID)
		if localizarErr != nil {
			s.logger.WithError(localizarErr).Error("Erro ao localizar boletos por duplicatas")
		} else {
			boletos = localizados
		}
	}

	return boletos, nil
}

// ConsultarBoletosPorDuplicata consulta boletos por duplicata
func (s *BankService) ConsultarBoletosPorDuplicata(duplicataID uint) ([]models.Boleto, error) {
	s.logger.WithField("duplicata_id", duplicataID).Info("Consultando boletos por duplicata")

	var boletos []models.Boleto
	if err := s.db.Where("duplicata_id = ?", duplicataID).Find(&boletos).Error; err != nil {
		return nil, fmt.Errorf("erro ao consultar boletos no banco: %w", err)
	}

	return boletos, nil
}

// GerarCodigoBarras gera código de barras para um boleto
func (s *BankService) GerarCodigoBarras(codigo string) (string, error) {
	s.logger.WithField("codigo", codigo).Info("Gerando código de barras")

	// Consulta o boleto
	boleto, err := s.ConsultarBoleto(codigo)
	if err != nil {
		return "", fmt.Errorf("erro ao consultar boleto: %w", err)
	}

	if boleto.CodigoBarras != "" {
		return boleto.CodigoBarras, nil
	}

	// Gera código de barras baseado no padrão FEBRABAN
	codigoBarras, err := s.gerarCodigoBarrasFEBRABAN(boleto)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar código de barras: %w", err)
	}

	// Atualiza o boleto no banco
	boleto.CodigoBarras = codigoBarras
	if err := s.db.Save(boleto).Error; err != nil {
		s.logger.WithError(err).Error("Erro ao atualizar código de barras no banco")
	}

	return codigoBarras, nil
}

// consultarAPIsBancarias consulta as APIs bancárias para localizar boletos
func (s *BankService) consultarAPIsBancarias(codigo string) (models.Boleto, error) {
	s.logger.Info("Consultando APIs bancárias")

	// Tenta consultar em diferentes bancos
	boletos := []models.Boleto{}

	// Consulta Itaú
	if s.config.Bank.Itau.URL != "" {
		if boleto, err := s.consultarItau(codigo); err == nil {
			boletos = append(boletos, boleto)
		}
	}

	// Consulta Bradesco
	if s.config.Bank.Bradesco.URL != "" {
		if boleto, err := s.consultarBradesco(codigo); err == nil {
			boletos = append(boletos, boleto)
		}
	}

	// Consulta Open Banking
	if s.config.Bank.OpenBanking.URL != "" {
		if boleto, err := s.consultarOpenBanking(codigo); err == nil {
			boletos = append(boletos, boleto)
		}
	}

	// Retorna o primeiro boleto encontrado
	if len(boletos) > 0 {
		return boletos[0], nil
	}

	// Mock de boleto para demonstração
	return models.Boleto{
		Numero:         codigo,
		Banco:          "001",
		CodigoBarras:   "00193373700000001000500940144816060680935031",
		LinhaDigitavel: "00190.00009 04441.601448 60606.809350 3 37370000000100",
		Valor:          1000.00,
		Vencimento:     time.Now().AddDate(0, 1, 0),
		Status:         "ABERTO",
	}, nil
}

// consultarItau consulta a API do Itaú
func (s *BankService) consultarItau(codigo string) (models.Boleto, error) {
	s.logger.Info("Consultando API do Itaú")

	// Implementação da consulta à API do Itaú
	// Esta é uma implementação simplificada

	return models.Boleto{
		Numero:         codigo,
		Banco:          "341",
		CodigoBarras:   "34191790010104351004791020150008787870026300",
		LinhaDigitavel: "34191.79001 01043.510047 91020.150008 8 787870026300",
		Valor:          1000.00,
		Vencimento:     time.Now().AddDate(0, 1, 0),
		Status:         "ABERTO",
	}, nil
}

// consultarBradesco consulta a API do Bradesco
func (s *BankService) consultarBradesco(codigo string) (models.Boleto, error) {
	s.logger.Info("Consultando API do Bradesco")

	// Implementação da consulta à API do Bradesco
	// Esta é uma implementação simplificada

	return models.Boleto{
		Numero:         codigo,
		Banco:          "237",
		CodigoBarras:   "23793381286000782713695000063305975820000126000",
		LinhaDigitavel: "23793.38128 60007.827136 95000.063305 9 75820000126000",
		Valor:          1000.00,
		Vencimento:     time.Now().AddDate(0, 1, 0),
		Status:         "ABERTO",
	}, nil
}

// consultarOpenBanking consulta a API do Open Banking
func (s *BankService) consultarOpenBanking(codigo string) (models.Boleto, error) {
	s.logger.Info("Consultando API do Open Banking")

	// Implementação da consulta à API do Open Banking
	// Esta é uma implementação simplificada

	return models.Boleto{
		Numero:         codigo,
		Banco:          "001",
		CodigoBarras:   "00193373700000001000500940144816060680935031",
		LinhaDigitavel: "00190.00009 04441.601448 60606.809350 3 37370000000100",
		Valor:          1000.00,
		Vencimento:     time.Now().AddDate(0, 1, 0),
		Status:         "ABERTO",
	}, nil
}

// localizarBoletosPorDuplicatas localiza boletos baseado nas duplicatas da NFe
func (s *BankService) localizarBoletosPorDuplicatas(nfeID uint) ([]models.Boleto, error) {
	s.logger.WithField("nfe_id", nfeID).Info("Localizando boletos por duplicatas")

	var duplicatas []models.Duplicata
	if err := s.db.Where("nfe_id = ?", nfeID).Find(&duplicatas).Error; err != nil {
		return nil, fmt.Errorf("erro ao consultar duplicatas: %w", err)
	}

	var boletos []models.Boleto
	for _, duplicata := range duplicatas {
		// Tenta localizar boletos baseado no número da duplicata
		codigo := fmt.Sprintf("DUP%s", duplicata.Numero)
		if boleto, err := s.ConsultarBoleto(codigo); err == nil {
			boleto.DuplicataID = &duplicata.ID
			boletos = append(boletos, *boleto)
		}
	}

	return boletos, nil
}

// gerarCodigoBarrasFEBRABAN gera código de barras no padrão FEBRABAN
func (s *BankService) gerarCodigoBarrasFEBRABAN(boleto *models.Boleto) (string, error) {
	s.logger.Info("Gerando código de barras FEBRABAN")

	// Implementação simplificada do padrão FEBRABAN
	// Na prática, você precisaria implementar o algoritmo completo

	// Exemplo de código de barras do Banco do Brasil
	return "00193373700000001000500940144816060680935031", nil
}
