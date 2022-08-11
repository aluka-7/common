package common

import (
	"encoding/json"
)

// NoProtectedURIProvider 当前运营子系统中对不需要进行安全拦截的地址提供器，各运营子系统只需要实现该接口并注册到中即可。
type NoProtectedURIProvider interface {
	// Patterns 获取不需要运营安全拦截的URI模式地址，符合ant-style规范。
	Patterns() []string
}

func NewNoProtectedURIProvider(patterns []string) NoProtectedURIProvider {
	return &noProtectedURIProvider{patterns}
}

type noProtectedURIProvider struct {
	patterns []string
}

func (p *noProtectedURIProvider) Patterns() []string {
	return p.patterns
}

// Result 纯状态返回
type Result struct {
	Code    int    `json:"code"` // 返回状态,接口返回码等于0表示SDK错误
	Message string `json:"msg"`  // 返回信息,业务描述
	Opt     string `json:"opt"`  // 操作编号
}

func (r *Result) SetOpt(opt string) {
	r.Opt = opt
}

//DtoResult 返回带有对象信息
type DtoResult struct {
	Result
	Data interface{} `json:"data"` // 具体返回数据的JSON格式
}
type Page struct {
	PageSize     int `json:"pageSize"`     // 每页大小
	TotalPages   int `json:"totalPages"`   // 总页数
	TotalRecords int `json:"totalRecords"` // 总记录数
}

// PageResult 返回带在分页信息
type PageResult struct {
	Result
	Pagination Page        `json:"pagination"`
	Opt        string      `json:"opt"`  // 操作编号
	List       interface{} `json:"data"` // 记录集
}

type Callback struct {
	CbType int    `json:"cbType"` // 业务数据JSON里的参数，回调类型代码
	Opt    string `json:"opt"`    // 操作编号
	Data   string `json:"data"`   // 具体返回数据的JSON格式
}

func NewCallBack(cbType int, v interface{}) Callback {
	data := Json(v, false)
	return Callback{CbType: cbType, Data: data}
}

func Json(v interface{}, indent bool) string {
	var n []byte
	if indent {
		n, _ = json.MarshalIndent(v, "", "\t")
	} else {
		n, _ = json.Marshal(v)
	}
	return string(n)
}
