package wrapx

import (
	"context"
	"reflect"

	"github.com/gin-gonic/gin"
)

// NewAPIFunc Custom context support
// type NewAPIFunc func(*gin.Context) interface{}

type GinBeforeAfter interface {
	GinBefore(req *GinBeforeAfterInfo) bool
	GinAfter(req *GinBeforeAfterInfo) bool
}

type _Base struct {
	isBigCamel  bool // big camel style.大驼峰命名规则
	isDev       bool // if is development
	apiFun      NewAPIFunc2
	apiType     reflect.Type
	outPath     string // output path.输出目录
	beforeAfter GinBeforeAfter
	isOutDoc    bool
	prefix      string
	service     string
}

// GinBeforeAfterInfo 对象调用前后执行中间件参数
type GinBeforeAfterInfo struct {
	C        *gin.Context
	FuncName string      // 函数名
	Req      interface{} // 调用前的请求参数
	Resp     interface{} // 调用后的返回参数
	Error    error
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context // 占位参数，可用于存储其他参数，前后连接可用

}

// type model struct {
// 	Group string // group 标记
// 	MP    map[string]map[string]DocModel
// }

// store the comment for the controller method. 生成注解路由
type genComment struct {
	RouterPath string
	Note       string // 注释
	Methods    []string
}

// router style list.路由规则列表
type genRouterInfo struct {
	GenComment  genComment
	HandFunName string
}

type ElementInfo struct {
	Name string // 参数名
	// URL      string      // web 访问参数
	Tag      string      // 标签
	Type     string      // 类型
	TypeRef  *StructInfo // 类型定义
	IsArray  bool        // 是否是数组
	Requierd bool        // 是否必须
	Note     string      // 注释
	Default  string      // 默认值
}

// StructInfo struct define
type StructInfo struct {
	Items []ElementInfo // 结构体元素
	Note  string        // 注释
	Name  string        //结构体名字
	Pkg   string        // 包名
}

// DocModel model
type DocModel struct {
	RouterPath string
	Methods    []string
	Note       string
	Req, Resp  *StructInfo
}

type MessageBody struct {
	State bool        `json:"state"`
	Code  int         `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type genInfo struct {
	List []genRouterInfo
	Tm   int64 //genout time
}
