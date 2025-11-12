package controllers

import (
	"net/http"

	"gogogo/global"
	"gogogo/models"
	"gogogo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProfile(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var user models.User
	if err := global.Db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": buildUserDTO(user)})
}

func ListMyPosts(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	page, pageSize := utils.GetPagination(ctx)

	var total int64
	if err := global.Db.Model(&models.Post{}).
		Where("author_id = ?", userID).
		Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count posts"})
		return
	}

	var posts []models.Post
	if err := global.Db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Where("author_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load posts"})
		return
	}

	response := make([]PostDTO, 0, len(posts))
	for _, post := range posts {
		response = append(response, buildPostDTO(post, false))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     response,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}
