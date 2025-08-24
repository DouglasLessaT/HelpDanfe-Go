package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/config"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"

	"github.com/antchfx/xmlquery"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// NFEService representa o serviço de NFe
type NFEService struct {
	config *config.Config
	db     *gorm.DB
	logger *logrus.Logger
}

// NewNFEService cria uma nova instância do serviço de NFe
func NewNFEService(cfg *config.Config, db *gorm.DB, logger *logrus.Logger) *NFEService {
	return &NFEService{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

// ConsultarNFe consulta uma NFe na SEFAZ
func (s *NFEService) ConsultarNFe(chaveAcesso string) (*models.NFe, error) {
	s.logger.WithField("chave_acesso", chaveAcesso).Info("Consultando NFe")

	// Verifica se já existe no cache/banco
	var nfe models.NFe
	if err := s.db.Where("chave_acesso = ?", chaveAcesso).First(&nfe).Error; err == nil {
		s.logger.Info("NFe encontrada no cache")
		return &nfe, nil
	}

	// Consulta na SEFAZ
	xmlData, err := s.consultarSEFAZ(chaveAcesso)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar SEFAZ: %w", err)
	}

	// Parse do XML
	nfe, err = s.parseXMLNFe(xmlData)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do XML: %w", err)
	}

	// Salva no banco
	if err := s.db.Create(&nfe).Error; err != nil {
		s.logger.WithError(err).Error("Erro ao salvar NFe no banco")
	}

	return &nfe, nil
}

// BaixarXML baixa o XML de uma NFe
func (s *NFEService) BaixarXML(chaveAcesso string) (string, error) {
	s.logger.WithField("chave_acesso", chaveAcesso).Info("Baixando XML da NFe")

	// Verifica se existe no banco
	var nfe models.NFe
	if err := s.db.Where("chave_acesso = ?", chaveAcesso).First(&nfe).Error; err != nil {
		return "", fmt.Errorf("NFe não encontrada: %w", err)
	}

	if nfe.XML == "" {
		return "", fmt.Errorf("XML não disponível para esta NFe")
	}

	return nfe.XML, nil
}

