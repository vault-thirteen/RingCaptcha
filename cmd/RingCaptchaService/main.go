package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	rcs "github.com/vault-thirteen/RingCaptcha/pkg/RCS"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/settings"
	"github.com/vault-thirteen/Versioneer"
)

func main() {
	showIntro()

	cla, err := rcs.ReadCLA()
	mustBeNoError(err)
	if cla.IsDefaultFile() {
		log.Println("Using the default configuration file.")
	}

	var stn *settings.Settings
	stn, err = settings.NewSettingsFromFile(cla.ConfigurationFilePath)
	mustBeNoError(err)

	log.Println("Server is starting ...")
	var srv *rcs.Server
	srv, err = rcs.NewServer(stn)
	mustBeNoError(err)

	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("HTTP Server: " + srv.GetListenDsn())

	serverMustBeStopped := srv.GetStopChannel()
	waitForQuitSignalFromOS(serverMustBeStopped)
	<-*serverMustBeStopped

	log.Println("Stopping the server ...")
	err = srv.Stop()
	if err != nil {
		log.Println(err)
	}
	log.Println("Server was stopped.")
	time.Sleep(time.Second)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func showIntro() {
	versioneer, err := ver.New()
	mustBeNoError(err)
	versioneer.ShowIntroText("Mailer Module")
	versioneer.ShowComponentsInfoText()
	fmt.Println()
}

func waitForQuitSignalFromOS(serverMustBeStopped *chan bool) {
	osSignals := make(chan os.Signal, 16)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for sig := range osSignals {
			switch sig {
			case syscall.SIGINT,
				syscall.SIGTERM:
				log.Println("quit signal from OS has been received: ", sig)
				*serverMustBeStopped <- true
			}
		}
	}()
}
