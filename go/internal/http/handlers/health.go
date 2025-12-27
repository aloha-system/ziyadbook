package handlers

import "github.com/gin-gonic/gin"

type HealthHandler struct{}

func (h HealthHandler) Register(r gin.IRoutes) {
	r.GET("/health", h.get)
}

func (h HealthHandler) get(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
