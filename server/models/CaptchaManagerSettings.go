package models

import "errors"

const (
	ErrImagesFolderIsNotSet = "images folder is not set"
	ErrImageWidthIsNotSet   = "image width is not set"
	ErrImageHeightIsNotSet  = "image height is not set"
	ErrImageTtlIsNotSet     = "image TTL is not set"
	ErrCacheSettingsError   = "error in cache settings"
)

type CaptchaManagerSettings struct {
	// Main settings.
	IsImageStorageUsed        bool
	IsImageServerEnabled      bool
	IsImageCleanupAtStartUsed bool
	IsCachingEnabled          bool

	// Image settings.
	ImagesFolder string
	ImageWidth   uint
	ImageHeight  uint
	ImageTtlSec  uint

	// HTTP server settings.
	HttpHost       string
	HttpPort       uint16
	HttpErrorsChan *chan error
	HttpServerName string

	// Cache settings.
	CacheSizeLimit   int
	CacheVolumeLimit int
	CacheRecordTtl   uint
}

func (s *CaptchaManagerSettings) Check() (err error) {
	if s.IsImageStorageUsed {
		if len(s.ImagesFolder) == 0 {
			return errors.New(ErrImagesFolderIsNotSet)
		}
	}

	if s.ImageWidth == 0 {
		return errors.New(ErrImageWidthIsNotSet)
	}
	if s.ImageHeight == 0 {
		return errors.New(ErrImageHeightIsNotSet)
	}
	if s.ImageTtlSec == 0 {
		return errors.New(ErrImageTtlIsNotSet)
	}

	if s.IsCachingEnabled {
		if (s.CacheSizeLimit <= 0) ||
			(s.CacheVolumeLimit <= 0) ||
			(s.CacheRecordTtl == 0) ||
			(s.CacheRecordTtl > s.ImageTtlSec) {
			return errors.New(ErrCacheSettingsError)
		}
	}

	return nil
}
