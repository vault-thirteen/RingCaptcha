package models

type HasCaptchaRequest struct {
	TaskId string
}

type HasCaptchaResponse struct {
	TaskId  string
	IsFound bool
}
