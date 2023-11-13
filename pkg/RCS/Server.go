package rcs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/osamingo/jsonrpc/v2"
	capman "github.com/vault-thirteen/RingCaptcha/pkg/RCS/CM"
	"github.com/vault-thirteen/RingCaptcha/pkg/RCS/settings"
)

type Server struct {
	settings *settings.Settings

	// HTTP server.
	listenDsn  string
	httpServer *http.Server

	// Channel for an external controller. When a message comes from this
	// channel, a controller must stop this server. The server does not stop
	// itself.
	mustBeStopped chan bool

	// Internal control structures.
	subRoutines *sync.WaitGroup
	mustStop    *atomic.Bool
	httpErrors  chan error

	// Captcha manager.
	captchaManager *capman.CaptchaManager

	// JSON-RPC handlers.
	jsonRpcHandlers *jsonrpc.MethodRepository

	// Diagnostic data.
	diag *DiagnosticData
}

func NewServer(stn *settings.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings:        stn,
		listenDsn:       fmt.Sprintf("%s:%d", stn.HttpSettings.Host, stn.HttpSettings.Port),
		mustBeStopped:   make(chan bool, 2),
		subRoutines:     new(sync.WaitGroup),
		mustStop:        new(atomic.Bool),
		httpErrors:      make(chan error, 8),
		jsonRpcHandlers: jsonrpc.NewMethodRepository(),
	}
	srv.mustStop.Store(false)

	err = srv.initJsonRpcHandlers()
	if err != nil {
		return nil, err
	}

	err = srv.initCaptchaManager()
	if err != nil {
		return nil, err
	}

	err = srv.initDiagnosticData()
	if err != nil {
		return nil, err
	}

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	return srv, nil
}

func (srv *Server) GetListenDsn() (dsn string) {
	return srv.listenDsn
}

func (srv *Server) GetCaptchaManagerListenDsn() (dsn string) {
	return srv.captchaManager.GetListenDsn()
}

func (srv *Server) GetStopChannel() *chan bool {
	return &srv.mustBeStopped
}

func (srv *Server) Start() (err error) {
	srv.startHttpServer()

	srv.subRoutines.Add(2)
	go srv.listenForHttpErrors()
	go srv.clearJunk()

	return nil
}

func (srv *Server) Stop() (err error) {
	srv.mustStop.Store(true)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	close(srv.httpErrors)

	srv.subRoutines.Wait()

	err = srv.captchaManager.Stop()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) startHttpServer() {
	go func() {
		var listenError error
		listenError = srv.httpServer.ListenAndServe()

		if (listenError != nil) && (listenError != http.ErrServerClosed) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) listenForHttpErrors() {
	defer srv.subRoutines.Done()

	for err := range srv.httpErrors {
		log.Println("Server error: " + err.Error())
		srv.mustBeStopped <- true
	}

	log.Println("HTTP error listener has stopped.")
}

func (srv *Server) clearJunk() {
	defer srv.subRoutines.Done()

	var err error
	var imageTTL = int(srv.settings.CaptchaSettings.ImageTTLSec)

main_loop:
	for {
		// Check & sleep.
		for i := 0; i < imageTTL; i++ {
			time.Sleep(time.Second)

			if srv.mustStop.Load() {
				break main_loop
			}
		}

		// Work.
		err = srv.captchaManager.ClearJunk()
		if err != nil {
			log.Println(err)
			srv.mustBeStopped <- true
		}
	}

	log.Println("Junk cleaner has stopped.")
}

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	srv.jsonRpcHandlers.ServeHTTP(rw, req)
}

func (srv *Server) initCaptchaManager() (err error) {
	srv.captchaManager, err = capman.NewCaptchaManager(
		srv.settings.CaptchaSettings.StoreImages,
		srv.settings.CaptchaSettings.ImagesFolder,
		srv.settings.CaptchaSettings.ImageWidth,
		srv.settings.CaptchaSettings.ImageHeight,
		srv.settings.CaptchaSettings.ImageTTLSec,
		srv.settings.CaptchaSettings.ClearImagesFolderAtStart,
		srv.settings.CaptchaSettings.UseHttpServerForImages,
		srv.settings.CaptchaSettings.HttpServerHost,
		srv.settings.CaptchaSettings.HttpServerPort,
		&srv.httpErrors,
		srv.settings.CaptchaSettings.HttpServerName,
	)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) initDiagnosticData() (err error) {
	srv.diag = &DiagnosticData{}

	return nil
}
