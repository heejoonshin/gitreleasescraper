package response

type Pagination struct {
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
	Page     int `json:"page"`
}
