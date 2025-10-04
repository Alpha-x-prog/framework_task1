package handlers

import (
	"context"
	"net/http"
	"strings"
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

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`    // обяз.
		Password string `json:"password"` // обяз., >= 6
		Role     string `json:"role"`     // опц., по умолч. "engineer"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}

	// простая валидация
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if !strings.Contains(req.Email, "@") || len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password too short"})
		return
	}
	if req.Role == "" {
		req.Role = "engineer"
	} // дефолтная роль

	// уже существует?
	exists, err := repo.EmailExists(c, h.DB, req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// найти role_id
	roleID, err := repo.GetRoleIDByName(c, h.DB, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown role"})
		return
	}

	// хеш пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// создать пользователя
	uid, err := repo.CreateUser(c, h.DB, req.Email, string(hash), roleID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// сразу выдать токен (автовход после регистрации)
	tok, _ := auth.Sign(uid, req.Role, h.Secret, 24*time.Hour)
	c.JSON(http.StatusCreated, gin.H{
		"token": tok,
		"user":  gin.H{"id": uid, "email": req.Email, "role": req.Role},
	})
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
