package controllers

import (
	"strconv"
)

type Resp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ToInt(s string) (int64, error) {
	res, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		res = 0
	}
	return res, err
}

func BuildResp(code, msg string, data interface{}) *Resp {
	return &Resp{code, msg, data}
}
