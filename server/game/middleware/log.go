package middleware

import (
	"Turing-Go/net"
	"log"
)

func Log() net.MiddlewareFunc {
	return func(next net.Handler) net.Handler {
		return func(req *net.WsMsgReq, resp *net.WsMsgResp) {
			log.Println("log 中间件录音：", req.Body.Name)
			log.Printf("log 中间件请求参数： %v \n", req.Body.Msg)
			// do something
			next(req, resp)
		}
	}
}
