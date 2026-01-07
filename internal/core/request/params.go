package request

type PaginationParams struct {
	Page     int `form:"page,omitempty" json:"page,omitempty"`
	PageSize int `form:"page_size,omitempty" json:"pageSize,omitempty"`
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}
