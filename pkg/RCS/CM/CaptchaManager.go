package capman

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/models"
	"github.com/vault-thirteen/RingCaptcha/pkg/captcha"
	cos "github.com/vault-thirteen/RingCaptcha/pkg/os"
	hdr "github.com/vault-thirteen/header"
)

const (
	ErrImagesFolderIsNotSet = "images folder is not set"
	ErrImageWidthIsNotSet   = "image width is not set"
	ErrImageHeightIsNotSet  = "image height is not set"
	ErrImageTTLIsNotSet     = "image TTL is not set"
)

const (
	ImageFileExt  = "png"
	ImageFormat   = "PNG"
	ImageMimeType = "image/png"
)

// CaptchaManager is a captcha manager.
type CaptchaManager struct {
	// Images.
	storeImages              bool
	imagesFolder             string
	imageWidth               uint
	imageHeight              uint
	imageTTLSec              uint
	registry                 *Registry
	clearImagesFolderAtStart bool

	// HTTP server.
	useHttpServerForImages bool
	httpServer             *http.Server
	listenDsn              string
	httpErrorsChan         *chan error
	httpServerName         string
}

func NewCaptchaManager(
	storeImages bool,
	imagesFolder string,
	imageWidth uint,
	imageHeight uint,
	imageTTLSec uint,
	clearImagesFolderAtStart bool,
	useHttpServerForImages bool,
	httpHost string,
	httpPort uint16,
	httpErrorsChan *chan error,
	httpServerName string,
) (cm *CaptchaManager, err error) {
	if storeImages {
		if len(imagesFolder) == 0 {
			return nil, errors.New(ErrImagesFolderIsNotSet)
		}
	}

	if imageWidth == 0 {
		return nil, errors.New(ErrImageWidthIsNotSet)
	}
	if imageHeight == 0 {
		return nil, errors.New(ErrImageHeightIsNotSet)
	}
	if imageTTLSec == 0 {
		return nil, errors.New(ErrImageTTLIsNotSet)
	}

	cm = &CaptchaManager{
		storeImages:              storeImages,
		imagesFolder:             imagesFolder,
		imageWidth:               imageWidth,
		imageHeight:              imageHeight,
		imageTTLSec:              imageTTLSec,
		registry:                 NewRegistry(storeImages, imagesFolder, imageTTLSec),
		clearImagesFolderAtStart: clearImagesFolderAtStart,
		useHttpServerForImages:   useHttpServerForImages,
	}

	if cm.clearImagesFolderAtStart {
		fmt.Print("Cleaning the images folder ... ")
		err = cm.clearImagesFolder()
		if err != nil {
			return nil, err
		}
		fmt.Println("Done")
	}

	if cm.useHttpServerForImages {
		err = cm.initHttpServer(httpHost, httpPort, httpErrorsChan, httpServerName)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Captcha manager has started.")

	return cm, nil
}

func (cm *CaptchaManager) initHttpServer(
	httpHost string,
	httpPort uint16,
	httpErrorsChan *chan error,
	httpServerName string,
) (err error) {
	cm.listenDsn = fmt.Sprintf("%s:%d", httpHost, httpPort)
	cm.httpServer = &http.Server{
		Addr:    cm.listenDsn,
		Handler: http.Handler(http.HandlerFunc(cm.httpRouter)),
	}
	cm.httpErrorsChan = httpErrorsChan
	cm.httpServerName = httpServerName

	go func() {
		var listenError error
		listenError = cm.httpServer.ListenAndServe()

		if (listenError != nil) && (listenError != http.ErrServerClosed) {
			*cm.httpErrorsChan <- listenError
		}
	}()

	return nil
}

func (cm *CaptchaManager) GetListenDsn() (dsn string) {
	return cm.listenDsn
}

func (cm *CaptchaManager) httpRouter(rw http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	if len(id) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if !cm.isIdRegistered(id) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if !cm.storeImages {
		// No images are stored. Sorry.
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	var fileContents []byte
	var err error
	fileContents, err = cm.readFile(id)
	if err != nil {
		log.Println(err) //TODO
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set(hdr.HttpHeaderContentType, ImageMimeType)
	rw.Header().Set(hdr.HttpHeaderServer, cm.httpServerName)
	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(fileContents)
	if err != nil {
		log.Println(err)
	}
}

func (cm *CaptchaManager) CreateCaptcha() (resp *models.CreateCaptchaResponse, err error) {
	var img *image.NRGBA
	var ringCount uint
	img, ringCount, err = captcha.CreateCaptchaImage(cm.imageWidth, cm.imageHeight, true, false)
	if err != nil {
		return nil, err
	}

	resp = &models.CreateCaptchaResponse{
		TaskId:              cm.createRandomUID(),
		ImageFormat:         ImageFormat,
		IsImageDataReturned: !cm.storeImages,
	}

	if resp.IsImageDataReturned {
		resp.ImageData, err = cm.getImageBinaryData(img)
		if err != nil {
			return nil, err
		}
	} else {
		err = cm.saveImage(resp.TaskId, img)
		if err != nil {
			return nil, err
		}
	}

	// Register the answer in database.
	err = cm.registry.AddRecord(resp.TaskId, ringCount)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cm *CaptchaManager) CheckCaptcha(req *models.CheckCaptchaRequest) (resp *models.CheckCaptchaResponse, err error) {
	var ok bool
	ok, err = cm.registry.CheckCaptcha(req.TaskId, req.Value)
	if err != nil {
		return nil, err
	}

	resp = &models.CheckCaptchaResponse{
		TaskId:    req.TaskId,
		IsSuccess: ok,
	}

	return resp, nil
}

func (cm *CaptchaManager) createRandomUID() (uid string) {
	return "RCS-" + uuid.New().String()
}

func (cm *CaptchaManager) saveImage(uid string, img *image.NRGBA) (err error) {
	err = cos.SaveImageAsPngFile(img, makeRecordFilePath(cm.imagesFolder, uid))
	if err != nil {
		return err
	}

	return nil
}

func (cm *CaptchaManager) getImageBinaryData(img *image.NRGBA) (data []byte, err error) {
	buf := new(bytes.Buffer)

	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (cm *CaptchaManager) ClearJunk() (err error) {
	return cm.registry.ClearJunk()
}

func (cm *CaptchaManager) clearImagesFolder() (err error) {
	var items []os.DirEntry
	items, err = os.ReadDir(cm.imagesFolder)
	if err != nil {
		return err
	}

	for _, item := range items {
		err = os.RemoveAll(filepath.Join(cm.imagesFolder, item.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (cm *CaptchaManager) Stop() (err error) {
	err = cm.stopHttpServer()
	if err != nil {
		return err
	}

	log.Println("Captcha manager has been stopped.")

	return nil
}

func (cm *CaptchaManager) stopHttpServer() (err error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = cm.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CaptchaManager) isIdRegistered(id string) (isRegistered bool) {
	return cm.registry.IsIdRegistered(id)
}

func (cm *CaptchaManager) readFile(id string) (data []byte, err error) {
	return cm.registry.ReadFile(id)
}
