package response

type (
	PaginationResponse struct {
		Page    int   `json:"page"`
		PerPage int   `json:"per_page"`
		MaxPage int64 `json:"max_page"`
		Count   int64 `json:"count"`
	}
)

func (p *PaginationResponse) GetLimit() int {
	return p.PerPage
}

func (p *PaginationResponse) GetPage() int {
	return p.Page
}
