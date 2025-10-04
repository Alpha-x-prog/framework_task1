package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"example/defects/app/backend/internal/auth"
	"example/defects/app/backend/internal/repo"
)

type AuthHandler struct {
	DB     *pgxpool.Pool
	Secret string
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`    // держим нижний регистр в JSON
		Password string `json:"password"` // и в curl/axios
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	u, err := repo.GetUserByEmail(context.Background(), h.DB, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	tok, _ := auth.Sign(u.ID, u.Role, h.Secret, 24*time.Hour)
	c.JSON(http.StatusOK, gin.H{
		"token": tok,
		"user":  gin.H{"id": u.ID, "email": u.Email, "role": u.Role},
	})
}
