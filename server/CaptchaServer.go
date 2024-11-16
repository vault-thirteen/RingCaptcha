package s

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/vault-thirteen/RingCaptcha/creator"
	"github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/RingCaptcha/registry"
	"github.com/vault-thirteen/auxie/header"
)

type CaptchaServer struct {
	// Settings.
	settings *m.CaptchaServerSettings

	// Storage guard.
	sg *sync.Mutex

	// Registry.
	registry *r.Registry

	// HTTP server.
	httpServer     *http.Server
	listenDsn      string
	httpErrorsChan *chan error
	httpServerName string
}

func NewCaptchaServer(s *m.CaptchaServerSettings) (cm *CaptchaServer, err error) {
	err = s.Check()
	if err != nil {
		return nil, err
	}

	cm = &CaptchaServer{
		settings: s,
		sg:       new(sync.Mutex),
	}

	rs := &m.RegistrySettings{
		// Main settings.
		IsImageStorageUsed:       cm.settings.IsImageStorageUsed,
		IsStorageCleaningEnabled: cm.settings.IsStorageCleaningEnabled,

		// Image settings.
		ImagesFolder:      cm.settings.ImagesFolder,
		FilesCountToClean: cm.settings.FilesCountToClean,

		// File cache settings.
		FileCacheSizeLimit:   cm.settings.FileCacheSizeLimit,
		FileCacheVolumeLimit: cm.settings.FileCacheVolumeLimit,
		FileCacheItemTtl:     cm.settings.FileCacheItemTtl,

		// Record cache settings.
		RecordCacheSizeLimit: cm.settings.RecordCacheSizeLimit,
		RecordCacheItemTtl:   cm.settings.RecordCacheItemTtl,
	}

	cm.registry, err = r.NewRegistry(rs)
	if err != nil {
		return nil, err
	}

	if cm.settings.IsImageCleanupAtStartUsed {
		err = cm.clearImagesFolder()
		if err != nil {
			return nil, err
		}
	}

	if cm.settings.IsImageServerEnabled {
		err = cm.initHttpServer()
		if err != nil {
			return nil, err
		}
	}

	return cm, nil
}

func (cs *CaptchaServer) Start() (err error) {
	cs.registry.Start()

	if cs.settings.IsImageServerEnabled {
		cs.startHttpServer()
	}

	log.Println(m.Msg_CaptchaManagerStart)

	return nil
}
func (cs *CaptchaServer) Stop() (err error) {
	if cs.settings.IsImageServerEnabled {
		err = cs.stopHttpServer()
		if err != nil {
			return err
		}
	}

	cs.registry.Stop()

	log.Println(m.Msg_CaptchaManagerStop)

	return nil
}
func (cs *CaptchaServer) GetListenDsn() (dsn string) {
	return cs.listenDsn
}

