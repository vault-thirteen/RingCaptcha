package m

type IEasyServer interface {
	Start()
	Stop()
	GetListenDsn() (dsn string)

	CreateCaptcha() (captcha *Captcha)
	CheckCaptcha(captcha *Captcha) (ok bool)
	HasCaptcha(captcha *Captcha) (ok bool)
	GetCaptchaImage(captcha *Captcha) (data []byte)
}
