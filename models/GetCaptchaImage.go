package m

import (
	"net/http"
)

type GetCaptchaImageRequest struct {
	TaskId string
}

type GetCaptchaImageResponse struct {
	TaskId    string
	ImageData []byte
}

func (req *GetCaptchaImageRequest) Check() (err error) {
	if req == nil {
		return NewErrorWithHttpStatusCode(Err_RequestIsAbsent, http.StatusBadRequest)
	}

	if len(req.TaskId) == 0 {
		return NewErrorWithHttpStatusCode(Err_IdIsNotSet, http.StatusBadRequest)
	}

	return nil
}

func NewGetImageRequestFromHttpRequest(httpReq *http.Request) (req *GetCaptchaImageRequest, err error) {
	req = &GetCaptchaImageRequest{
		TaskId: httpReq.URL.Query().Get(QueryKeyId),
	}

	return req, nil
}
