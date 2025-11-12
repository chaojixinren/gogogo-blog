package controllers

import (
	"errors"
	"net/http"
	"strings"

	"gogogo/global"
	"gogogo/models"
	"gogogo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type registerRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"displayName"`
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(ctx *gin.Context) {
	var input registerRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)

	if len(input.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
		return
	}

	var existing models.User
	if err := global.Db.Where("username = ?", input.Username).First(&existing).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate username"})
		return
	}

	if input.Email != "" {
		if err := global.Db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate email"})
			return
		}
	}

	hashedPwd, err := utils.HashPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	var emailPtr *string
	if input.Email != "" {
		email := input.Email
		emailPtr = &email
	}

	user := models.User{
		Username:    input.Username,
		Email:       emailPtr,
		DisplayName: input.DisplayName,
		Password:    hashedPwd,
	}

	if strings.TrimSpace(user.DisplayName) == "" {
		user.DisplayName = user.Username
	}

	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  buildUserDTO(user),
	})
}

func Login(ctx *gin.Context) {
	var input loginRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  buildUserDTO(user),
	})
}
