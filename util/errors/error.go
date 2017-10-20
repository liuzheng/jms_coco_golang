package errors

type Error interface {
	Error() string
	GetCode() int
	SetCode(int)
	SetError(string)
}

//API错误
type UtilError struct {
	ErrorMsg  string `json:"error"`
	ErrorCode int `json:"code"`
}

// New returns an error that formats as the given text.
func New(text string, code int) Error {
	return &UtilError{ErrorMsg:text, ErrorCode:code}
}

func (e *UtilError) Error() string {
	return e.ErrorMsg
}

func (e *UtilError) GetCode() int {
	return e.ErrorCode
}

func (e *UtilError)SetCode(code int) {
	e.ErrorCode = code
}

func (e *UtilError)SetError(err string) {
	e.ErrorMsg = err
}