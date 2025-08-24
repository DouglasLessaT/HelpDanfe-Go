package handlers

import (
	"fmt"
	"net/http"

	"github.com/Douglaslessat/HelpDanfe-Go/internal/models"
	"github.com/Douglaslessat/HelpDanfe-Go/internal/services"

	"github.com/gin-gonic/gin"
)

// ConsultarBoleto handler para consultar um boleto específico
func ConsultarBoleto(bankService *services.BankService) gin.HandlerFunc {
	return func(c *gin.Context) {
		codigo := c.Param("codigo")
		
		if codigo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Código do boleto é obrigatório",
			})
			return
		}

		// Consulta boleto
		boleto, err := bankService.ConsultarBoleto(codigo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao consultar boleto",
				"error":   err.Error(),
			})
			return
		}

		if boleto == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Boleto não encontrado",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Boleto consultado com sucesso",
			"data":    boleto,
		})
	}
}

// ConsultarMultiplosBoletos handler para consultar múltiplos boletos
func ConsultarMultiplosBoletos(bankService *services.BankService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ConsultaBoletosRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Dados inválidos",
				"error":   err.Error(),
			})
			return
		}

		if len(req.Codigos) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Pelo menos um código de boleto deve ser fornecido",
			})
			return
		}

		// Consulta múltiplos boletos
		boletos, err := bankService.ConsultarMultiplosBoletos(req.Codigos)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao consultar boletos",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, models.ConsultaBoletosResponse{
			Success: true,
			Message: "Boletos consultados com sucesso",
			Data:    boletos,
		})
	}
}

// ConsultarBoletosPorDuplicata handler para consultar boletos por duplicata
func ConsultarBoletosPorDuplicata(bankService *services.BankService) gin.HandlerFunc {
	return func(c *gin.Context) {
		duplicataID := c.Param("duplicata_id")
		
		if duplicataID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID da duplicata é obrigatório",
			})
			return
		}

		// Converte string para uint
		var id uint
		if _, err := fmt.Sscanf(duplicataID, "%d", &id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID da duplicata inválido",
			})
			return
		}

		// Consulta boletos por duplicata
		boletos, err := bankService.ConsultarBoletosPorDuplicata(id)
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
			"data":    boletos,
		})
	}
}

// GerarCodigoBarras handler para gerar código de barras de um boleto
func GerarCodigoBarras(bankService *services.BankService) gin.HandlerFunc {
	return func(c *gin.Context) {
		codigo := c.Param("codigo")
		
		if codigo == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Código do boleto é obrigatório",
			})
			return
		}

		// Gera código de barras
		codigoBarras, err := bankService.GerarCodigoBarras(codigo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Erro ao gerar código de barras",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Código de barras gerado com sucesso",
			"data": gin.H{
				"codigo_barras": codigoBarras,
			},
		})
	}
}
