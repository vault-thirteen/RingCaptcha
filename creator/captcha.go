package c

import (
	"bytes"
	"errors"
	"image/png"

	"github.com/google/uuid"
	"github.com/vault-thirteen/RingCaptcha/models"
)

func CreateCaptcha(w, h uint) (c *m.Captcha, err error) {
	c = &m.Captcha{
		Id:          createRandomUID(),
		ImageFormat: m.ImageFormat,
	}

	c.Image, c.RingCount, err = CreateCaptchaImage(w, h, true, false)
	if err != nil {
		return nil, err
	}

	err = prepareImageBinaryData(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func createRandomUID() (uid string) {
	return m.RUIDPrefix + uuid.New().String()
}

func prepareImageBinaryData(c *m.Captcha) (err error) {
	if c.ImageFormat != m.ImageFormat {
		return errors.New(m.Err_ImageFormatIsInvalid)
	}

	buf := new(bytes.Buffer)

	err = png.Encode(buf, c.Image)
	if err != nil {
		return err
	}

	c.ImageData = buf.Bytes()

	return nil
}
