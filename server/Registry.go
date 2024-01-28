package server

import (
	"errors"
	"fmt"
	"os"
	"sync"

	sfs "github.com/vault-thirteen/Simple-File-Server"
)

const (
	ErrFDuplicateId   = "duplicate id: %v"
	ErrFAbsentId      = "absent id: %v"
	ErrFileIsNotFound = "file is not found"
)

type Registry struct {
	storeImages  bool
	imagesFolder string
	recordTtlSec float64
	records      map[string]*RegistryRecord
	guard        *sync.Mutex
	fileServer   *sfs.SimpleFileServer
}

func NewRegistry(
	storeImages bool,
	imagesFolder string,
	recordTtlSec uint,
	isCachingEnabled bool,
	cacheSizeLimit int,
	cacheVolumeLimit int,
	cacheRecordTtl uint,
) (r *Registry, err error) {
	r = &Registry{
		storeImages:  storeImages,
		imagesFolder: imagesFolder,
		recordTtlSec: float64(recordTtlSec),
		records:      make(map[string]*RegistryRecord),
		guard:        new(sync.Mutex),
	}

	r.fileServer, err = sfs.NewSimpleFileServer(
		imagesFolder,
		[]string{},
		isCachingEnabled,
		cacheSizeLimit,
		cacheVolumeLimit,
		cacheRecordTtl,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Registry) IsIdRegistered(id string) (isRegistered bool) {
	r.guard.Lock()
	defer r.guard.Unlock()

	return r.isIdRegistered(id)
}

func (r *Registry) isIdRegistered(id string) (isRegistered bool) {
	_, isRegistered = r.records[id]
	return isRegistered
}

func (r *Registry) AddRecord(id string, answer uint) (err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	if r.isIdRegistered(id) {
		return fmt.Errorf(ErrFDuplicateId, id)
	}

	r.records[id] = NewRegistryRecord(id, answer, r.recordTtlSec)

	return nil
}

func (r *Registry) GetSize() (size int) {
	r.guard.Lock()
	defer r.guard.Unlock()

	return len(r.records)
}

func (r *Registry) ClearJunk() (err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	for id, rec := range r.records {
		if rec.IsExpired() {
			delete(r.records, id)

			_ = r.fileServer.ForgetFile(MakeFileName(id))

			if r.storeImages {
				err = r.deleteRecordFile(id)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *Registry) deleteRecordFile(id string) (err error) {
	return os.Remove(MakeRecordFilePath(r.imagesFolder, id))
}

func (r *Registry) GetImageFile(id string) (data []byte, err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	if !r.isIdRegistered(id) {
		return nil, fmt.Errorf(ErrFAbsentId, id)
	}

	var fileExists bool
	data, fileExists, err = r.fileServer.GetFile(MakeFileName(id))
	if err != nil {
		return nil, err
	}
	if !fileExists {
		return nil, errors.New(ErrFileIsNotFound)
	}

	return data, nil
}

func (r *Registry) CheckCaptcha(id string, value uint) (ok bool, err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	rec, isRegistered := r.records[id]
	if !isRegistered {
		return false, fmt.Errorf(ErrFAbsentId, id)
	}

	// Captcha is deleted after a first guess.
	err = r.removeCaptcha(id)
	if err != nil {
		return false, err
	}

	if value != rec.Answer {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *Registry) removeCaptcha(id string) (err error) {
	delete(r.records, id)

	_ = r.fileServer.ForgetFile(MakeFileName(id))

	if r.storeImages {
		err = r.deleteRecordFile(id)
		if err != nil {
			return err
		}
	}

	return nil
}
