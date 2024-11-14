package models

type CreateCaptchaRequest struct{}

type CreateCaptchaResponse struct {
	TaskId              string
	ImageFormat         string
	IsImageDataReturned bool
	ImageData           []byte
}
