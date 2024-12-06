package middleware

import (
	"Turing-Go/net"
	"fmt"
	"log"
)

func Log() net.MiddlewareFunc {
	return func(next net.Handler) net.Handler {
		return func(req *net.WsMsgReq, resp *net.WsMsgResp) {
			log.Println("请求路由", req.Body.Name)
			log.Println("请求数据", fmt.Sprintf("%v", req.Body.Msg))
			// do something
			next(req, resp)
		}
	}
}
