package main

import (
	"fmt"
	"time"

	"github.com/vault-thirteen/RingCaptcha/server"
	"github.com/vault-thirteen/RingCaptcha/server/models"
)

func main() {
	httpErrorsChan := make(chan error)

	var s = &models.CaptchaManagerSettings{
		// Main settings.
		IsImageStorageUsed:        true,
		IsImageServerEnabled:      true,
		IsImageCleanupAtStartUsed: true,
		IsCachingEnabled:          true,

		// Image settings.
		ImagesFolder: "images",
		ImageWidth:   256,
		ImageHeight:  256,
		ImageTtlSec:  60,

		// HTTP server settings.
		HttpHost:       "localhost",
		HttpPort:       2000,
		HttpErrorsChan: &httpErrorsChan,
		HttpServerName: "RCS",

		// Cache settings.
		CacheSizeLimit:   64,
		CacheVolumeLimit: 16_000_000,
		CacheRecordTtl:   60,
	}

	cm, err := server.NewCaptchaManager(s)
	mustBeNoError(err)
	err = cm.Start()
	mustBeNoError(err)

	var resp *models.CreateCaptchaResponse
	resp, err = cm.CreateCaptcha()
	mustBeNoError(err)
	fmt.Println("TaskId:", resp.TaskId)

	time.Sleep(1 * time.Minute)

	err = cm.Stop()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}
