package m

import (
	"errors"
)

type CaptchaServerSettings struct {
	// Main settings.
	IsImageStorageUsed        bool
	IsImageServerEnabled      bool
	IsImageCleanupAtStartUsed bool
	IsStorageCleaningEnabled  bool

	// Image settings.
	ImagesFolder      string
	ImageWidth        uint
	ImageHeight       uint
	FilesCountToClean int

	// HTTP server settings.
	HttpHost       string
	HttpPort       uint16
	HttpErrorsChan *chan error
	HttpServerName string

	// File cache settings.
	FileCacheSizeLimit   int
	FileCacheVolumeLimit int
	FileCacheItemTtl     uint

	// Record cache settings.
	RecordCacheSizeLimit int
	RecordCacheItemTtl   uint
}

func (s *CaptchaServerSettings) Check() (err error) {
	if s.IsImageStorageUsed {
		if len(s.ImagesFolder) == 0 {
			return errors.New(Err_ImagesFolderIsNotSet)
		}
	}

	if s.ImageWidth == 0 {
		return errors.New(Err_ImageWidthIsNotSet)
	}
	if s.ImageHeight == 0 {
		return errors.New(Err_ImageHeightIsNotSet)
	}

	if (s.FilesCountToClean <= 0) ||
		(len(s.HttpHost) == 0) ||
		(s.HttpPort == 0) ||
		(s.HttpErrorsChan == nil) ||
		(len(s.HttpServerName) == 0) ||
		(s.FileCacheSizeLimit <= 0) ||
		(s.FileCacheVolumeLimit <= 0) ||
		(s.FileCacheItemTtl <= 0) ||
		(s.RecordCacheSizeLimit <= 0) ||
		(s.RecordCacheItemTtl <= 0) {
		return errors.New(Err_CacheSettingsError)
	}

	if (s.IsStorageCleaningEnabled) && (!s.IsImageStorageUsed) {
		return errors.New(Err_CacheSettingsError)
	}

	if s.FileCacheItemTtl != s.RecordCacheItemTtl {
		return errors.New(Err_CacheSettingsError)
	}

	return nil
}
