package controllers

import (
	"strings"

	"gogogo/global"
	"gogogo/models"
	"gogogo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func currentUserID(ctx *gin.Context) (uint, bool) {
	idValue, exists := ctx.Get("userID")
	if !exists {
		return 0, false
	}

	switch id := idValue.(type) {
	case uint:
		return id, true
	case int:
		return uint(id), true
	case int64:
		return uint(id), true
	case float64:
		return uint(id), true
	default:
		return 0, false
	}
}

func optionalUserID(ctx *gin.Context) (uint, bool) {
	if id, ok := currentUserID(ctx); ok {
		return id, true
	}

	rawHeader := ctx.GetHeader("Authorization")
	if rawHeader == "" {
		return 0, false
	}

	token := rawHeader
	if parts := strings.SplitN(rawHeader, " ", 2); len(parts) == 2 {
		token = parts[1]
	}

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		return 0, false
	}

	ctx.Set("userID", claims.UserID)
	ctx.Set("username", claims.Username)
	return claims.UserID, true
}

func loadUserByID(_ *gin.Context, id uint) (*models.User, error) {
	var user models.User
	if err := global.Db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
