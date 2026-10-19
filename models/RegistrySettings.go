package m

type RegistrySettings struct {
	// Main settings.
	IsImageStorageUsed       bool
	IsStorageCleaningEnabled bool

	// Image settings.
	ImagesFolder      string
	FilesCountToClean int

	// File cache settings.
	FileCacheSizeLimit   int
	FileCacheVolumeLimit int
	FileCacheItemTtl     uint

	// Record cache settings.
	RecordCacheSizeLimit int
	RecordCacheItemTtl   uint
}
