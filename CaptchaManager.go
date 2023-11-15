package rc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
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
	QueryKeyId    = "id"
	RUIDPrefix    = "RCS-"
)

const (
	MsgCleaningImagesFolder = "Cleaning the images folder ... "
	MsgDone                 = "Done"
	MsgCaptchaManagerStart  = "Captcha manager has started"
	MsgCaptchaManagerStop   = "Captcha manager has been stopped"
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

	// Storage guard.
	sg *sync.Mutex
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
		sg:                       new(sync.Mutex),
	}

	if cm.clearImagesFolderAtStart {
		fmt.Print(MsgCleaningImagesFolder)
		err = cm.clearImagesFolder()
		if err != nil {
			return nil, err
		}
		fmt.Println(MsgDone)
	}

	if cm.useHttpServerForImages {
		err = cm.initHttpServer(httpHost, httpPort, httpErrorsChan, httpServerName)
		if err != nil {
			return nil, err
		}
	}

	log.Println(MsgCaptchaManagerStart)

	return cm, nil
}

func (cm *CaptchaManager) initHttpServer(
	httpHost string,
	httpPort uint16,
	httpErrorsChan *chan error,
	httpServerName string,
) (err error) {
	cm.listenDsn = net.JoinHostPort(httpHost, strconv.FormatUint(uint64(httpPort), 10))
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
	id := req.URL.Query().Get(QueryKeyId)

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
		log.Println(err)
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

func (cm *CaptchaManager) CreateCaptcha() (resp *CreateCaptchaResponse, err error) {
	var img *image.NRGBA
	var ringCount uint
	img, ringCount, err = CreateCaptchaImage(cm.imageWidth, cm.imageHeight, true, false)
	if err != nil {
		return nil, err
	}

	resp = &CreateCaptchaResponse{
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

func (cm *CaptchaManager) CheckCaptcha(req *CheckCaptchaRequest) (resp *CheckCaptchaResponse, err error) {
	var ok bool
	ok, err = cm.registry.CheckCaptcha(req.TaskId, req.Value)
	if err != nil {
		return nil, err
	}

	resp = &CheckCaptchaResponse{
		TaskId:    req.TaskId,
		IsSuccess: ok,
	}

	return resp, nil
}

func (cm *CaptchaManager) createRandomUID() (uid string) {
	return RUIDPrefix + uuid.New().String()
}

func (cm *CaptchaManager) saveImage(uid string, img *image.NRGBA) (err error) {
	cm.sg.Lock()
	defer cm.sg.Unlock()

	err = SaveImageAsPngFile(img, makeRecordFilePath(cm.imagesFolder, uid))
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

	log.Println(MsgCaptchaManagerStop)

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
