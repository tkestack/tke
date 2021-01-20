package proxy

import (
	gojson "encoding/json"
)

type BaseApiResponse interface {
	Code() int
	Message() string
	RawData() *gojson.RawMessage
	SetData(data interface{})
}

type Response struct {
	// 错误码
	ResultCode int `json:"result_code,omitempty"`

	// 错误信息
	Msg string `json:"message,omitempty"`

	// 返回结构体，可选
	Data interface{} `json:"data,omitempty"`
}

type Page struct {
	Items      interface{} `json:"items"`
	TotalCount int         `json:"total_count"`
}

var _ BaseApiResponse = ApiResponse{}

// var _ BaseApiResponse = ApiPageResponse{}

// ApiResponse - API 响应结构体
type ApiResponse struct {
	// 错误码
	ResultCode int `json:"result_code,omitempty"`

	// 错误信息
	Msg string `json:"message,omitempty"`

	// 返回结构体，可选
	Raw *gojson.RawMessage `json:"data,omitempty"`

	Data interface{} `json:"-"`
}

type ApiError struct {
	Response `json:",inline"`
	Status   int `json:"-"`
}

func (a ApiResponse) Code() int {
	return a.ResultCode
}

func (a ApiResponse) Message() string {
	return a.Msg
}

func (a ApiResponse) RawData() *gojson.RawMessage {
	return a.Raw
}

func (a ApiResponse) SetData(data interface{}) {
	a.Data = data
}

type ApiPageResponse struct {
	// 错误码
	ResultCode int `json:"result_code,omitempty"`

	// 错误信息
	Msg string `json:"message,omitempty"`

	Data ApiPage `json:"data"`
}

func (p ApiPageResponse) Code() int {
	return p.ResultCode
}

func (p ApiPageResponse) Message() string {
	return p.Msg
}

func (p ApiPageResponse) RawData() *gojson.RawMessage {
	return p.Data.Items
}

func (p ApiPageResponse) SetData(data interface{}) {
	p.Data = data.(ApiPage)
}

type ApiPage struct {
	Items      *gojson.RawMessage `json:"items,omitempty"`
	TotalCount int                `json:"total_count"`
}

type PageItems []interface{}

func (p PageItems) Filter(filter func(interface{}) bool) PageItems {
	results := make([]interface{}, 0)
	for i, l := 0, len(p); i < l; i++ {
		item := p[i]
		if filter(item) {
			results = append(results, item)
		}
	}

	return results
}

func (p PageItems) First(filter func(interface{}) bool) interface{} {
	var result interface{} = nil
	for i, l := 0, len(p); i < l; i++ {
		item := p[i]
		if filter(item) {
			result = item
			break
		}
	}
	return result
}

func (p PageItems) Paging(offset, limit int) (int, PageItems) {
	total := len(p)

	if offset > total {
		offset = total
	} else if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit >= 200 {
		limit = 10
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return total, p[offset:end]
}

type PageQuery struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
}


type TCMErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}