package apiv1

// PageReq 公共请求参数
type PageReq struct {
	DateRange []string `p:"dateRange"` //日期范围
	PageIndex int      `p:"pageIndex"` //当前页码
	PageSize  int      `p:"pageSize"`  //每页数
	OrderBy   string   //排序方式
}

// ListRes 列表公共返回
type ListRes struct {
	CurrentPage int `json:"currentPage"`
	Total       int `json:"total"`
}
