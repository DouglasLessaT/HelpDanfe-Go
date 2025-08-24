package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServeStatic serve arquivos estáticos do frontend
func ServeStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Configura o diretório de arquivos estáticos
		fs := http.FileServer(http.Dir("./web"))
		
		// Se a requisição for para a raiz, serve o index.html
		if c.Request.URL.Path == "/" {
			c.File("./web/index.html")
			return
		}
		
		// Para outras requisições, serve os arquivos estáticos
		fs.ServeHTTP(c.Writer, c.Request)
	}
}

// ServeIndex serve o arquivo index.html para rotas SPA
func ServeIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verifica se é uma requisição para arquivo estático
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || 
		   strings.HasPrefix(path, "/css/") || 
		   strings.HasPrefix(path, "/js/") || 
		   strings.HasPrefix(path, "/fonts/") ||
		   strings.Contains(path, ".") {
			c.Next()
			return
		}
		
		// Para outras rotas, serve o index.html (SPA)
		c.File("./web/index.html")
	}
}
