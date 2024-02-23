package handlers

// Result structure for easily json serialization
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// UpdateAllFields in Result and return it
func (res *Result) UpdateAllFields(code int, msg string, data interface{}) *Result {
	res.Code = code
	res.Message = msg
	res.Data = data
	return res
}

// UpdateDataField in Result and return it
func (res *Result) UpdateDataField(data interface{}) *Result {
	res.Data = data
	return res
}
