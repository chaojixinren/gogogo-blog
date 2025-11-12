package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gogogo/global"
	"gogogo/models"
	"gogogo/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type postRequest struct {
	Title        string     `json:"title" binding:"required"`
	Summary      string     `json:"summary"`
	Content      string     `json:"content" binding:"required"`
	Status       string     `json:"status"`
	Slug         string     `json:"slug"`
	CoverImage   string     `json:"coverImage"`
	CategoryID   *uint      `json:"categoryId"`
	CategorySlug string     `json:"categorySlug"`
	Tags         []string   `json:"tags"`
	PublishedAt  *time.Time `json:"publishedAt"`
}

type updatePostRequest struct {
	Title        *string    `json:"title"`
	Summary      *string    `json:"summary"`
	Content      *string    `json:"content"`
	Status       *string    `json:"status"`
	Slug         *string    `json:"slug"`
	CoverImage   *string    `json:"coverImage"`
	CategoryID   *uint      `json:"categoryId"`
	CategorySlug *string    `json:"categorySlug"`
	Tags         *[]string  `json:"tags"`
	PublishedAt  *time.Time `json:"publishedAt"`
}

func ListPosts(ctx *gin.Context) {
	listPostsWithScopes(ctx, models.PostStatusPublished)
}

func GetPostByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	post, err := loadPostWithRelations(uint(id))
	if err != nil {
		handlePostLoadError(ctx, err)
		return
	}

	if post.Status != models.PostStatusPublished {
		if userID, ok := optionalUserID(ctx); !ok || userID != post.AuthorID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": buildPostDTO(post, true)})
}

func GetPostBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if strings.TrimSpace(slug) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid slug"})
		return
	}

	var post models.Post
	if err := global.Db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", "approved = ?", true).
		Preload("Comments.User").
		Where("slug = ?", slug).
		First(&post).Error; err != nil {
		handlePostLoadError(ctx, err)
		return
	}

	if post.Status != models.PostStatusPublished {
		if userID, ok := optionalUserID(ctx); !ok || userID != post.AuthorID {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": buildPostDTO(post, true)})
}

func CreatePost(ctx *gin.Context) {
	userID, ok := currentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var input postRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if strings.TrimSpace(input.Content) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "content is required"})
		return
	}

	status := sanitizeStatus(input.Status)
	slug := input.Slug
	if slug == "" {
		slug = utils.Slugify(input.Title)
	}

	var err error
	if slug, err = ensureUniquePostSlug(slug, 0); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate slug"})
		return
	}

	var category *models.Category
	if input.CategoryID != nil || input.CategorySlug != "" {
		category, err = resolveCategory(input.CategoryID, input.CategorySlug)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load category"})
			return
		}
	}

	tags, err := findOrCreateTags(input.Tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process tags"})
		return
	}

	post := models.Post{
		Title:      input.Title,
		Summary:    input.Summary,
		Content:    input.Content,
		Status:     status,
		Slug:       slug,
		CoverImage: input.CoverImage,
		AuthorID:   userID,
		Tags:       tags,
	}

	if category != nil {
		post.CategoryID = &category.ID
	}

	if status == models.PostStatusPublished {
		if input.PublishedAt != nil {
			post.PublishedAt = input.PublishedAt
		} else {
			now := time.Now()
			post.PublishedAt = &now
		}
	}

	if err := global.Db.Create(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	post, err = loadPostWithRelations(post.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load post"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": buildPostDTO(post, true)})
}

