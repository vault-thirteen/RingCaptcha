package rc

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	ErrFDuplicateId = "duplicate id: %v"
	ErrFAbsentId    = "absent id: %v"
)

type Registry struct {
	storeImages  bool
	imagesFolder string
	recordTTLSec float64
	data         map[string]*RegistryRecord
	guard        *sync.Mutex
}

type RegistryRecord struct {
	Id     string
	Answer uint

	// Time of creation.
	ToC time.Time
}

func NewRegistry(
	storeImages bool,
	imagesFolder string,
	recordTTLSec uint,
) *Registry {
	return &Registry{
		storeImages:  storeImages,
		imagesFolder: imagesFolder,
		recordTTLSec: float64(recordTTLSec),
		data:         make(map[string]*RegistryRecord),
		guard:        new(sync.Mutex),
	}
}

func (r *Registry) AddRecord(id string, answer uint) (err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	_, alreadyExists := r.data[id]
	if alreadyExists {
		return fmt.Errorf(ErrFDuplicateId, id)
	}

	r.data[id] = &RegistryRecord{
		Id:     id,
		Answer: answer,
		ToC:    time.Now(),
	}

	return nil
}

func (r *Registry) GetSize() (size int) {
	r.guard.Lock()
	defer r.guard.Unlock()

	return len(r.data)
}

func (r *Registry) ClearJunk() (err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	now := time.Now()

	for id, rec := range r.data {
		if now.Sub(rec.ToC).Seconds() > r.recordTTLSec {
			delete(r.data, id)

			err = r.deleteRecordFile(id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Registry) deleteRecordFile(id string) (err error) {
	err = os.Remove(makeRecordFilePath(r.imagesFolder, id))
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) IsIdRegistered(id string) (isRegistered bool) {
	r.guard.Lock()
	defer r.guard.Unlock()

	_, isRegistered = r.data[id]

	return isRegistered
}

func (r *Registry) ReadFile(id string) (data []byte, err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	_, isRegistered := r.data[id]

	if !isRegistered {
		return nil, fmt.Errorf(ErrFAbsentId, id)
	}

	data, err = os.ReadFile(makeRecordFilePath(r.imagesFolder, id))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Registry) CheckCaptcha(id string, value uint) (ok bool, err error) {
	r.guard.Lock()
	defer r.guard.Unlock()

	rec, isRegistered := r.data[id]
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
	delete(r.data, id)

	if r.storeImages {
		err = r.deleteRecordFile(id)
		if err != nil {
			return err
		}
	} else {
		// We do not store anything, so nothing can be deleted.
	}

	return nil
}
