package services

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"
	"github.com/jung-kurt/gofpdf"
	"github.com/sirupsen/logrus"
)

// PDFService representa o serviço de geração de PDFs
type PDFService struct {
	logger *logrus.Logger
}

// NewPDFService cria uma nova instância do serviço de PDF
func NewPDFService(logger *logrus.Logger) *PDFService {
	return &PDFService{
		logger: logger,
	}
}

// GerarDANFE gera o DANFE (Documento Auxiliar da Nota Fiscal Eletrônica) em PDF
func (s *PDFService) GerarDANFE(nfe *models.NFe) ([]byte, error) {
	s.logger.WithField("chave_acesso", nfe.ChaveAcesso).Info("Gerando DANFE")

	// Cria novo documento PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Configura fonte
	pdf.SetFont("Arial", "B", 12)

	// Cabeçalho
	s.adicionarCabecalho(pdf, nfe)

	// Dados do emitente
	s.adicionarEmitente(pdf, nfe)

	// Dados do destinatário
	s.adicionarDestinatario(pdf, nfe)

	// Dados da NFe
	s.adicionarDadosNFe(pdf, nfe)

	// Valores
	s.adicionarValores(pdf, nfe)

	// Duplicatas
	s.adicionarDuplicatas(pdf, nfe)

	// Código de barras
	s.adicionarCodigoBarras(pdf, nfe)

	// Gera o PDF
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// adicionarCabecalho adiciona o cabeçalho do DANFE
func (s *PDFService) adicionarCabecalho(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "DANFE - Documento Auxiliar da Nota Fiscal Eletrônica")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Chave de Acesso: "+nfe.ChaveAcesso)
	pdf.Ln(10)
}

// adicionarEmitente adiciona os dados do emitente
func (s *PDFService) adicionarEmitente(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "EMITENTE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "CNPJ: "+nfe.EmitenteCNPJ)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Nome: "+nfe.EmitenteNome)
	pdf.Ln(6)
	pdf.Cell(0, 6, "IE: "+nfe.EmitenteIE)
	pdf.Ln(10)
}

// adicionarDestinatario adiciona os dados do destinatário
func (s *PDFService) adicionarDestinatario(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "DESTINATÁRIO")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "CNPJ: "+nfe.DestinatarioCNPJ)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Nome: "+nfe.DestinatarioNome)
	pdf.Ln(6)
	pdf.Cell(0, 6, "IE: "+nfe.DestinatarioIE)
	pdf.Ln(10)
}

// adicionarDadosNFe adiciona os dados da NFe
func (s *PDFService) adicionarDadosNFe(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "DADOS DA NFE")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "Número: "+nfe.Numero)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Série: "+nfe.Serie)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Data de Emissão: "+nfe.DataEmissao.Format("02/01/2006 15:04:05"))
	pdf.Ln(6)
	if nfe.DataAutorizacao != nil {
		pdf.Cell(0, 6, "Data de Autorização: "+nfe.DataAutorizacao.Format("02/01/2006 15:04:05"))
		pdf.Ln(6)
	}
	pdf.Cell(0, 6, "Status: "+nfe.Status)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Ambiente: "+nfe.Ambiente)
	pdf.Ln(6)
	pdf.Cell(0, 6, "UF: "+nfe.UF)
	pdf.Ln(10)
}

// adicionarValores adiciona os valores da NFe
func (s *PDFService) adicionarValores(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "VALORES")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, "Valor Total: R$ "+formatarValor(nfe.ValorTotal))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Valor dos Produtos: R$ "+formatarValor(nfe.ValorProdutos))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Valor dos Impostos: R$ "+formatarValor(nfe.ValorImpostos))
	pdf.Ln(10)
}

// adicionarDuplicatas adiciona as duplicatas da NFe
func (s *PDFService) adicionarDuplicatas(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	if len(nfe.Duplicatas) == 0 {
		return
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "DUPLICATAS")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	for _, dup := range nfe.Duplicatas {
		pdf.Cell(0, 6, fmt.Sprintf("Número: %s, Vencimento: %s, Valor: R$ %s",
			dup.Numero,
			dup.Vencimento.Format("02/01/2006"),
			formatarValor(dup.Valor)))
		pdf.Ln(6)
	}
	pdf.Ln(10)
}

// adicionarCodigoBarras adiciona o código de barras da NFe
func (s *PDFService) adicionarCodigoBarras(pdf *gofpdf.Fpdf, nfe *models.NFe) {
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "CÓDIGO DE BARRAS")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 6, nfe.ChaveAcesso)
	pdf.Ln(10)

	// Aqui você poderia adicionar a geração de código de barras real
	// usando uma biblioteca como github.com/boombuler/barcode
}

// formatarValor formata um valor float64 para string
func formatarValor(valor float64) string {
	return strconv.FormatFloat(valor, 'f', 2, 64)
}

// GerarRelatorioBoletos gera um relatório de boletos em PDF
func (s *PDFService) GerarRelatorioBoletos(boletos []models.Boleto) ([]byte, error) {
	s.logger.Info("Gerando relatório de boletos")

	// Cria novo documento PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Configura fonte
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Relatório de Boletos")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, fmt.Sprintf("Total de Boletos: %d", len(boletos)))
	pdf.Ln(10)

	// Cabeçalho da tabela
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(30, 8, "Banco")
	pdf.Cell(40, 8, "Número")
	pdf.Cell(50, 8, "Vencimento")
	pdf.Cell(30, 8, "Valor")
	pdf.Cell(30, 8, "Status")
	pdf.Ln(8)

	// Dados dos boletos
	pdf.SetFont("Arial", "", 10)
	for _, boleto := range boletos {
		pdf.Cell(30, 6, boleto.Banco)
		pdf.Cell(40, 6, boleto.Numero)
		pdf.Cell(50, 6, boleto.Vencimento.Format("02/01/2006"))
		pdf.Cell(30, 6, "R$ "+formatarValor(boleto.Valor))
		pdf.Cell(30, 6, boleto.Status)
		pdf.Ln(6)
	}

	// Gera o PDF
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
