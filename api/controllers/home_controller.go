package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome to my pretty API"})
}
