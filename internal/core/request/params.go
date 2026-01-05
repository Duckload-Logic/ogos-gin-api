package request

type PaginationParams struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`
}

func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}
