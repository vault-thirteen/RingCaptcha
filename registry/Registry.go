package r

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vault-thirteen/Cache/NVL"
	"github.com/vault-thirteen/RingCaptcha/creator"
	"github.com/vault-thirteen/RingCaptcha/models"
	"github.com/vault-thirteen/Simple-File-Server"
)

type Registry struct {
	// Settings.
	settings *m.RegistrySettings

	// Data structures.
	cache      *nvl.Cache[string, *RegistryRecord]
	fileServer *sfs.SimpleFileServer
	guard      *sync.Mutex

	// Control structures.
	cleanerWG *sync.WaitGroup
	mustStop  atomic.Bool
}

func NewRegistry(s *m.RegistrySettings) (r *Registry, err error) {
	r = &Registry{
		settings:  s,
		cache:     nvl.NewCache[string, *RegistryRecord](s.RecordCacheSizeLimit, s.RecordCacheItemTtl),
		guard:     new(sync.Mutex),
		cleanerWG: new(sync.WaitGroup),
	}

	if s.IsImageStorageUsed {
		r.fileServer, err = sfs.NewSimpleFileServer(
			r.settings.ImagesFolder,
			[]string{},
			true,
			r.settings.FileCacheSizeLimit,
			r.settings.FileCacheVolumeLimit,
			r.settings.FileCacheItemTtl,
		)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Registry) Start() {
	r.mustStop.Store(false)

	if r.settings.IsStorageCleaningEnabled {
		r.cleanerWG.Add(1)
		go r.runCleaner()
	}
}
func (r *Registry) Stop() {
	r.mustStop.Store(true)

	if r.settings.IsStorageCleaningEnabled {
		r.cleanerWG.Wait()
	}
}

func (r *Registry) CreateCaptcha(captcha *m.Captcha) (err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	if r.cache.RecordExists(captcha.Id) {
		return m.NewErrorWithHttpStatusCode(m.Err_IdIsDuplicate, http.StatusConflict)
	}

	err = r.cache.AddRecord(captcha.Id, NewRegistryRecord(captcha.Id, captcha.RingCount, r.settings.RecordCacheItemTtl))
	if err != nil {
		return err
	}

	if r.settings.IsImageStorageUsed {
		err = r.fileServer.CreateFile(c.MakeFileName(captcha.Id), captcha.ImageData)
		if err != nil {
			return err
		}
	}

	return nil
}
func (r *Registry) CheckCaptcha(captcha *m.Captcha) (ok bool, err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	if !r.cache.RecordExists(captcha.Id) {
		return false, m.NewErrorWithHttpStatusCode(m.Err_IdIsNotFound, http.StatusNotFound)
	}

	var rr *RegistryRecord
	rr, err = r.cache.GetRecord(captcha.Id)
	if err != nil {
		return false, err
	}

	// N.B. Captcha is deleted after a first guess.
	// We do not instantly erase images from storage to reduce the wear of
	// hardware devices. Outdated images are periodically cleaned by the
	// cleaner.
	r.cache.RemoveRecord(captcha.Id)

	if rr.Id != captcha.Id {
		return false, errors.New(m.Err_Anomaly)
	}

	if captcha.RingCount != rr.Answer {
		return false, nil
	}

	return true, nil
}
func (r *Registry) HasCaptcha(captcha *m.Captcha) (exists bool, err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	if !r.cache.RecordExists(captcha.Id) {
		return false, nil
	}

	return true, nil
}
func (r *Registry) GetCaptchaImage(captcha *m.Captcha) (data []byte, err error) {
	if !r.settings.IsImageStorageUsed {
		return nil, m.NewErrorWithHttpStatusCode(m.Err_FileStorageIsDisabled, http.StatusBadRequest)
	}

	r.guard.Lock()
	defer r.guard.Unlock()

	if !r.cache.RecordExists(captcha.Id) {
		return nil, m.NewErrorWithHttpStatusCode(m.Err_IdIsNotFound, http.StatusNotFound)
	}

	return r.fileServer.GetFile(c.MakeFileName(captcha.Id))
}

func (r *Registry) runCleaner() {
	defer r.cleanerWG.Done()

	fmt.Println(m.Msg_ImageCleanerHasStarted)

	var ts = 0
	var err error
	for {
		if r.mustStop.Load() {
			break
		}

		// Each minute do ...
		if ts == 60 {
			ts = 0

			err = r.cleanOutdatedFiles()
			if err != nil {
				log.Println(err)
			}
		}

		// Next tick.
		time.Sleep(5 * time.Second)
		ts += 5
	}

	fmt.Println(m.Msg_ImageCleanerHasStopped)
}
func (r *Registry) cleanOutdatedFiles() (err error) {
	var filesCount int
	filesCount, err = r.fileServer.CountFiles(".")
	if err != nil {
		return err
	}

	if filesCount < r.settings.FilesCountToClean {
		return nil
	}

	var fileNames []string
	fileNames, err = r.fileServer.ListFileNames(".")
	if err != nil {
		return err
	}

	var id string
	var exists bool
	var filesToDelete = []string{}
	for _, fileName := range fileNames {
		if filepath.Ext(fileName) != m.FileExtFullPng {
			continue
		}

		id = c.FileNameWithoutExtension(fileName)

		exists, err = r.HasCaptcha(m.NewCaptchaWithId(id))
		if err != nil {
			return err
		}

		if !exists {
			filesToDelete = append(filesToDelete, fileName)
		}
	}

	if len(filesToDelete) < r.settings.FilesCountToClean {
		return nil
	}

	log.Println(fmt.Sprintf(m.MsgF_CleaningImages, len(filesToDelete)))

	for _, fileToDelete := range filesToDelete {
		err = r.fileServer.DeleteFile(fileToDelete)
		if err != nil {
			return err
		}
	}

	return nil
}