func UpdatePost(ctx *gin.Context) {
	post, err := loadPostParam(ctx.Param("id"))
	if err != nil {
		handlePostLoadError(ctx, err)
		return
	}

	userID, ok := currentUserID(ctx)
	if !ok || post.AuthorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "not allowed to edit this post"})
		return
	}

	var input updatePostRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != nil {
		if strings.TrimSpace(*input.Title) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
			return
		}
		post.Title = *input.Title
	}

	if input.Summary != nil {
		post.Summary = *input.Summary
	}

	if input.Content != nil {
		if strings.TrimSpace(*input.Content) == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "content cannot be empty"})
			return
		}
		post.Content = *input.Content
	}

	if input.Status != nil {
		status := sanitizeStatus(*input.Status)
		post.Status = status
		if status == models.PostStatusPublished {
			if input.PublishedAt != nil {
				publishAt := *input.PublishedAt
				post.PublishedAt = &publishAt
			} else if post.PublishedAt == nil {
				now := time.Now()
				post.PublishedAt = &now
			}
		} else if status == models.PostStatusDraft {
			post.PublishedAt = nil
		}
	}

	if input.PublishedAt != nil && post.Status == models.PostStatusPublished {
		publishAt := *input.PublishedAt
		post.PublishedAt = &publishAt
	}

	if input.Slug != nil {
		slug := *input.Slug
		if slug == "" {
			slug = utils.Slugify(post.Title)
		}
		if slug, err = ensureUniquePostSlug(slug, post.ID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update slug"})
			return
		}
		post.Slug = slug
	}

	if input.CoverImage != nil {
		post.CoverImage = *input.CoverImage
	}

	if input.CategoryID != nil {
		if *input.CategoryID == 0 {
			post.CategoryID = nil
		} else {
			category, catErr := resolveCategory(input.CategoryID, "")
			if catErr != nil {
				if errors.Is(catErr, gorm.ErrRecordNotFound) {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
				return
			}
			post.CategoryID = &category.ID
		}
	} else if input.CategorySlug != nil {
		slug := strings.TrimSpace(*input.CategorySlug)
		if slug == "" {
			post.CategoryID = nil
		} else {
			category, catErr := resolveCategory(nil, slug)
			if catErr != nil {
				if errors.Is(catErr, gorm.ErrRecordNotFound) {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
					return
				}
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
				return
			}
			post.CategoryID = &category.ID
		}
	}

	if input.Tags != nil {
		tags, tagErr := findOrCreateTags(*input.Tags)
		if tagErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tags"})
			return
		}
		if err := global.Db.Model(&post).Association("Tags").Replace(tags); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store tags"})
			return
		}
		post.Tags = tags
	}

	if err := global.Db.Save(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	post, err = loadPostWithRelations(post.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to refresh post"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": buildPostDTO(post, true)})
}

func DeletePost(ctx *gin.Context) {
	post, err := loadPostParam(ctx.Param("id"))
	if err != nil {
		handlePostLoadError(ctx, err)
		return
	}

	userID, ok := currentUserID(ctx)
	if !ok || post.AuthorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "not allowed to delete this post"})
		return
	}

	if err := global.Db.Delete(&post).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func ListPostsByCategory(ctx *gin.Context) {
	category, err := loadCategoryParam(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load category"})
		return
	}

	listPostsWithScopes(ctx, models.PostStatusPublished, func(db *gorm.DB) *gorm.DB {
		return db.Where("category_id = ?", category.ID)
	})
}

func ListPostsByTag(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if strings.TrimSpace(slug) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag slug"})
		return
	}

	var tag models.Tag
	if err := global.Db.Where("slug = ?", slug).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tag"})
		return
	}

	listPostsWithScopes(ctx, models.PostStatusPublished, func(db *gorm.DB) *gorm.DB {
		return db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id = ?", tag.ID).
			Distinct("posts.id")
	})
}

