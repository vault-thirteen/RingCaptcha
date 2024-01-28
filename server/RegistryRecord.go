package server

import "time"

// RegistryRecord stores an answer and timing information.
type RegistryRecord struct {
	Id     string
	Answer uint
	TTL    float64   // TTL is set in seconds.
	ToC    time.Time // Time of creation.
}

func NewRegistryRecord(id string, answer uint, ttl float64) *RegistryRecord {
	return &RegistryRecord{
		Id:     id,
		Answer: answer,
		TTL:    ttl,
		ToC:    time.Now(),
	}
}

func (rr *RegistryRecord) IsExpired() bool {
	return time.Now().Sub(rr.ToC).Seconds() > rr.TTL
}
