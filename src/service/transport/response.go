package transport

type APIResponse struct {
	Data  interface{}    `json:"data""`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Msg  string       `json:"msg"`
	Code ResponseCode `json:"code"`
}

type ResponseCode int

const (
	ErrorCodeInvalidParameter ResponseCode = 1
	ErrorCodePermissionDenied ResponseCode = 2
	ErrorCodeInternal         ResponseCode = 3
	ErrorCodeNotFound         ResponseCode = 4
	ErrorCodeEmpty            ResponseCode = 5
	ErrorCodeNotImplemented   ResponseCode = 6
	ErrorCodeUnauthorized     ResponseCode = 7
)

type Error struct {
	error
	Msg  string
	Code ResponseCode
}

func (e Error) Error() string {
	return e.Msg
}
