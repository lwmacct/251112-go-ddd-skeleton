package pagination

import "math"

// Pagination 分页结构
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewPagination 创建分页对象
func NewPagination(page, pageSize int, total int64) *Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 最大每页100条
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
	return p.PageSize
}

// HasNext 是否有下一页
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

// HasPrev 是否有上一页
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// PaginatedResponse 分页响应结构
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

// NewPaginatedResponse 创建分页响应
func NewPaginatedResponse(data interface{}, pagination *Pagination) *PaginatedResponse {
	return &PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}
}

// ParsePaginationParams 解析分页参数
func ParsePaginationParams(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return
}
