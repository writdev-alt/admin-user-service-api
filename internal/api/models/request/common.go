package request

type PageInfo struct {
	PageNumber int `form:"pageNumber" json:"pageNumber"`
	PageSize   int `form:"pageSize" json:"pageSize"`
}
