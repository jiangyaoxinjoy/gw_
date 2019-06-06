package model

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
}
