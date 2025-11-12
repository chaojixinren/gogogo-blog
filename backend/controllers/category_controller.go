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

type categoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type updateCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}

func ListCategories(ctx *gin.Context) {
	var categories []models.Category
	if err := global.Db.Order("name ASC").Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load categories"})
		return
	}

	result := make([]CategoryDTO, 0, len(categories))
	for _, category := range categories {
		if dto := buildCategoryDTO(&category); dto != nil {
			result = append(result, *dto)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

func CreateCategory(ctx *gin.Context) {
	var input categoryRequest
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

	if err := ensureUniqueCategorySlug(slug, 0); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:        input.Name,
		Slug:        slug,
		Description: input.Description,
	}

	if err := global.Db.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": buildCategoryDTO(&category)})
}

func UpdateCategory(ctx *gin.Context) {
	category, err := loadCategoryParam(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load category"})
		return
	}

	var input updateCategoryRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
			return
		}
		category.Name = *input.Name
	}

	if input.Slug != nil {
		slug := *input.Slug
		if slug == "" {
			slug = utils.Slugify(category.Name)
		}
		if err := ensureUniqueCategorySlug(slug, category.ID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		category.Slug = slug
	}

	if input.Description != nil {
		category.Description = *input.Description
	}

	if err := global.Db.Save(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": buildCategoryDTO(&category)})
}

func DeleteCategory(ctx *gin.Context) {
	category, err := loadCategoryParam(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load category"})
		return
	}

	if err := global.Db.Model(&models.Post{}).Where("category_id = ?", category.ID).Update("category_id", nil).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to detach posts"})
		return
	}

	if err := global.Db.Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ensureUniqueCategorySlug(slug string, excludeID uint) error {
	slug = utils.Slugify(slug)
	if slug == "" {
		return errors.New("invalid slug")
	}

	var count int64
	query := global.Db.Model(&models.Category{}).Where("slug = ?", slug)
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
