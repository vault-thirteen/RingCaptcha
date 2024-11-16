package m

import (
	"image"
)

type Captcha struct {
	Id          string
	Image       *image.NRGBA
	ImageData   []byte
	ImageFormat string
	RingCount   uint
}

func NewCaptchaWithId(id string) (captcha *Captcha) {
	return &Captcha{Id: id}
}

func NewCaptchaWithIdAndAnswer(id string, answer uint) (captcha *Captcha) {
	return &Captcha{Id: id, RingCount: answer}
}
