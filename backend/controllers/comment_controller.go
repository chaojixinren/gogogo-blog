package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gogogo/global"
	"gogogo/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type commentRequest struct {
	AuthorName string `json:"authorName"`
	Body       string `json:"body" binding:"required"`
}

func ListComments(ctx *gin.Context) {
	post, err := loadPostSummary(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load post"})
		return
	}

	query := global.Db.Where("post_id = ?", post.ID).Order("created_at ASC").Preload("User")
	if userID, ok := optionalUserID(ctx); !ok || userID != post.AuthorID {
		query = query.Where("approved = ?", true)
	}

	var comments []models.Comment
	if err := query.Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load comments"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": buildCommentDTOs(comments)})
}

func CreateComment(ctx *gin.Context) {
	post, err := loadPostSummary(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load post"})
		return
	}

	var input commentRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.AuthorName = strings.TrimSpace(input.AuthorName)
	if strings.TrimSpace(input.Body) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "comment body is required"})
		return
	}

	var comment models.Comment
	comment.PostID = post.ID
	comment.Body = input.Body

	if userID, ok := optionalUserID(ctx); ok {
		user, userErr := loadUserByID(ctx, userID)
		if userErr != nil {
			if errors.Is(userErr, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			return
		}
		comment.UserID = &user.ID
		comment.AuthorName = user.DisplayName
		if comment.AuthorName == "" {
			comment.AuthorName = user.Username
		}
		comment.Approved = true
	} else {
		if input.AuthorName == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "author name is required for guest comments"})
			return
		}
		comment.AuthorName = input.AuthorName
		comment.Approved = true
	}

	if err := global.Db.Create(&comment).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	if err := global.Db.Preload("User").First(&comment, comment.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load comment"})
		return
	}

	dto := buildCommentDTOs([]models.Comment{comment})
	ctx.JSON(http.StatusCreated, gin.H{"data": dto[0]})
}

func loadPostSummary(param string) (models.Post, error) {
	param = strings.TrimSpace(param)
	var post models.Post
	if param == "" {
		return post, gorm.ErrRecordNotFound
	}

	if id, err := strconv.ParseUint(param, 10, 64); err == nil {
		err = global.Db.Select("id", "author_id").First(&post, id).Error
		return post, err
	}

	err := global.Db.Select("id", "author_id").Where("slug = ?", param).First(&post).Error
	return post, err
}
