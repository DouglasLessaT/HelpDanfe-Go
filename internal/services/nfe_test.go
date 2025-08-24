package services

import (
	"testing"

	"helpdanfe-go/internal/config"
	"helpdanfe-go/internal/models"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto migrate
	db.AutoMigrate(&models.NFe{}, &models.Duplicata{}, &models.Boleto{})

	return db
}

func setupTestConfig() *config.Config {
	return &config.Config{
		SEFAZ: config.SEFAZConfig{
			Ambiente: "homologacao",
			UF:       "SP",
		},
	}
}

func TestConsultarNFe(t *testing.T) {
	db := setupTestDB()
	cfg := setupTestConfig()
	logger := logrus.New()

	service := NewNFEService(cfg, db, logger)

	// Teste com chave válida
	chave := "12345678901234567890123456789012345678901234"
	nfe, err := service.ConsultarNFe(chave)

	assert.NoError(t, err)
	assert.NotNil(t, nfe)
	assert.Equal(t, chave, nfe.ChaveAcesso)
	assert.Equal(t, "AUTORIZADA", nfe.Status)
}

func TestBaixarXML(t *testing.T) {
	db := setupTestDB()
	cfg := setupTestConfig()
	logger := logrus.New()

	service := NewNFEService(cfg, db, logger)

	// Primeiro consulta a NFe
	chave := "12345678901234567890123456789012345678901234"
	_, err := service.ConsultarNFe(chave)
	assert.NoError(t, err)

	// Depois baixa o XML
	xml, err := service.BaixarXML(chave)

	assert.NoError(t, err)
	assert.NotEmpty(t, xml)
	assert.Contains(t, xml, "nfeProc")
}

func TestParseXMLNFe(t *testing.T) {
	db := setupTestDB()
	cfg := setupTestConfig()
	logger := logrus.New()

	service := NewNFEService(cfg, db, logger)

	// XML de teste
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<nfeProc xmlns="http://www.portalfiscal.inf.br/nfe" versao="4.00">
  <NFe>
    <infNFe Id="NFe12345678901234567890123456789012345678901234" versao="4.00">
      <ide>
        <serie>1</serie>
        <nNF>123456</nNF>
        <dhEmi>2024-01-01T10:00:00-03:00</dhEmi>
      </ide>
      <emit>
        <CNPJ>12345678000123</CNPJ>
        <xNome>EMPRESA EXEMPLO LTDA</xNome>
        <IE>123456789</IE>
      </emit>
      <dest>
        <CNPJ>98765432000198</CNPJ>
        <xNome>CLIENTE EXEMPLO LTDA</xNome>
        <IE>987654321</IE>
      </dest>
      <total>
        <ICMSTot>
          <vNF>1000.00</vNF>
          <vProd>1000.00</vProd>
          <vICMS>180.00</vICMS>
        </ICMSTot>
      </total>
    </infNFe>
  </NFe>
  <protNFe versao="4.00">
    <infProt>
      <chNFe>12345678901234567890123456789012345678901234</chNFe>
      <dhRecbto>2024-01-01T10:05:00-03:00</dhRecbto>
    </infProt>
  </protNFe>
</nfeProc>`

	nfe, err := service.parseXMLNFe(xmlData)

	assert.NoError(t, err)
	assert.Equal(t, "12345678901234567890123456789012345678901234", nfe.ChaveAcesso)
	assert.Equal(t, "123456", nfe.Numero)
	assert.Equal(t, "1", nfe.Serie)
	assert.Equal(t, "12345678000123", nfe.EmitenteCNPJ)
	assert.Equal(t, "EMPRESA EXEMPLO LTDA", nfe.EmitenteNome)
	assert.Equal(t, "98765432000198", nfe.DestinatarioCNPJ)
	assert.Equal(t, "CLIENTE EXEMPLO LTDA", nfe.DestinatarioNome)
	assert.Equal(t, "AUTORIZADA", nfe.Status)
}
