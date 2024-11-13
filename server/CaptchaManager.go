package server

import (
	"bytes"
	"context"
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
	"github.com/vault-thirteen/RingCaptcha"
	"github.com/vault-thirteen/RingCaptcha/server/models"
	hdr "github.com/vault-thirteen/auxie/header"
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
	// Settings.
	settings *models.CaptchaManagerSettings

	// Storage guard.
	sg *sync.Mutex

	// Registry.
	registry *Registry

	// HTTP server.
	httpServer     *http.Server
	listenDsn      string
	httpErrorsChan *chan error
	httpServerName string
}

func NewCaptchaManager(s *models.CaptchaManagerSettings) (cm *CaptchaManager, err error) {
	err = s.Check()
	if err != nil {
		return nil, err
	}

	cm = &CaptchaManager{
		settings: s,
		sg:       new(sync.Mutex),
	}

	cm.registry, err = NewRegistry(
		cm.settings.IsImageStorageUsed,
		cm.settings.ImagesFolder,
		cm.settings.ImageTtlSec,
		cm.settings.IsCachingEnabled,
		cm.settings.CacheSizeLimit,
		cm.settings.CacheVolumeLimit,
		cm.settings.CacheRecordTtl,
	)
	if err != nil {
		return nil, err
	}

	if cm.settings.IsImageCleanupAtStartUsed {
		fmt.Print(MsgCleaningImagesFolder)
		err = cm.clearImagesFolder()
		if err != nil {
			return nil, err
		}
		fmt.Println(MsgDone)
	}

	if cm.settings.IsImageServerEnabled {
		err = cm.initHttpServer(
			cm.settings.HttpHost,
			cm.settings.HttpPort,
			cm.settings.HttpErrorsChan,
			cm.settings.HttpServerName,
		)
		if err != nil {
			return nil, err
		}
	}

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

	return nil
}

func (cm *CaptchaManager) Start() (err error) {
	go func() {
		var listenError error
		listenError = cm.httpServer.ListenAndServe()

		if (listenError != nil) && (listenError != http.ErrServerClosed) {
			*cm.httpErrorsChan <- listenError
		}
	}()

	log.Println(MsgCaptchaManagerStart)

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

	if !cm.settings.IsImageStorageUsed {
		// No images are stored. Sorry.
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	var fileContents []byte
	var err error
	fileContents, err = cm.getImageFile(id)
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

func (cm *CaptchaManager) CreateCaptcha() (resp *models.CreateCaptchaResponse, err error) {
	var img *image.NRGBA
	var ringCount uint
	img, ringCount, err = rc.CreateCaptchaImage(cm.settings.ImageWidth, cm.settings.ImageHeight, true, false)
	if err != nil {
		return nil, err
	}

	resp = &models.CreateCaptchaResponse{
		TaskId:              cm.createRandomUID(),
		ImageFormat:         ImageFormat,
		IsImageDataReturned: !cm.settings.IsImageStorageUsed,
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

	// Register the answer.
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
	return RUIDPrefix + uuid.New().String()
}

func (cm *CaptchaManager) saveImage(uid string, img *image.NRGBA) (err error) {
	cm.sg.Lock()
	defer cm.sg.Unlock()

	err = rc.SaveImageAsPngFile(img, MakeRecordFilePath(cm.settings.ImagesFolder, uid))
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
	items, err = os.ReadDir(cm.settings.ImagesFolder)
	if err != nil {
		return err
	}

	for _, item := range items {
		err = os.RemoveAll(filepath.Join(cm.settings.ImagesFolder, item.Name()))
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

func (cm *CaptchaManager) getImageFile(id string) (data []byte, err error) {
	return cm.registry.GetImageFile(id)
}
