package common

type Result struct {
	Code   int         `json:"code"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func Error(code int, errMsg string) *Result {
	return &Result{
		Code:   code,
		ErrMsg: errMsg,
		Data:   nil,
	}
}

func Success(code int, data interface{}) *Result {
	return &Result{
		Code:   code,
		ErrMsg: "",
		Data:   data,
	}
}
