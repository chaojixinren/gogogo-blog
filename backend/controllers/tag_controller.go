package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gogogo/global"
	"gogogo/models"
	"gogogo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

type updateTagRequest struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}

func ListTags(ctx *gin.Context) {
	var tags []models.Tag
	if err := global.Db.Order("name ASC").Find(&tags).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tags"})
		return
	}

	result := make([]TagDTO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, TagDTO{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func CreateTag(ctx *gin.Context) {
	var input tagRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(input.Name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	slug := input.Slug
	if slug == "" {
		slug = utils.Slugify(input.Name)
	}

	if err := ensureUniqueTagSlug(slug, 0); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := models.Tag{
		Name: input.Name,
		Slug: slug,
	}

	if err := global.Db.Create(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tag"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": TagDTO{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}})
}

func UpdateTag(ctx *gin.Context) {
	tag, err := loadTagParam(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tag"})
		return
	}

	var input updateTagRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
			return
		}
		tag.Name = *input.Name
	}

	if input.Slug != nil {
		slug := *input.Slug
		if slug == "" {
			slug = utils.Slugify(tag.Name)
		}
		if err := ensureUniqueTagSlug(slug, tag.ID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tag.Slug = slug
	}

	if err := global.Db.Save(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tag"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": TagDTO{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}})
}

func DeleteTag(ctx *gin.Context) {
	tag, err := loadTagParam(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tag"})
		return
	}

	if err := global.Db.Model(&tag).Association("Posts").Clear(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to detach tag"})
		return
	}

	if err := global.Db.Delete(&tag).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tag"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ensureUniqueTagSlug(slug string, excludeID uint) error {
	slug = utils.Slugify(slug)
	if slug == "" {
		return errors.New("invalid slug")
	}

	var count int64
	query := global.Db.Model(&models.Tag{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("slug already in use")
	}
	return nil
}

func loadTagParam(param string) (models.Tag, error) {
	param = strings.TrimSpace(param)
	var tag models.Tag
	if param == "" {
		return tag, gorm.ErrRecordNotFound
	}

	if id, err := strconv.ParseUint(param, 10, 64); err == nil {
		if err := global.Db.First(&tag, id).Error; err != nil {
			return tag, err
		}
		return tag, nil
	}

	if err := global.Db.Where("slug = ?", param).First(&tag).Error; err != nil {
		return tag, err
	}
	return tag, nil
}
