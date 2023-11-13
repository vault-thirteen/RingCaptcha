package models

type CommonResult struct {
	// Time taken to perform the request, in milliseconds.
	TimeSpent int64 `json:"timeSpent"`
}

type PingParams struct{}

type PingResult struct {
	OK bool `json:"ok"`
}

type CreateCaptchaParams struct{}

type CreateCaptchaResult struct {
	CommonResult

	TaskId              string `json:"taskId"`
	ImageFormat         string `json:"imageFormat"`
	IsImageDataReturned bool   `json:"isImageDataReturned"`
	ImageDataB64        string `json:"imageDataB64,omitempty"`
}

type CheckCaptchaParams struct {
	TaskId string `json:"taskId"`
	Value  uint   `json:"value"`
}

type CheckCaptchaResult struct {
	CommonResult

	TaskId    string `json:"taskId"`
	IsSuccess bool   `json:"isSuccess"`
}

type ShowDiagnosticDataParams struct{}

type ShowDiagnosticDataResult struct {
	CommonResult

	TotalRequestsCount      uint64 `json:"totalRequestsCount"`
	SuccessfulRequestsCount uint64 `json:"successfulRequestsCount"`
}
