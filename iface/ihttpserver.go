package iface

import "net/http"

// HttpCodeMsgHandler http内部错误码转外部错误码和文本
type HttpCodeMsgHandler func(code uint32, inner bool) (uint32, string)

type HttpServerResType bool

type IHttpServer interface {
	Init()
	Run() error
	Stop() error
	SetSSLFile(certFile, keyFile string)
	// RespJson 返回json数据
	RespJson(w http.ResponseWriter, data interface{}) HttpServerResType
	// RespErr 返回错误码
	RespErr(w http.ResponseWriter, errCode uint32) HttpServerResType
	// RespCustom 自定义返回
	RespCustom(err error) HttpServerResType
}
