package assetcache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/url"
	"time"

	_ "modernc.org/sqlite" // Register the pure-Go SQLite driver used by database/sql.
)

const (
	defaultMetricsPath = "/Library/Application Support/Apple/AssetCache/Metrics/Metrics.db"
	coreDataEpoch      = 978307200
)

// Observation is a labelled value from the latest metrics interval.
type Observation struct {
	Source string
	Value  float64
}

// Interval contains the latest immutable row from Apple's metrics database.
type Interval struct {
	Timestamp               time.Time
	PeriodSeconds           float64
	BytesServed             []Observation
	Replies                 []Observation
	Requests                []Observation
	BytesImported           float64
	Imports                 float64
	BytesDropped            float64
	BytesPurged             float64
	RequestsRejectedNoSpace float64
	CachePressureRatio      float64
}

// MetricsSource reads Apple's metrics database in read-only mode.
type MetricsSource struct {
	path string
}

// NewMetricsSource returns a source using the system Content Caching metrics database.
func NewMetricsSource() MetricsSource {
	return MetricsSource{path: defaultMetricsPath}
}

// Read returns the latest metrics interval. A nil interval means the database is readable but empty.
func (s MetricsSource) Read(ctx context.Context) (*Interval, error) {
	path := s.path
	if path == "" {
		path = defaultMetricsPath
	}

	databaseURL := url.URL{Scheme: "file", Path: path}
	query := databaseURL.Query()
	query.Set("mode", "ro")
	databaseURL.RawQuery = query.Encode()

	database, err := sql.Open("sqlite", databaseURL.String())
	if err != nil {
		return nil, fmt.Errorf("open metrics database: %w", err)
	}
	database.SetMaxOpenConns(1)

	interval, err := readLatestInterval(ctx, database)
	closeErr := database.Close()
	if err != nil {
		return nil, fmt.Errorf("read metrics database: %w", err)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close metrics database: %w", closeErr)
	}
	return interval, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func readLatestInterval(ctx context.Context, database *sql.DB) (*Interval, error) {
	row := database.QueryRowContext(ctx, `
		SELECT
			ZCREATIONDATE + ?,
			COALESCE(ZPERIOD, 0),
			COALESCE(ZBYTESFROMCACHETOCHILD, 0) + COALESCE(ZBYTESFROMCACHETOCLIENT, 0) + COALESCE(ZBYTESFROMCACHETOPEER, 0),
			COALESCE(ZBYTESFROMORIGINTOCHILD, 0) + COALESCE(ZBYTESFROMORIGINTOCLIENT, 0) + COALESCE(ZBYTESFROMORIGINTOPEER, 0),
			COALESCE(ZBYTESFROMPARENTTOCHILD, 0) + COALESCE(ZBYTESFROMPARENTTOCLIENT, 0) + COALESCE(ZBYTESFROMPARENTTOPEER, 0),
			COALESCE(ZBYTESFROMPEERTOCHILD, 0) + COALESCE(ZBYTESFROMPEERTOCLIENT, 0),
			COALESCE(ZREPLIESFROMCACHETOCHILD, 0) + COALESCE(ZREPLIESFROMCACHETOCLIENT, 0) + COALESCE(ZREPLIESFROMCACHETOPEER, 0),
			COALESCE(ZREPLIESFROMORIGINTOCHILD, 0) + COALESCE(ZREPLIESFROMORIGINTOCLIENT, 0) + COALESCE(ZREPLIESFROMORIGINTOPEER, 0),
			COALESCE(ZREPLIESFROMPARENTTOCHILD, 0) + COALESCE(ZREPLIESFROMPARENTTOCLIENT, 0) + COALESCE(ZREPLIESFROMPARENTTOPEER, 0),
			COALESCE(ZREPLIESFROMPEERTOCHILD, 0) + COALESCE(ZREPLIESFROMPEERTOCLIENT, 0),
			COALESCE(ZREQUESTSFROMCHILD, 0),
			COALESCE(ZREQUESTSFROMCLIENT, 0),
			COALESCE(ZREQUESTSFROMPEER, 0),
			COALESCE(ZBYTESIMPORTEDBYHTTP, 0) + COALESCE(ZBYTESIMPORTEDBYXPC, 0),
			COALESCE(ZIMPORTSBYHTTP, 0) + COALESCE(ZIMPORTSBYXPC, 0),
			COALESCE(ZBYTESDROPPED, 0),
			COALESCE(ZBYTESPURGEDTOTAL, 0),
			COALESCE(ZREQUESTSREJECTEDFORNOSPACE, 0),
			CASE
				WHEN COALESCE(ZREQUESTSREJECTEDFORNOSPACE, 0) > 0 THEN 1.0
				WHEN COALESCE(ZBYTESPURGEDYOUNGERTHAN1DAY, 0) > 0 THEN 0.8
				WHEN COALESCE(ZBYTESPURGEDYOUNGERTHAN7DAYS, 0) > 0 THEN 0.6
				WHEN COALESCE(ZBYTESPURGEDYOUNGERTHAN30DAYS, 0) > 0 THEN 0.4
				WHEN COALESCE(ZBYTESPURGEDTOTAL, 0) > 0 THEN 0.2
				ELSE 0.0
			END
		FROM ZMETRIC
		ORDER BY ZCREATIONDATE DESC
		LIMIT 1
	`, coreDataEpoch)

	interval, err := scanInterval(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return interval, nil
}

func scanInterval(row rowScanner) (*Interval, error) {
	var (
		timestampSeconds float64
		periodSeconds    float64
		cacheBytes       float64
		originBytes      float64
		parentBytes      float64
		peerBytes        float64
		cacheReplies     float64
		originReplies    float64
		parentReplies    float64
		peerReplies      float64
		childRequests    float64
		clientRequests   float64
		peerRequests     float64
		bytesImported    float64
		imports          float64
		bytesDropped     float64
		bytesPurged      float64
		rejectedNoSpace  float64
		pressureRatio    float64
	)
	if err := row.Scan(
		&timestampSeconds,
		&periodSeconds,
		&cacheBytes,
		&originBytes,
		&parentBytes,
		&peerBytes,
		&cacheReplies,
		&originReplies,
		&parentReplies,
		&peerReplies,
		&childRequests,
		&clientRequests,
		&peerRequests,
		&bytesImported,
		&imports,
		&bytesDropped,
		&bytesPurged,
		&rejectedNoSpace,
		&pressureRatio,
	); err != nil {
		return nil, err
	}

	seconds, fraction := math.Modf(timestampSeconds)
	return &Interval{
		Timestamp:     time.Unix(int64(seconds), int64(fraction*float64(time.Second))).UTC(),
		PeriodSeconds: periodSeconds,
		BytesServed: []Observation{
			{Source: "cache", Value: cacheBytes},
			{Source: "origin", Value: originBytes},
			{Source: "parent", Value: parentBytes},
			{Source: "peer", Value: peerBytes},
		},
		Replies: []Observation{
			{Source: "cache", Value: cacheReplies},
			{Source: "origin", Value: originReplies},
			{Source: "parent", Value: parentReplies},
			{Source: "peer", Value: peerReplies},
		},
		Requests: []Observation{
			{Source: "child", Value: childRequests},
			{Source: "client", Value: clientRequests},
			{Source: "peer", Value: peerRequests},
		},
		BytesImported:           bytesImported,
		Imports:                 imports,
		BytesDropped:            bytesDropped,
		BytesPurged:             bytesPurged,
		RequestsRejectedNoSpace: rejectedNoSpace,
		CachePressureRatio:      pressureRatio,
	}, nil
}
