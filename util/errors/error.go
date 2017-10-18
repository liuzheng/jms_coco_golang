package errors

type Error interface {
	Error() string
	GetCode() int
	SetCode(int)
	SetError(string)
}

//API错误
type RespError struct {
	ErrorMsg  string `json:"error"`
	ErrorCode int
}

// New returns an error that formats as the given text.
func New(text string, code int) Error {
	return &RespError{ErrorMsg:text, ErrorCode:code}
}

func (e *RespError) Error() string {
	return e.ErrorMsg
}

func (e *RespError) GetCode() int {
	return e.ErrorCode
}

func (e *RespError)SetCode(code int) {
	e.ErrorCode = code
}

func (e *RespError)SetError(err string) {
	e.ErrorMsg = err
}