package middleware

import (
	"Turing-Go/constant"
	"Turing-Go/net"
)

func CheckRole() net.MiddlewareFunc {
	return func(next net.Handler) net.Handler {
		return func(req *net.WsMsgReq, resp *net.WsMsgResp) {
			_, err := req.Conn.GetProperty("role")
			if err != nil {
				resp.Body.Code = constant.SessionInvalid
				return
			}
			// do something
			next(req, resp)
		}
	}
}
