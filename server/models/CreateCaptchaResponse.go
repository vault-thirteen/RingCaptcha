package models

type CreateCaptchaResponse struct {
	TaskId              string
	ImageFormat         string
	IsImageDataReturned bool
	ImageData           []byte
}
