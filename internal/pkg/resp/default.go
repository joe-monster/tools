package resp

const SUCCESS_CODE = 0
const GENERAL_FAIL_CODE = 1

type DataType struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

func Fail(err error) DataType {
	var rtn DataType
	rtn.Code = GENERAL_FAIL_CODE
	rtn.Msg = err.Error()
	return rtn
}

func Success(data interface{}) DataType {
	var rtn DataType
	rtn.Code = SUCCESS_CODE
	rtn.Msg = ""
	rtn.Data = data
	return rtn
}
