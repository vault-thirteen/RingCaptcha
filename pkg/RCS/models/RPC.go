package models

type CreateCaptchaParams struct {
}

type CreateCaptchaResult struct {
	TaskId              string `json:"taskId"`
	ImageFormat         string `json:"imageFormat"`
	IsImageDataReturned bool   `json:"isImageDataReturned"`
	ImageDataB64        string `json:"imageDataB64,omitempty"`
	TimeSpent           int64  `json:"timeSpent"` // Time taken to make the action, in milliseconds.
}

type CheckCaptchaParams struct {
	TaskId string `json:"taskId"`
	Value  uint   `json:"value"`
}

type CheckCaptchaResult struct {
	TaskId    string `json:"taskId"`
	IsSuccess bool   `json:"isSuccess"`
	TimeSpent int64  `json:"timeSpent"` // Time taken to make the action, in milliseconds.
}

type PingParams struct{}

type PingResult struct {
	OK bool `json:"ok"`
}
