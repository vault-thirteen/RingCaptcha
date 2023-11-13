package rcs

import "sync/atomic"

type DiagnosticData struct {
	// Number of all incoming requests.
	totalRequestsCount atomic.Uint64

	// Number of successfully finished requests.
	successfulRequestsCount atomic.Uint64
}

func (dd *DiagnosticData) getTotalRequestsCount() (trc uint64) {
	return dd.totalRequestsCount.Load()
}

func (dd *DiagnosticData) incTotalRequestsCount() {
	dd.totalRequestsCount.Add(1)
}

func (dd *DiagnosticData) getSuccessfulRequestsCount() (src uint64) {
	return dd.successfulRequestsCount.Load()
}

func (dd *DiagnosticData) incSuccessfulRequestsCount() {
	dd.successfulRequestsCount.Add(1)
}
