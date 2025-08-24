package handlers

import (
	"net/http"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ConsultarNFe handler para consultar uma NFe
func ConsultarNFe(nfeService *services.NFEService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithField("handler", "ConsultarNFe")
		
		// Verifica o tipo de conteúdo para determinar como processar a requisição
		contentType := c.GetHeader("Content-Type")
		
		var req models.ConsultaNFeRequest
		var tipoCertificado string
		
		if contentType == "application/json" {
			// Requisição JSON (certificado do sistema)
			var jsonReq struct {
				ChaveAcesso     string `json:"chave_acesso" binding:"required"`
				TipoCertificado string `json:"tipo_certificado"`
			}
			
			if err := c.ShouldBindJSON(&jsonReq); err != nil {
				logger.WithError(err).Error("Erro ao fazer bind dos dados JSON")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Dados inválidos",
					"error":   err.Error(),
				})
				return
			}
			
			req.ChaveAcesso = jsonReq.ChaveAcesso
			tipoCertificado = jsonReq.TipoCertificado
			
		} else {
			// Requisição form-data (certificado de arquivo)
			if err := c.ShouldBind(&req); err != nil {
				logger.WithError(err).Error("Erro ao fazer bind dos dados form")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Dados inválidos",
					"error":   err.Error(),
				})
				return
			}
			
			tipoCertificado = c.PostForm("tipo_certificado")
		}

		// Valida chave de acesso
		if len(req.ChaveAcesso) != 44 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Chave de acesso deve ter 44 dígitos",
			})
			return
		}
		
		logger.WithFields(logrus.Fields{
			"chave_acesso":     req.ChaveAcesso,
			"tipo_certificado": tipoCertificado,
		}).Info("Iniciando consulta de NFe")
		
		// Processa baseado no tipo de certificado
		if tipoCertificado == "sistema" {
			// Usa certificado do navegador
			cert, err := ExtrairCertificadoDaRequisicao(c)
			if err != nil {
				logger.WithError(err).Error("Erro ao extrair certificado da requisição")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Certificado digital não encontrado",
					"error":   "Verifique se o certificado está instalado e o navegador está configurado corretamente.",
				})
				return
			}
			
			// Valida o certificado ICP-Brasil
			if !ValidarCertificadoICP(cert) {
				logger.Error("Certificado não é válido para ICP-Brasil")
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Certificado digital inválido",
					"error":   "Certificado digital não é válido para ICP-Brasil ou está expirado.",
				})
				return
			}
			
			logger.WithField("subject", cert.Subject.String()).Info("Usando certificado do sistema")
		}

		// Consulta NFe
		nfe, err := nfeService.ConsultarNFe(req.ChaveAcesso)
		if err != nil {
			logger.WithError(err).Error("Erro ao consultar NFe")
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao consultar NFe",
				"error":   err.Error(),
			})
			return
		}
		
		logger.Info("NFe consultada com sucesso")

		c.JSON(http.StatusOK, models.ConsultaNFeResponse{
			Success: true,
			Message: "NFe consultada com sucesso",
			Data:    nfe,
		})
	}
}

// BaixarXMLNFe handler para baixar XML da NFe
func BaixarXMLNFe(nfeService *services.NFEService) gin.HandlerFunc {
	return func(c *gin.Context) {
		chave := c.Param("chave")
		
		if len(chave) != 44 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Chave de acesso deve ter 44 dígitos",
			})
			return
		}

		xml, err := nfeService.BaixarXML(chave)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao baixar XML da NFe",
				"error":   err.Error(),
			})
			return
		}

		c.Header("Content-Type", "application/xml")
		c.Header("Content-Disposition", "attachment; filename=nfe_"+chave+".xml")
		c.Data(http.StatusOK, "application/xml", []byte(xml))
	}
}

// GerarPDFNFe handler para gerar PDF da NFe
func GerarPDFNFe(nfeService *services.NFEService, pdfService *services.PDFService) gin.HandlerFunc {
	return func(c *gin.Context) {
		chave := c.Param("chave")
		
		if len(chave) != 44 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Chave de acesso deve ter 44 dígitos",
			})
			return
		}

		// Obtém NFe
		nfe, err := nfeService.ConsultarNFe(chave)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao consultar NFe",
				"error":   err.Error(),
			})
			return
		}

		// Gera PDF
		pdf, err := pdfService.GerarDANFE(nfe)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao gerar PDF",
				"error":   err.Error(),
			})
			return
		}

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=danfe_"+chave+".pdf")
		c.Data(http.StatusOK, "application/pdf", pdf)
	}
}

// ConsultarBoletosNFe handler para consultar boletos de uma NFe
func ConsultarBoletosNFe(nfeService *services.NFEService, bankService *services.BankService) gin.HandlerFunc {
	return func(c *gin.Context) {
		chave := c.Param("chave")
		
		if len(chave) != 44 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Chave de acesso deve ter 44 dígitos",
			})
			return
		}

		// Obtém NFe
		nfe, err := nfeService.ConsultarNFe(chave)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao consultar NFe",
				"error":   err.Error(),
			})
			return
		}

		// Consulta boletos
		boletos, err := bankService.ConsultarBoletosPorNFe(nfe.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao consultar boletos",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Boletos consultados com sucesso",
			"data": gin.H{
				"nfe":     nfe,
				"boletos": boletos,
			},
		})
	}
}
