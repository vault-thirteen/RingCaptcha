package main

import (
	"time"

	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/random"

	"github.com/vault-thirteen/RingCaptcha/easy_server"
	"github.com/vault-thirteen/RingCaptcha/models"
)

func main() {
	showUsageExample()
}

func showUsageExample() {
	var easyServer = es.NewEasyServer(getSettings())
	easyServer.Start()
	defer easyServer.Stop()

	var c *m.Captcha
	var ok bool
	var guess int
	var imageFile []byte

	// Go language is quite stupid. It writes message into stdout and stderr
	// in random order, so here we wait until it shows all the messages.
	// Since this code is just a simple example, here we do not care about
	// performance.
	time.Sleep(time.Second)

	var useIntroduction = true

	if useIntroduction {
		var n = 5
		for i := 0; i < n; i++ {
			c = easyServer.CreateCaptcha()
			printResultOfCreateCaptcha(c)

			ok = easyServer.HasCaptcha(c)
			printResultOfHasCaptcha(ok)

			imageFile = easyServer.GetCaptchaImage(c)
			printResultOfGetCaptchaImage(c, imageFile)
			time.Sleep(10 * time.Second)

			guess = getRandomAnswer()
			c.RingCount = uint(guess)
			ok = easyServer.CheckCaptcha(c)
			printResultOfCheckCaptcha(uint(guess), ok)
			time.Sleep(time.Second)

			cio.PrintMessage()

			if ok {
				break
			}
		}
	}

	cio.PrintMessage("Now it is your turn to guess.")
	cio.PrintMessage()
	time.Sleep(5 * time.Second)

	for {
		c = easyServer.CreateCaptcha()
		printResultOfCreateCaptcha(c)

		ok = easyServer.HasCaptcha(c)
		printResultOfHasCaptcha(ok)

		imageFile = easyServer.GetCaptchaImage(c)
		printResultOfGetCaptchaImage(c, imageFile)
		time.Sleep(time.Second)

		guess = getNumberFromUserInput()
		c.RingCount = uint(guess)
		ok = easyServer.CheckCaptcha(c)
		printResultOfCheckCaptcha(uint(guess), ok)
		time.Sleep(time.Second)

		cio.PrintMessage()

		if ok {
			break
		}
	}

	cio.PrintMessage("Have a good time.")
}

func getSettings() (settings *m.CaptchaServerSettings) {
	httpErrorsChan := make(chan error)

	settings = &m.CaptchaServerSettings{
		// Main settings.
		IsImageStorageUsed:        true,
		IsImageServerEnabled:      true,
		IsImageCleanupAtStartUsed: true,
		IsStorageCleaningEnabled:  true,

		// Image settings.
		ImagesFolder: "test\\4",
		ImageWidth:   256,
		ImageHeight:  256,

		// This number is only an example.
		// To save your SSD or HDD, this count must be at least 1'000 to
		// reduce the wear of hardware. With current settings, one image uses
		// approximately 25 KB, so 1'000 images would take 25 MB of space.
		FilesCountToClean: 5,

		// HTTP server settings.
		HttpHost:       "localhost",
		HttpPort:       2000,
		HttpErrorsChan: &httpErrorsChan,
		HttpServerName: "RCS",

		// File cache settings.
		FileCacheSizeLimit:   50,
		FileCacheVolumeLimit: 1_000_000,
		FileCacheItemTtl:     60,

		// Record cache settings.
		RecordCacheSizeLimit: 50,
		RecordCacheItemTtl:   60,
	}

	return settings
}
func getRandomAnswer() int {
	x, err := random.Uint(1, 10)
	mustBeNoError(err)
	return int(x)
}
func printResultOfCreateCaptcha(c *m.Captcha) {
	cio.PrintMessage(cio.FormatMessage("A new captcha was created. ID=%v.", c.Id))
}
func printResultOfHasCaptcha(ok bool) {
	msg := "The captcha should exist, "

	if !ok {
		msg = msg + "but someone has stolen it."
	} else {
		msg = msg + "and it does still exist."
	}

	cio.PrintMessage(msg)
}
func printResultOfGetCaptchaImage(c *m.Captcha, imageFile []byte) {
	if len(imageFile) < 1000 {
		msg := "something bad has happened"
		panic(msg)
	}

	url := `http://localhost:2000?id=` + c.Id // This address depends on the settings.

	msg := "Unfortunately, Go language still does not support graphical interface. \r\n"
	msg = msg + cio.FormatMessage("Open your browser to see the image: %v", url)

	cio.PrintMessage(msg)
}
func getNumberFromUserInput() (guess int) {
	cio.Print("Enter your number: ")
	return cio.GetInputIntRetry()
}

func printResultOfCheckCaptcha(guess uint, ok bool) {
	msg := cio.FormatMessage("We tried to guess the answer with number %v and ", guess)

	if !ok {
		msg = msg + "we failed."
	} else {
		msg = msg + "we succeeded."
	}

	cio.PrintMessage(msg)
}
func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
