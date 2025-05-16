package controllers

import (
	"encoding/json"
	"github.com/ybazli/auth-service/src/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ybazli/auth-service/src/config"
	"github.com/ybazli/auth-service/src/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful. Please login."})
}

func Login(c *gin.Context) {
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokens, err := utils.GenerateTokenPair(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	sessionData := models.Session{
		UserID:       user.ID,
		RefreshToken: tokens.RefreshToken,
		UserAgent:    c.Request.UserAgent(),
		IP:           c.ClientIP(),
		CreatedAt:    time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(sessionData)

	// store refresh token in redis
	err = config.RedisClient.Set(
		config.Ctx,
		utils.SessionKey(tokens.SessionID),
		data,
		7*24*time.Hour,
	).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"session_id":    tokens.SessionID,
	})
}

func RefreshToken(c *gin.Context) {
	type RefreshToken struct {
		SessionID    string `json:"session_id" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var input RefreshToken

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	sessionData, err := config.RedisClient.Get(config.Ctx, utils.SessionKey(input.SessionID)).Result()
	var session models.Session
	err = json.Unmarshal([]byte(sessionData), &session)

	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Session not found or expired")
		return
	}

	if session.RefreshToken != input.RefreshToken {
		utils.Error(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	var user models.User
	if err := config.DB.First(&user, session.UserID).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "User not found.")
	}

	newTokens, err := utils.GenerateTokenPair(user)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token.")
		return
	}
	newSession := models.Session{
		UserID:       user.ID,
		RefreshToken: newTokens.RefreshToken,
		UserAgent:    c.Request.UserAgent(),
		IP:           c.ClientIP(),
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	newSessionData, _ := json.Marshal(newSession)
	config.RedisClient.Del(config.Ctx, utils.SessionKey(input.SessionID))

	config.RedisClient.Set(config.Ctx, utils.SessionKey(newTokens.SessionID), newSessionData, 7*24*time.Hour)

	utils.Success(c, gin.H{
		"access_token":  newTokens.AccessToken,
		"refresh_token": newTokens.RefreshToken,
		"session_id":    newTokens.SessionID,
	})
}