func listPostsWithScopes(ctx *gin.Context, defaultStatus string, extraScopes ...func(*gorm.DB) *gorm.DB) {
	page, pageSize := utils.GetPagination(ctx)
	includeContent := ctx.DefaultQuery("includeContent", "false") == "true"

	allScopes := append(extraScopes, postFilters(ctx, defaultStatus))

	countDB := global.Db.Model(&models.Post{})
	for _, scope := range allScopes {
		countDB = scope(countDB)
	}

	var total int64
	if err := countDB.Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count posts"})
		return
	}

	queryDB := global.Db.Model(&models.Post{})
	for _, scope := range allScopes {
		queryDB = scope(queryDB)
	}

	var posts []models.Post
	if err := queryDB.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC, created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load posts"})
		return
	}

	response := make([]PostDTO, 0, len(posts))
	for _, post := range posts {
		response = append(response, buildPostDTO(post, includeContent))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":     response,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

func postFilters(ctx *gin.Context, defaultStatus string) func(*gorm.DB) *gorm.DB {
	status := ctx.DefaultQuery("status", defaultStatus)
	category := ctx.Query("category")
	tag := ctx.Query("tag")
	author := ctx.Query("author")
	search := ctx.Query("search")

	return func(db *gorm.DB) *gorm.DB {
		if status != "" && status != "all" {
			db = db.Where("posts.status = ?", status)
		}

		if category != "" {
			db = db.Joins("JOIN categories ON categories.id = posts.category_id").
				Where("categories.slug = ?", category).
				Distinct("posts.id")
		}

		if tag != "" {
			db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
				Joins("JOIN tags ON tags.id = post_tags.tag_id").
				Where("tags.slug = ?", tag).
				Distinct("posts.id")
		}

		if author != "" {
			db = db.Joins("JOIN users ON users.id = posts.author_id").
				Where("users.username = ?", author).
				Distinct("posts.id")
		}

		if search != "" {
			like := "%" + search + "%"
			db = db.Where("posts.title LIKE ? OR posts.summary LIKE ?", like, like)
		}

		return db
	}
}

func findOrCreateTags(values []string) ([]models.Tag, error) {
	result := make([]models.Tag, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		slug := utils.Slugify(value)
		var tag models.Tag
		err := global.Db.Where("slug = ?", slug).First(&tag).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tag = models.Tag{Name: value, Slug: slug}
				if err := global.Db.Create(&tag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		result = append(result, tag)
	}
	return result, nil
}

func sanitizeStatus(status string) string {
	switch strings.ToLower(status) {
	case models.PostStatusPublished:
		return models.PostStatusPublished
	case models.PostStatusArchived:
		return models.PostStatusArchived
	default:
		return models.PostStatusDraft
	}
}

func ensureUniquePostSlug(slug string, excludeID uint) (string, error) {
	base := utils.Slugify(slug)
	if base == "" {
		base = "post"
	}

	candidate := base
	var count int64
	for i := 0; ; i++ {
		query := global.Db.Model(&models.Post{}).Where("slug = ?", candidate)
		if excludeID > 0 {
			query = query.Where("id <> ?", excludeID)
		}
		if err := query.Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i+1)
	}
}

func resolveCategory(id *uint, slug string) (*models.Category, error) {
	if id == nil && strings.TrimSpace(slug) == "" {
		return nil, nil
	}

	var category models.Category
	var err error
	if id != nil && *id > 0 {
		err = global.Db.First(&category, *id).Error
	} else {
		err = global.Db.Where("slug = ?", slug).First(&category).Error
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func loadPostWithRelations(id uint) (models.Post, error) {
	var post models.Post
	err := global.Db.
		Preload("Author").
		Preload("Category").
		Preload("Tags").
		Preload("Comments", "approved = ?", true).
		Preload("Comments.User").
		First(&post, id).Error
	return post, err
}

func loadPostParam(idParam string) (models.Post, error) {
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return models.Post{}, gorm.ErrRecordNotFound
	}
	return loadPostWithRelations(uint(id))
}

func handlePostLoadError(ctx *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load post"})
}

func loadCategoryParam(param string) (models.Category, error) {
	param = strings.TrimSpace(param)
	var category models.Category
	if param == "" {
		return category, gorm.ErrRecordNotFound
	}

	if id, err := strconv.ParseUint(param, 10, 64); err == nil {
		if err := global.Db.First(&category, id).Error; err != nil {
			return category, err
		}
		return category, nil
	}

	if err := global.Db.Where("slug = ?", param).First(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}
