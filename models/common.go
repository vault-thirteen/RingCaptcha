package m

const (
	ImageFileExt            = "png"
	FileExtFullPng          = "." + ImageFileExt
	ImageFormat             = "PNG"
	ImageMimeType           = "image/png"
	QueryKeyId              = "id"
	RUIDPrefix              = "RCS-"
	CaptchaImageMinWidth    = 128
	CaptchaImageMinHeight   = 128
	RingMinCount            = 3
	RingMaxCount            = 6
	RingMinRadius           = 24
	BrushOuterRadiusMin     = 2
	BrushOuterRadiusMax     = 32
	ColourComponentMaxValue = 65535
	KR                      = 0.5
	KDMin                   = 1.0
	KDMax                   = 1.5

	// C1 is a maximum colour channel value of a pixel in an 'image.RGBA' of
	// the built-in library.
	C1 = float64(65_535)
)

const (
	Err_Anomaly                       = "anomaly"
	Err_BrushRadiusRatio              = "brush radius ratio error"
	Err_CacheSettingsError            = "error in cache settings"
	Err_CanvasIsTooSmall              = "canvas is too small"
	Err_DensityCoefficient            = "density coefficient error"
	Err_Dimensions                    = "dimensions error"
	ErrF_FileExtensionMismatch        = "file extension mismatch: png vs %s"
	Err_FileStorageIsDisabled         = "file storage is disabled"
	Err_IdIsDuplicate                 = "duplicate ID"
	Err_IdIsNotFound                  = "ID is not found"
	Err_IdIsNotSet                    = "ID is not set"
	Err_ImageFormatIsInvalid          = "image format is invalid"
	Err_ImagesFolderIsNotSet          = "images folder is not set"
	Err_ImagesHaveDifferentDimensions = "images have different dimensions"
	Err_ImageHeightIsNotSet           = "image height is not set"
	Err_ImageWidthIsNotSet            = "image width is not set"
	Err_RequestIsAbsent               = "request is absent"
)

const (
	Msg_CaptchaManagerStart    = "Captcha manager has started"
	Msg_CaptchaManagerStop     = "Captcha manager has been stopped"
	Msg_CleaningImagesFolder   = "Cleaning the images folder ... "
	Msg_Done                   = "Done"
	Msg_Failure                = "Failure"
	Msg_ImageCleanerHasStarted = "image cleaner has started"
	Msg_ImageCleanerHasStopped = "image cleaner has stopped"
	MsgF_CleaningImages        = "cleaning %v images ... "
)
