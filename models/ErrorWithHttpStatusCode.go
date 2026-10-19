package m

type ErrorWithHttpStatusCode struct {
	msg            string
	httpStatusCode int
}

func NewErrorWithHttpStatusCode(msg string, httpStatusCode int) *ErrorWithHttpStatusCode {
	return &ErrorWithHttpStatusCode{msg: msg, httpStatusCode: httpStatusCode}
}

func NewErrorWithHttpStatusCodeFromError(e error) *ErrorWithHttpStatusCode {
	if e == nil {
		return nil
	}

	return e.(*ErrorWithHttpStatusCode)
}

func (e *ErrorWithHttpStatusCode) Error() string {
	return e.msg
}

func (e *ErrorWithHttpStatusCode) GetHttpStatusCode() int {
	return e.httpStatusCode
}
