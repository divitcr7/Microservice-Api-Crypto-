package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/streamdp/ccd/config"
)

// SendHTML show a beautiful page with small intro and instruction
func SendHTML(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"year":    time.Now().Year(),
		"version": config.Version,
	})
}

// SendOK using for HEAD request and send 200 and nil body
func SendOK(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
