package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPagination reads page and pageSize query params, returning sane defaults.
func GetPagination(ctx *gin.Context) (int, int) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	if pageSize > 50 {
		pageSize = 50
	}

	return page, pageSize
}
