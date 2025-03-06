package constant

const (
	CODE_SUCCESS = "0000"
	CODE_ERROR   = "9999"

	STATUS_SUCCESS       = "success"
	STATUS_GENERIC_ERROR = "Internal server error"
	STATUS_BAD_INPUT     = "Bad request"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespSuccess() *Response {
	return &Response{
		Code:    CODE_SUCCESS,
		Message: STATUS_SUCCESS,
		Data:    nil,
	}
}

func RespSuccessWithData(data interface{}) *Response {
	return &Response{
		Code:    CODE_SUCCESS,
		Message: STATUS_SUCCESS,
		Data:    data,
	}
}
