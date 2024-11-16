package m

import (
	"net/http"
)

type CheckCaptchaRequest struct {
	TaskId string
	Value  uint
}

type CheckCaptchaResponse struct {
	TaskId    string
	IsSuccess bool
}

func (req *CheckCaptchaRequest) Check() (err error) {
	if req == nil {
		return NewErrorWithHttpStatusCode(Err_RequestIsAbsent, http.StatusBadRequest)
	}

	if len(req.TaskId) == 0 {
		return NewErrorWithHttpStatusCode(Err_IdIsNotSet, http.StatusBadRequest)
	}

	return nil
}
