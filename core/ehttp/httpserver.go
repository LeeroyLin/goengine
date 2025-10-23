package ehttp

import (
	"encoding/json"
	"fmt"
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/iface"
	"net/http"
	"strings"
)

type WebRespData struct {
	ErrCode uint32      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type HttpServer struct {
	Server             *http.Server
	mux                *http.ServeMux
	IsHttps            bool
	IP                 string
	Port               int
	addr               string
	certFile           string
	keyFile            string
	closeChan          chan interface{}
	httpCodeMsgHandler iface.HttpCodeMsgHandler
}

func NewHttpServer(ip string, port int, isHttps bool, httpCodeMsgHandler iface.HttpCodeMsgHandler) *HttpServer {
	addr := fmt.Sprintf("%s:%d", ip, port)

	s := &HttpServer{
		Server: &http.Server{
			Addr: addr,
		},
		IP:                 ip,
		Port:               port,
		IsHttps:            isHttps,
		closeChan:          make(chan interface{}),
		httpCodeMsgHandler: httpCodeMsgHandler,
	}

	s.addr = fmt.Sprintf("%s:%d", ip, port)

	return s
}

func (s *HttpServer) Init() {
	s.mux = http.NewServeMux()
	s.Server.Handler = s.mux
}

func (s *HttpServer) Run() {
	s.closeChan = make(chan interface{})

	if s.IsHttps {
		s.runAsHttps()
	} else {
		s.runAsHttp()
	}

	select {
	case <-s.closeChan:
		return
	}
}

func (s *HttpServer) Stop() {
	select {
	case <-s.closeChan:
		return
	default:
		close(s.closeChan)

		err := s.Server.Close()
		if err != nil {
			elog.Fatal("[HttpServer] close err:", s.addr, err)
			return
		}

		elog.Fatal("[HttpServer] closed.", s.addr)
	}
}

func (s *HttpServer) SetSSLFile(certFile, keyFile string) {
	s.certFile = certFile
	s.keyFile = keyFile
}

// RespJson 返回json数据
func (s *HttpServer) RespJson(w http.ResponseWriter, data interface{}) iface.HttpServerResType {
	// 设置响应头为JSON格式
	w.Header().Set("Content-Type", "application/json")

	errCode, msg := s.httpCodeMsgHandler(INNER_HTTP_SUCCESS)

	err := json.NewEncoder(w).Encode(WebRespData{
		ErrCode: errCode,
		Message: msg,
		Data:    data,
	})
	if err != nil {
		elog.Error("[HttpServer] resp json err:", s.addr, err)
		return false
	}

	return true
}

// RespErr 返回错误码
func (s *HttpServer) RespErr(w http.ResponseWriter, errCode uint32) iface.HttpServerResType {
	// 设置响应头为JSON格式
	w.Header().Set("Content-Type", "application/json")

	errCode, msg := s.httpCodeMsgHandler(errCode)

	err := json.NewEncoder(w).Encode(WebRespData{
		ErrCode: errCode,
		Message: msg,
	})
	if err != nil {
		elog.Error("[HttpServer] resp err code err:", s.addr, err)
		return false
	}

	return true
}

// RespCustom 自定义返回
func (s *HttpServer) RespCustom(err error) iface.HttpServerResType {
	if err != nil {
		return false
	}

	return true
}

func (s *HttpServer) runAsHttps() {
	// 打印服务器启动信息
	elog.Info("[HttpServer] Start https...", s.addr)
	elog.Info("[HttpServer] Listen https://", s.addr)

	elog.Info("[HttpServer] CertFile:", s.certFile)
	elog.Info("[HttpServer] KeyFile:", s.keyFile)

	// 启动HTTPS服务器
	err := s.Server.ListenAndServeTLS(s.certFile, s.keyFile)
	if err != nil {
		elog.Fatal("Start https server failed:", err)
	}
}

func (s *HttpServer) runAsHttp() {
	// 打印服务器启动信息
	elog.Info("Start http...", s.addr)
	elog.Info("Listen http://", s.addr)

	// 启动HTTP服务器
	err := s.Server.ListenAndServe()
	if err != nil {
		elog.Fatal("Start http server failed:", err)
	}
}

// HandlePostFunc 注册Post事件
func HandlePostFunc[T interface{}](s *HttpServer, pattern string, handler func(http.ResponseWriter, *http.Request, T) iface.HttpServerResType) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 不是POST方法
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			s.RespErr(w, INNER_HTTP_POST_ONLY)
			return
		}

		contentType := r.Header.Get("Content-Type")

		// 请求结构不是json
		if !strings.Contains(contentType, "application/json") {
			s.RespErr(w, INNER_HTTP_NEED_JSON_TYPE)
			return
		}

		var reqData T

		// 解析JSON格式的请求体
		err := json.NewDecoder(r.Body).Decode(&reqData)

		if err != nil {
			elog.Error("[HttpServer] decode req json data err:", s.addr, err)
			s.RespErr(w, INNER_HTTP_WRONG_REQ_DATA)
			return
		}

		// 回调
		handler(w, r, reqData)
	})
}

// HandleGetFunc 注册Get事件
func HandleGetFunc(s *HttpServer, pattern string, handler func(http.ResponseWriter, *http.Request) iface.HttpServerResType) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 不是Get方法
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			s.RespErr(w, INNER_HTTP_GET_ONLY)
			return
		}

		// 回调
		handler(w, r)
	})
}

// HandleCustomFunc 注册自定义事件
func HandleCustomFunc(s *HttpServer, pattern string, handler func(http.ResponseWriter, *http.Request) iface.HttpServerResType) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 回调
		handler(w, r)
	})
}
