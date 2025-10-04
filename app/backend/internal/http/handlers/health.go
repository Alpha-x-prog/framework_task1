package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func Me(c *gin.Context) {
	uid, _ := c.Get("uid")
	role, _ := c.Get("role")
	c.JSON(http.StatusOK, gin.H{"uid": uid, "role": role})
}
