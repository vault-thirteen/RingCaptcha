package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/RingCaptcha/easy_server"
	"github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/auxie/random"
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
	var guess uint
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
			c.RingCount = guess
			ok = easyServer.CheckCaptcha(c)
			printResultOfCheckCaptcha(guess, ok)
			time.Sleep(time.Second)

			fmt.Println()

			if ok {
				break
			}
		}
	}

	fmt.Println("Now it is your turn to guess.")
	fmt.Println()
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
		c.RingCount = guess
		ok = easyServer.CheckCaptcha(c)
		printResultOfCheckCaptcha(guess, ok)
		time.Sleep(time.Second)

		fmt.Println()

		if ok {
			break
		}
	}

	fmt.Println("Have a good time.")
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
func getRandomAnswer() (x uint) {
	var err error
	x, err = random.Uint(1, 10)
	mustBeNoError(err)
	return x
}
func printResultOfCreateCaptcha(c *m.Captcha) {
	fmt.Println(fmt.Sprintf("A new captcha was created. ID=%v.", c.Id))
}
func printResultOfHasCaptcha(ok bool) {
	msg := "The captcha should exist,"

	if !ok {
		msg = msg + fmt.Sprintf(" but some has stolen it.")
	} else {
		msg = msg + fmt.Sprintf(" and it does still exist.")
	}

	fmt.Println(msg)
}
func printResultOfGetCaptchaImage(c *m.Captcha, imageFile []byte) {
	if len(imageFile) < 1000 {
		msg := "something bad has happened"
		panic(msg)
	}

	url := `http://localhost:2000?id=` + c.Id // This address depends on the settings.

	msg := "Unfortunately, Go language still does not support graphical interface. \r\n"
	msg = msg + fmt.Sprintf("Open your browser to see the image: %v", url)

	fmt.Println(msg)
}
func getNumberFromUserInput() (guess uint) {
	fmt.Print("Enter your number: ")

	var err error
	for {
		_, err = fmt.Scanf("%d", &guess)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		break
	}

	return guess
}

func printResultOfCheckCaptcha(guess uint, ok bool) {
	msg := fmt.Sprintf("We tried to guess the answer with number %v and", guess)

	if !ok {
		msg = msg + fmt.Sprintf(" we failed.")
	} else {
		msg = msg + fmt.Sprintf(" we succeeded.")
	}

	fmt.Println(msg)
}
func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
