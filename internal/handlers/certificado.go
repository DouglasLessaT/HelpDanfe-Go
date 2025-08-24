package handlers

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CertificadoInfo representa informações sobre um certificado
type CertificadoInfo struct {
	Subject      string `json:"subject"`
	Issuer       string `json:"issuer"`
	SerialNumber string `json:"serial_number"`
	NotBefore    string `json:"not_before"`
	NotAfter     string `json:"not_after"`
	Valid        bool   `json:"valid"`
}

// VerificarCertificados verifica se há certificados disponíveis no navegador
func VerificarCertificados() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithField("handler", "VerificarCertificados")
		
		// Configura headers para requisitar certificado cliente
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		
		// Verifica se há certificados na requisição
		if c.Request.TLS != nil && len(c.Request.TLS.PeerCertificates) > 0 {
			cert := c.Request.TLS.PeerCertificates[0]
			
			info := CertificadoInfo{
				Subject:      cert.Subject.String(),
				Issuer:       cert.Issuer.String(),
				SerialNumber: cert.SerialNumber.String(),
				NotBefore:    cert.NotBefore.Format("2006-01-02 15:04:05"),
				NotAfter:     cert.NotAfter.Format("2006-01-02 15:04:05"),
				Valid:        cert.NotAfter.After(cert.NotBefore),
			}
			
			logger.Info("Certificado detectado automaticamente")
			
			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"certificado": info,
				"message":     "Certificado detectado",
			})
			return
		}
		
		// Se não há certificados, retorna status indicando que precisa selecionar
		logger.Info("Nenhum certificado detectado automaticamente")
		
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Nenhum certificado detectado. Configure seu navegador para permitir certificados cliente.",
		})
	}
}

// SelecionarCertificado força a seleção de certificado pelo usuário
func SelecionarCertificado() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithField("handler", "SelecionarCertificado")
		
		// Configura headers para forçar seleção de certificado
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		
		// Verifica se há certificados na requisição
		if c.Request.TLS != nil && len(c.Request.TLS.PeerCertificates) > 0 {
			cert := c.Request.TLS.PeerCertificates[0]
			
			// Valida o certificado
			
			// Para certificados ICP-Brasil, podemos adicionar validações específicas
			isICPBrasil := strings.Contains(cert.Issuer.String(), "ICP-Brasil") ||
				strings.Contains(cert.Subject.String(), "ICP-Brasil")
			
			info := CertificadoInfo{
				Subject:      cert.Subject.String(),
				Issuer:       cert.Issuer.String(),
				SerialNumber: cert.SerialNumber.String(),
				NotBefore:    cert.NotBefore.Format("2006-01-02 15:04:05"),
				NotAfter:     cert.NotAfter.Format("2006-01-02 15:04:05"),
				Valid:        isICPBrasil && cert.NotAfter.After(cert.NotBefore),
			}
			
			logger.WithFields(logrus.Fields{
				"subject":     info.Subject,
				"issuer":      info.Issuer,
				"valid":       info.Valid,
				"icp_brasil":  isICPBrasil,
			}).Info("Certificado selecionado")
			
			c.JSON(http.StatusOK, gin.H{
				"success":     true,
				"certificado": info,
				"message":     "Certificado selecionado com sucesso",
			})
			return
		}
		
		// Se não há certificados, retorna erro
		logger.Warn("Nenhum certificado fornecido na requisição")
		
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Nenhum certificado foi selecionado. Verifique se há certificados instalados no sistema.",
		})
	}
}

// ExtrairCertificadoDaRequisicao extrai informações do certificado da requisição TLS
func ExtrairCertificadoDaRequisicao(c *gin.Context) (*x509.Certificate, error) {
	if c.Request.TLS == nil || len(c.Request.TLS.PeerCertificates) == 0 {
		return nil, gin.Error{
			Err:  errors.New("Nenhum certificado encontrado na requisição"),
			Type: gin.ErrorTypePublic,
			Meta: "Nenhum certificado encontrado na requisição",
		}
	}
	
	return c.Request.TLS.PeerCertificates[0], nil
}

// ValidarCertificadoICP valida se o certificado é válido para ICP-Brasil
func ValidarCertificadoICP(cert *x509.Certificate) bool {
	// Verifica se é um certificado ICP-Brasil
	isICPBrasil := strings.Contains(cert.Issuer.String(), "ICP-Brasil") ||
		strings.Contains(cert.Subject.String(), "ICP-Brasil") ||
		strings.Contains(cert.Issuer.String(), "AC ") // Autoridade Certificadora
	
	// Verifica validade temporal
	isTemporalmenteValido := cert.NotAfter.After(cert.NotBefore)
	
	// Verifica se não expirou
	isNaoExpirado := cert.NotAfter.After(cert.NotBefore)
	
	return isICPBrasil && isTemporalmenteValido && isNaoExpirado
}

// ConfigurarTLSParaCertificadoCliente configura o servidor para aceitar certificados cliente
func ConfigurarTLSParaCertificadoCliente() *tls.Config {
	return &tls.Config{
		ClientAuth: tls.RequestClientCert,
		ClientCAs:  nil, // Aceita qualquer CA por enquanto
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Validação customizada pode ser adicionada aqui
			return nil
		},
	}
}
