package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 10
)

// Pagination query params
type PaginationQuery struct {
	Limit   int    `json:"limit,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

// Set page limit
func (q *PaginationQuery) SetLimit(limitQuery string) error {
	if limitQuery == "" {
		q.Limit = defaultLimit
		return nil
	}
	n, err := strconv.Atoi(limitQuery)
	if err != nil {
		return err
	}
	if n > 500 {
		n = 500
	}
	q.Limit = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Limit = 0

		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// Set order by
func (q *PaginationQuery) SetOrderBy(orderByQuery string) {
	q.OrderBy = orderByQuery
}

// Get offset
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Limit
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Limit
}

// Get OrderBy
func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}

// Get OrderBy
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&limit=%v&orderBy=%s", q.GetPage(), q.GetLimit(), q.GetOrderBy())
}

// Get pagination query struct from
func GetPaginationFromCtx(c *gin.Context) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.DefaultQuery("Page", "1")); err != nil {
		return nil, err
	}
	if err := q.SetLimit(c.DefaultQuery("Limit", "50")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.Query("OrderBy"))

	return q, nil
}

// Get total pages int
func GetTotalPages(totalCount int64, pageLimit int) int {
	d := float64(totalCount) / float64(pageLimit)
	return int(math.Ceil(d))
}

// Get has more
// func GetHasMore(currentPage int, totalCount int64, pageLimit int) bool {
// 	return currentPage < totalCount/pageLimit
// }

// GetHasMore checks if there are more pages available
func GetHasMore(currentPage int, totalCount int64, pageLimit int) bool {
	totalPages := totalCount / int64(pageLimit)
	if totalCount%int64(pageLimit) != 0 {
		totalPages++
	}
	return currentPage < int(totalPages)
}