// consultarSEFAZ consulta a SEFAZ para obter dados da NFe
func (s *NFEService) consultarSEFAZ(chaveAcesso string) (string, error) {
	s.logger.Info("Consultando SEFAZ")

	// Implementação da consulta SOAP para SEFAZ
	// Esta é uma implementação simplificada - na prática, você precisaria:
	// 1. Carregar o certificado digital
	// 2. Fazer requisição SOAP para o webservice da SEFAZ
	// 3. Tratar a resposta

	// Mock da resposta para demonstração
	xmlMock := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<nfeProc xmlns="http://www.portalfiscal.inf.br/nfe" versao="4.00">
  <NFe>
    <infNFe Id="NFe%s" versao="4.00">
      <ide>
        <cUF>35</cUF>
        <cNF>12345678</cNF>
        <natOp>Venda de mercadoria</natOp>
        <mod>55</mod>
        <serie>1</serie>
        <nNF>123456</nNF>
        <dhEmi>2024-01-01T10:00:00-03:00</dhEmi>
        <tpNF>1</tpNF>
        <idDest>1</idDest>
        <cMunFG>3550308</cMunFG>
        <tpImp>1</tpImp>
        <tpEmis>1</tpEmis>
        <cDV>1</cDV>
        <tpAmb>2</tpAmb>
        <finNFe>1</finNFe>
        <indFinal>1</indFinal>
        <indPres>1</indPres>
        <procEmi>0</procEmi>
        <verProc>1.0</verProc>
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
          <vBC>1000.00</vBC>
          <vICMS>180.00</vICMS>
          <vProd>1000.00</vProd>
          <vNF>1000.00</vNF>
        </ICMSTot>
      </total>
      <cobr>
        <dup>
          <nDup>001</nDup>
          <dVenc>2024-02-01</dVenc>
          <vDup>1000.00</vDup>
        </dup>
      </cobr>
    </infNFe>
  </NFe>
  <protNFe versao="4.00">
    <infProt>
      <tpAmb>2</tpAmb>
      <verAplic>1.0</verAplic>
      <chNFe>%s</chNFe>
      <dhRecbto>2024-01-01T10:05:00-03:00</dhRecbto>
      <nProt>123456789012345</nProt>
      <digVal>abc123</digVal>
      <cStat>100</cStat>
      <xMotivo>Autorizado o uso da NF-e</xMotivo>
    </infProt>
  </protNFe>
</nfeProc>`, chaveAcesso, chaveAcesso)

	return xmlMock, nil
}

// parseXMLNFe faz o parse do XML da NFe
func (s *NFEService) parseXMLNFe(xmlData string) (models.NFe, error) {
	var nfe models.NFe

	// Parse do XML usando xmlquery
	doc, err := xmlquery.Parse(strings.NewReader(xmlData))
	if err != nil {
		return nfe, fmt.Errorf("erro ao fazer parse do XML: %w", err)
	}

	// Extrai dados básicos
	if ide := xmlquery.FindOne(doc, "//ide"); ide != nil {
		if serie := xmlquery.FindOne(ide, "serie"); serie != nil {
			nfe.Serie = serie.InnerText()
		}
		if numero := xmlquery.FindOne(ide, "nNF"); numero != nil {
			nfe.Numero = numero.InnerText()
		}
		if dataEmissao := xmlquery.FindOne(ide, "dhEmi"); dataEmissao != nil {
			if t, err := time.Parse("2006-01-02T15:04:05-07:00", dataEmissao.InnerText()); err == nil {
				nfe.DataEmissao = t
			}
		}
	}

	// Extrai dados do emitente
	if emit := xmlquery.FindOne(doc, "//emit"); emit != nil {
		if cnpj := xmlquery.FindOne(emit, "CNPJ"); cnpj != nil {
			nfe.EmitenteCNPJ = cnpj.InnerText()
		}
		if nome := xmlquery.FindOne(emit, "xNome"); nome != nil {
			nfe.EmitenteNome = nome.InnerText()
		}
		if ie := xmlquery.FindOne(emit, "IE"); ie != nil {
			nfe.EmitenteIE = ie.InnerText()
		}
	}

	// Extrai dados do destinatário
	if dest := xmlquery.FindOne(doc, "//dest"); dest != nil {
		if cnpj := xmlquery.FindOne(dest, "CNPJ"); cnpj != nil {
			nfe.DestinatarioCNPJ = cnpj.InnerText()
		}
		if nome := xmlquery.FindOne(dest, "xNome"); nome != nil {
			nfe.DestinatarioNome = nome.InnerText()
		}
		if ie := xmlquery.FindOne(dest, "IE"); ie != nil {
			nfe.DestinatarioIE = ie.InnerText()
		}
	}

	// Extrai valores
	if total := xmlquery.FindOne(doc, "//total/ICMSTot"); total != nil {
		if valorTotal := xmlquery.FindOne(total, "vNF"); valorTotal != nil {
			if v, err := parseFloat(valorTotal.InnerText()); err == nil {
				nfe.ValorTotal = v
			}
		}
		if valorProdutos := xmlquery.FindOne(total, "vProd"); valorProdutos != nil {
			if v, err := parseFloat(valorProdutos.InnerText()); err == nil {
				nfe.ValorProdutos = v
			}
		}
		if valorImpostos := xmlquery.FindOne(total, "vICMS"); valorImpostos != nil {
			if v, err := parseFloat(valorImpostos.InnerText()); err == nil {
				nfe.ValorImpostos = v
			}
		}
	}

	// Extrai duplicatas
	if cobr := xmlquery.FindOne(doc, "//cobr"); cobr != nil {
		duplicatas := xmlquery.Find(cobr, "dup")
		for _, dup := range duplicatas {
			var duplicata models.Duplicata
			
			if numero := xmlquery.FindOne(dup, "nDup"); numero != nil {
				duplicata.Numero = numero.InnerText()
			}
			if vencimento := xmlquery.FindOne(dup, "dVenc"); vencimento != nil {
				if t, err := time.Parse("2006-01-02", vencimento.InnerText()); err == nil {
					duplicata.Vencimento = t
				}
			}
			if valor := xmlquery.FindOne(dup, "vDup"); valor != nil {
				if v, err := parseFloat(valor.InnerText()); err == nil {
					duplicata.Valor = v
				}
			}
			
			nfe.Duplicatas = append(nfe.Duplicatas, duplicata)
		}
	}

	// Define outros campos
	nfe.ChaveAcesso = extractChaveAcesso(xmlData)
	nfe.XML = xmlData
	nfe.Status = "AUTORIZADA"
	nfe.Ambiente = s.config.SEFAZ.Ambiente
	nfe.UF = s.config.SEFAZ.UF

	// Define data de autorização
	if prot := xmlquery.FindOne(doc, "//protNFe/infProt"); prot != nil {
		if dataRecbto := xmlquery.FindOne(prot, "dhRecbto"); dataRecbto != nil {
			if t, err := time.Parse("2006-01-02T15:04:05-07:00", dataRecbto.InnerText()); err == nil {
				nfe.DataAutorizacao = &t
			}
		}
	}

	return nfe, nil
}

// extractChaveAcesso extrai a chave de acesso do XML
func extractChaveAcesso(xmlData string) string {
	doc, err := xmlquery.Parse(strings.NewReader(xmlData))
	if err != nil {
		return ""
	}

	if chNFe := xmlquery.FindOne(doc, "//chNFe"); chNFe != nil {
		return chNFe.InnerText()
	}

	return ""
}

// parseFloat converte string para float64
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
