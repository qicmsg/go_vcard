package result

const (
	OK = 200
)

var resultMsg = map[int]string{
	OK: "成功",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func ResultMsg(code int) string {
	return resultMsg[code]
}

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
