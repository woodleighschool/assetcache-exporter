// Package assetcache reads Apple Content Caching status and interval metrics.
package assetcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

const defaultStatusCommand = "/usr/bin/AssetCacheManagerUtil"

// Status is the stable subset of AssetCacheManagerUtil status used by the exporter.
type Status struct {
	Activated               bool
	Active                  bool
	CacheUsedBytes          float64
	CacheLimitBytes         float64
	CacheFreeBytes          float64
	PersonalCacheUsedBytes  float64
	PersonalCacheLimitBytes float64
	PersonalCacheFreeBytes  float64
	MaxCachePressureRatio   float64
	CacheStatus             string
	StartupStatus           string
	RegistrationStatus      int
	ParentCount             int
	PeerCount               int
	TotalsStartTime         time.Time
	BytesDroppedTotal       float64
	BytesImportedTotal      float64
	BytesServedToChildren   float64
	BytesServedToClients    float64
	BytesServedToPeers      float64
	BytesStoredFromOrigin   float64
	BytesStoredFromParents  float64
	BytesStoredFromPeers    float64
}

// StatusSource reads status from AssetCacheManagerUtil.
type StatusSource struct {
	command string
}

// NewStatusSource returns a source using the system AssetCacheManagerUtil.
func NewStatusSource() StatusSource {
	return StatusSource{command: defaultStatusCommand}
}

// Read returns the current Content Caching status.
func (s StatusSource) Read(ctx context.Context) (Status, error) {
	command := s.command
	if command == "" {
		command = defaultStatusCommand
	}

	output, err := exec.CommandContext(ctx, command, "-j", "status").Output()
	if err != nil {
		return Status{}, fmt.Errorf("run AssetCacheManagerUtil: %w", err)
	}
	status, err := parseStatus(output)
	if err != nil {
		return Status{}, fmt.Errorf("parse AssetCacheManagerUtil status: %w", err)
	}
	return status, nil
}

type statusEnvelope struct {
	Result *statusPayload `json:"result"`
}

type statusPayload struct {
	Activated                    *bool   `json:"Activated"`
	Active                       *bool   `json:"Active"`
	CacheUsed                    float64 `json:"CacheUsed"`
	CacheLimit                   float64 `json:"CacheLimit"`
	CacheFree                    float64 `json:"CacheFree"`
	PersonalCacheUsed            float64 `json:"PersonalCacheUsed"`
	PersonalCacheLimit           float64 `json:"PersonalCacheLimit"`
	PersonalCacheFree            float64 `json:"PersonalCacheFree"`
	MaxCachePressureLast1Hour    float64 `json:"MaxCachePressureLast1Hour"`
	CacheStatus                  string  `json:"CacheStatus"`
	StartupStatus                string  `json:"StartupStatus"`
	RegistrationStatus           int     `json:"RegistrationStatus"`
	Parents                      []any   `json:"Parents"`
	Peers                        []any   `json:"Peers"`
	TotalBytesAreSince           string  `json:"TotalBytesAreSince"`
	TotalBytesDropped            float64 `json:"TotalBytesDropped"`
	TotalBytesImported           float64 `json:"TotalBytesImported"`
	TotalBytesReturnedToChildren float64 `json:"TotalBytesReturnedToChildren"`
	TotalBytesReturnedToClients  float64 `json:"TotalBytesReturnedToClients"`
	TotalBytesReturnedToPeers    float64 `json:"TotalBytesReturnedToPeers"`
	TotalBytesStoredFromOrigin   float64 `json:"TotalBytesStoredFromOrigin"`
	TotalBytesStoredFromParents  float64 `json:"TotalBytesStoredFromParents"`
	TotalBytesStoredFromPeers    float64 `json:"TotalBytesStoredFromPeers"`
}

func parseStatus(data []byte) (Status, error) {
	var envelope statusEnvelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return Status{}, err
	}
	if envelope.Result == nil {
		return Status{}, errors.New("response has no result")
	}

	payload := envelope.Result
	if payload.Activated == nil || payload.Active == nil {
		return Status{}, errors.New("result has no Activated or Active status")
	}
	var totalsStartTime time.Time
	if payload.TotalBytesAreSince != "" {
		var err error
		totalsStartTime, err = time.Parse("2006-01-02 15:04:05 -0700", payload.TotalBytesAreSince)
		if err != nil {
			return Status{}, fmt.Errorf("parse TotalBytesAreSince: %w", err)
		}
	}

	return Status{
		Activated:               *payload.Activated,
		Active:                  *payload.Active,
		CacheUsedBytes:          payload.CacheUsed,
		CacheLimitBytes:         payload.CacheLimit,
		CacheFreeBytes:          payload.CacheFree,
		PersonalCacheUsedBytes:  payload.PersonalCacheUsed,
		PersonalCacheLimitBytes: payload.PersonalCacheLimit,
		PersonalCacheFreeBytes:  payload.PersonalCacheFree,
		MaxCachePressureRatio:   payload.MaxCachePressureLast1Hour / 100,
		CacheStatus:             payload.CacheStatus,
		StartupStatus:           payload.StartupStatus,
		RegistrationStatus:      payload.RegistrationStatus,
		ParentCount:             len(payload.Parents),
		PeerCount:               len(payload.Peers),
		TotalsStartTime:         totalsStartTime,
		BytesDroppedTotal:       payload.TotalBytesDropped,
		BytesImportedTotal:      payload.TotalBytesImported,
		BytesServedToChildren:   payload.TotalBytesReturnedToChildren,
		BytesServedToClients:    payload.TotalBytesReturnedToClients,
		BytesServedToPeers:      payload.TotalBytesReturnedToPeers,
		BytesStoredFromOrigin:   payload.TotalBytesStoredFromOrigin,
		BytesStoredFromParents:  payload.TotalBytesStoredFromParents,
		BytesStoredFromPeers:    payload.TotalBytesStoredFromPeers,
	}, nil
}
