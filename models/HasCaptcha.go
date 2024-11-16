package m

import (
	"net/http"
)

type HasCaptchaRequest struct {
	TaskId string
}

type HasCaptchaResponse struct {
	TaskId  string
	IsFound bool
}

func (req *HasCaptchaRequest) Check() (err error) {
	if req == nil {
		return NewErrorWithHttpStatusCode(Err_RequestIsAbsent, http.StatusBadRequest)
	}

	if len(req.TaskId) == 0 {
		return NewErrorWithHttpStatusCode(Err_IdIsNotSet, http.StatusBadRequest)
	}

	return nil
}
