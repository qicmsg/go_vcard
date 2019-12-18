package utils

import "vcard/pkg/utils/result"

func Result(code int, data interface{}, msg string) result.Result {
	if msg == "" {
		msg = result.ResultMsg(code)
	}

	return result.Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
