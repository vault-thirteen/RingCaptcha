package models

type CheckCaptchaResponse struct {
	TaskId    string
	IsFound   bool
	IsSuccess bool
}
