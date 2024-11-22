package main

import (
	"fmt"
)

func check(s string) (string, error) {
	e := &MyError{}
	if s == "" {
		return "", e.Fail(400, "s is empty")
	} else {
		return s, e.Success()
	}
}

type MyError struct {
	Code int
	Msg  string
}

func (e MyError) Error() string {
	return e.Msg
}
func (e MyError) Fail(code int, msg string) MyError {
	e.Code = code
	e.Msg = msg
	return e
}
func (e MyError) Success() MyError {
	e.Code = 200
	e.Msg = "success"
	return e
}
func main() {
	s, err := check("")
	if err != nil {
		fmt.Printf("err: %v\n,code: %v ", err.Error(), err.(MyError).Code)
	} else {
		fmt.Printf("s: %v\n", s)
	}
}
