package assetcache

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func TestMetricsSourceReadsLatestInterval(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Metrics test.db")
	database, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	createMetricsTable(t, database)
	if _, err := database.Exec(`
		INSERT INTO ZMETRIC (
			ZCREATIONDATE, ZPERIOD,
			ZBYTESFROMCACHETOCLIENT, ZBYTESFROMORIGINTOCLIENT,
			ZREPLIESFROMCACHETOCLIENT, ZREPLIESFROMORIGINTOCLIENT,
			ZREQUESTSFROMCLIENT, ZBYTESIMPORTEDBYHTTP, ZIMPORTSBYHTTP,
			ZBYTESDROPPED, ZBYTESPURGEDTOTAL, ZBYTESPURGEDYOUNGERTHAN7DAYS,
			ZREQUESTSREJECTEDFORNOSPACE
		) VALUES
			(100, 60, 1, 2, 3, 4, 5, 6, 7, 8, 9, 1, 0),
			(200.25, 60, 10, 20, 30, 40, 50, 60, 70, 80, 90, 1, 0)
	`); err != nil {
		t.Fatalf("insert fixture rows: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close fixture database: %v", err)
	}

	source := MetricsSource{path: path}
	interval, err := source.Read(context.Background())
	if err != nil {
		t.Fatalf("read metrics source: %v", err)
	}
	if interval == nil {
		t.Fatal("expected an interval")
	}
	wantTime := time.Unix(coreDataEpoch+200, 250_000_000).UTC()
	if !interval.Timestamp.Equal(wantTime) {
		t.Fatalf("timestamp = %v, want %v", interval.Timestamp, wantTime)
	}
	if got, want := observationValue(t, interval.BytesServed, "cache"), float64(10); got != want {
		t.Fatalf("cache bytes = %v, want %v", got, want)
	}
	if got, want := observationValue(t, interval.BytesServed, "origin"), float64(20); got != want {
		t.Fatalf("origin bytes = %v, want %v", got, want)
	}
	if got, want := observationValue(t, interval.Replies, "cache"), float64(30); got != want {
		t.Fatalf("cache replies = %v, want %v", got, want)
	}
	if got, want := observationValue(t, interval.Requests, "client"), float64(50); got != want {
		t.Fatalf("client requests = %v, want %v", got, want)
	}
	if got, want := interval.CachePressureRatio, 0.6; got != want {
		t.Fatalf("pressure ratio = %v, want %v", got, want)
	}
}

func TestMetricsSourceRejectsIncompatibleSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Metrics.db")
	database, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	if _, err := database.Exec(`CREATE TABLE ZMETRIC (ZCREATIONDATE REAL)`); err != nil {
		t.Fatalf("create incompatible table: %v", err)
	}
	if err := database.Close(); err != nil {
		t.Fatalf("close fixture database: %v", err)
	}

	if _, err := (MetricsSource{path: path}).Read(context.Background()); err == nil {
		t.Fatal("expected incompatible schema to fail")
	}
}

func createMetricsTable(t *testing.T, database *sql.DB) {
	t.Helper()
	if _, err := database.Exec(`
		CREATE TABLE ZMETRIC (
			ZCREATIONDATE REAL NOT NULL,
			ZPERIOD INTEGER,
			ZBYTESFROMCACHETOCHILD INTEGER,
			ZBYTESFROMCACHETOCLIENT INTEGER,
			ZBYTESFROMCACHETOPEER INTEGER,
			ZBYTESFROMORIGINTOCHILD INTEGER,
			ZBYTESFROMORIGINTOCLIENT INTEGER,
			ZBYTESFROMORIGINTOPEER INTEGER,
			ZBYTESFROMPARENTTOCHILD INTEGER,
			ZBYTESFROMPARENTTOCLIENT INTEGER,
			ZBYTESFROMPARENTTOPEER INTEGER,
			ZBYTESFROMPEERTOCHILD INTEGER,
			ZBYTESFROMPEERTOCLIENT INTEGER,
			ZREPLIESFROMCACHETOCHILD INTEGER,
			ZREPLIESFROMCACHETOCLIENT INTEGER,
			ZREPLIESFROMCACHETOPEER INTEGER,
			ZREPLIESFROMORIGINTOCHILD INTEGER,
			ZREPLIESFROMORIGINTOCLIENT INTEGER,
			ZREPLIESFROMORIGINTOPEER INTEGER,
			ZREPLIESFROMPARENTTOCHILD INTEGER,
			ZREPLIESFROMPARENTTOCLIENT INTEGER,
			ZREPLIESFROMPARENTTOPEER INTEGER,
			ZREPLIESFROMPEERTOCHILD INTEGER,
			ZREPLIESFROMPEERTOCLIENT INTEGER,
			ZREQUESTSFROMCHILD INTEGER,
			ZREQUESTSFROMCLIENT INTEGER,
			ZREQUESTSFROMPEER INTEGER,
			ZBYTESIMPORTEDBYHTTP INTEGER,
			ZBYTESIMPORTEDBYXPC INTEGER,
			ZIMPORTSBYHTTP INTEGER,
			ZIMPORTSBYXPC INTEGER,
			ZBYTESDROPPED INTEGER,
			ZBYTESPURGEDTOTAL INTEGER,
			ZBYTESPURGEDYOUNGERTHAN1DAY INTEGER,
			ZBYTESPURGEDYOUNGERTHAN7DAYS INTEGER,
			ZBYTESPURGEDYOUNGERTHAN30DAYS INTEGER,
			ZREQUESTSREJECTEDFORNOSPACE INTEGER
		)
	`); err != nil {
		t.Fatalf("create fixture table: %v", err)
	}
}

func observationValue(t *testing.T, observations []Observation, source string) float64 {
	t.Helper()
	for _, observation := range observations {
		if observation.Source == source {
			return observation.Value
		}
	}
	t.Fatalf("observation %q not found", source)
	return 0
}
