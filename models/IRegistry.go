package m

type IRegistry interface {
	Start()
	Stop()

	CreateCaptcha(captcha *Captcha) (err error)
	CheckCaptcha(captcha *Captcha) (ok bool, isFound bool, err error)
	HasCaptcha(captcha *Captcha) (exists bool, err error)
	GetCaptchaImage(captcha *Captcha) (data []byte, err error)
}