func (cs *CaptchaServer) CreateCaptcha() (resp *m.CreateCaptchaResponse, err error) {
	var captcha *m.Captcha
	captcha, err = c.CreateCaptcha(cs.settings.ImageWidth, cs.settings.ImageHeight)
	if err != nil {
		return nil, err
	}

	err = cs.registry.CreateCaptcha(captcha)
	if err != nil {
		return nil, err
	}

	return cs.composeResponseForCreateCaptcha(captcha)
}
func (cs *CaptchaServer) CheckCaptcha(req *m.CheckCaptchaRequest) (resp *m.CheckCaptchaResponse, err error) {
	err = req.Check()
	if err != nil {
		return nil, err
	}

	var captcha = m.NewCaptchaWithIdAndAnswer(req.TaskId, req.Value)

	resp = &m.CheckCaptchaResponse{TaskId: req.TaskId}
	resp.IsSuccess, err = cs.registry.CheckCaptcha(captcha)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func (cs *CaptchaServer) HasCaptcha(req *m.HasCaptchaRequest) (resp *m.HasCaptchaResponse, err error) {
	err = req.Check()
	if err != nil {
		return nil, err
	}

	var captcha = m.NewCaptchaWithId(req.TaskId)

	resp = &m.HasCaptchaResponse{TaskId: req.TaskId}
	resp.IsFound, err = cs.registry.HasCaptcha(captcha)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
func (cs *CaptchaServer) GetCaptchaImage(req *m.GetCaptchaImageRequest) (resp *m.GetCaptchaImageResponse, err error) {
	if !cs.settings.IsImageStorageUsed {
		return nil, m.NewErrorWithHttpStatusCode(m.Err_FileStorageIsDisabled, http.StatusForbidden)
	}

	err = req.Check()
	if err != nil {
		return nil, err
	}

	var captcha = m.NewCaptchaWithId(req.TaskId)

	resp = &m.GetCaptchaImageResponse{TaskId: req.TaskId}
	resp.ImageData, err = cs.registry.GetCaptchaImage(captcha)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (cs *CaptchaServer) clearImagesFolder() (err error) {
	fmt.Print(m.Msg_CleaningImagesFolder)
	defer func() {
		if err != nil {
			fmt.Println(m.Msg_Failure)
		} else {
			fmt.Println(m.Msg_Done)
		}
	}()

	imgFolder := cs.settings.ImagesFolder

	var items []os.DirEntry
	items, err = os.ReadDir(imgFolder)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if filepath.Ext(item.Name()) != m.FileExtFullPng {
			continue
		}

		err = os.Remove(filepath.Join(imgFolder, item.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
func (cs *CaptchaServer) initHttpServer() (err error) {
	host := cs.settings.HttpHost
	port := cs.settings.HttpPort
	cs.listenDsn = net.JoinHostPort(host, strconv.FormatUint(uint64(port), 10))

	cs.httpServer = &http.Server{
		Addr:    cs.listenDsn,
		Handler: http.Handler(http.HandlerFunc(cs.httpRouter)),
	}

	cs.httpErrorsChan = cs.settings.HttpErrorsChan
	cs.httpServerName = cs.settings.HttpServerName

	return nil
}
func (cs *CaptchaServer) startHttpServer() {
	go func() {
		var listenError error
		listenError = cs.httpServer.ListenAndServe()

		if (listenError != nil) && (listenError != http.ErrServerClosed) {
			*cs.httpErrorsChan <- listenError
		}
	}()
}
func (cs *CaptchaServer) stopHttpServer() (err error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = cs.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (cs *CaptchaServer) composeResponseForCreateCaptcha(c *m.Captcha) (resp *m.CreateCaptchaResponse, err error) {
	resp = &m.CreateCaptchaResponse{
		TaskId:              c.Id,
		ImageFormat:         c.ImageFormat,
		IsImageDataReturned: !cs.settings.IsImageStorageUsed,
	}

	if resp.IsImageDataReturned {
		resp.ImageData = c.ImageData
	}

	return resp, nil
}

func (cs *CaptchaServer) httpRouter(rw http.ResponseWriter, httpReq *http.Request) {
	cs.httpHandler_GetCaptchaImage(rw, httpReq)
}
func (cs *CaptchaServer) httpHandler_GetCaptchaImage(rw http.ResponseWriter, httpReq *http.Request) {
	req, err := m.NewGetImageRequestFromHttpRequest(httpReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp *m.GetCaptchaImageResponse
	resp, err = cs.GetCaptchaImage(req)
	if err != nil {
		cs.httpHandler_ProcessError(err, rw)
		return
	}

	rw.Header().Set(header.HttpHeaderContentType, m.ImageMimeType)
	rw.Header().Set(header.HttpHeaderServer, cs.httpServerName)
	rw.WriteHeader(http.StatusOK)

	_, err = rw.Write(resp.ImageData)
	if err != nil {
		log.Println(err)
	}
	return
}
func (cs *CaptchaServer) httpHandler_ProcessError(err error, rw http.ResponseWriter) {
	var ewhsc *m.ErrorWithHttpStatusCode
	if errors.As(err, &ewhsc) {
		rw.WriteHeader(ewhsc.GetHttpStatusCode())
		return
	}

	log.Println(err)
	rw.WriteHeader(http.StatusInternalServerError)
	return
}
