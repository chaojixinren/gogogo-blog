package controllers

import (
	"time"

	"gogogo/models"
)

type UserDTO struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email,omitempty"`
	DisplayName string    `json:"displayName"`
	Bio         string    `json:"bio,omitempty"`
	AvatarURL   string    `json:"avatarUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CategoryDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TagDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentDTO struct {
	ID         uint      `json:"id"`
	AuthorName string    `json:"authorName"`
	Body       string    `json:"body"`
	Approved   bool      `json:"approved"`
	CreatedAt  time.Time `json:"createdAt"`
	User       *UserDTO  `json:"user,omitempty"`
}

type PostDTO struct {
	ID          uint         `json:"id"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	Content     string       `json:"content"`
	Slug        string       `json:"slug"`
	Status      string       `json:"status"`
	CoverImage  string       `json:"coverImage,omitempty"`
	PublishedAt *time.Time   `json:"publishedAt,omitempty"`
	Author      UserDTO      `json:"author"`
	Category    *CategoryDTO `json:"category,omitempty"`
	Tags        []TagDTO     `json:"tags"`
	Comments    []CommentDTO `json:"comments,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func buildUserDTO(user models.User) UserDTO {
	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	return UserDTO{
		ID:          user.ID,
		Username:    user.Username,
		Email:       email,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
	}
}

func buildCategoryDTO(category *models.Category) *CategoryDTO {
	if category == nil {
		return nil
	}

	return &CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
	}
}

func buildTagDTOs(tags []models.Tag) []TagDTO {
	result := make([]TagDTO, 0, len(tags))
	for _, tag := range tags {
		result = append(result, TagDTO{
			ID:        tag.ID,
			Name:      tag.Name,
			Slug:      tag.Slug,
			CreatedAt: tag.CreatedAt,
		})
	}
	return result
}

func buildCommentDTOs(comments []models.Comment) []CommentDTO {
	result := make([]CommentDTO, 0, len(comments))
	for _, comment := range comments {
		var userDTO *UserDTO
		if comment.User != nil {
			dto := buildUserDTO(*comment.User)
			dto.Email = ""
			userDTO = &dto
		}

		result = append(result, CommentDTO{
			ID:         comment.ID,
			AuthorName: comment.AuthorName,
			Body:       comment.Body,
			Approved:   comment.Approved,
			CreatedAt:  comment.CreatedAt,
			User:       userDTO,
		})
	}
	return result
}

func buildPostDTO(post models.Post, includeContent bool) PostDTO {
	content := post.Content
	if !includeContent {
		content = ""
	}

	author := buildUserDTO(post.Author)
	author.Email = ""

	dto := PostDTO{
		ID:          post.ID,
		Title:       post.Title,
		Summary:     post.Summary,
		Content:     content,
		Slug:        post.Slug,
		Status:      post.Status,
		CoverImage:  post.CoverImage,
		PublishedAt: post.PublishedAt,
		Author:      author,
		Category:    buildCategoryDTO(post.Category),
		Tags:        buildTagDTOs(post.Tags),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}

	if len(post.Comments) > 0 {
		dto.Comments = buildCommentDTOs(post.Comments)
	}

	return dto
}
