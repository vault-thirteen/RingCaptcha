package models

type CheckCaptchaRequest struct {
	TaskId string
	Value  uint
}

type CheckCaptchaResponse struct {
	TaskId    string
	IsFound   bool
	IsSuccess bool
}
