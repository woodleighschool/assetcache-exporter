package assetcache

import (
	"slices"
	"testing"
	"time"
)

func TestParseStatus(t *testing.T) {
	data := []byte(`{
		"name": "status",
		"result": {
			"Activated": true,
			"Active": true,
			"Port": 54321,
			"PrivateAddresses": ["192.0.2.10", "2001:db8::10"],
			"CacheUsed": 40520685527,
			"CacheLimit": 199000000000,
			"CacheFree": 141430352896,
			"PersonalCacheUsed": 123,
			"PersonalCacheLimit": 456,
			"PersonalCacheFree": 333,
			"MaxCachePressureLast1Hour": 40,
			"CacheStatus": "OK",
			"StartupStatus": "OK",
			"RegistrationStatus": 1,
			"Parents": [{}, {}],
			"Peers": [{}],
			"TotalBytesAreSince": "2026-08-07 00:54:51 +0000",
			"TotalBytesDropped": 10,
			"TotalBytesImported": 20,
			"TotalBytesReturnedToChildren": 30,
			"TotalBytesReturnedToClients": 40,
			"TotalBytesReturnedToPeers": 50,
			"TotalBytesStoredFromOrigin": 60,
			"TotalBytesStoredFromParents": 70,
			"TotalBytesStoredFromPeers": 80
		}
	}`)

	status, err := parseStatus(data)
	if err != nil {
		t.Fatalf("parse status: %v", err)
	}
	if !status.Activated || !status.Active {
		t.Fatalf("expected active and activated status, got %+v", status)
	}
	if status.Port == nil || *status.Port != 54321 {
		t.Fatalf("port = %v, want 54321", status.Port)
	}
	if !slices.Equal(status.PrivateAddresses, []string{"192.0.2.10", "2001:db8::10"}) {
		t.Fatalf("private addresses = %v", status.PrivateAddresses)
	}
	if got, want := status.MaxCachePressureRatio, 0.4; got != want {
		t.Fatalf("pressure ratio = %v, want %v", got, want)
	}
	if got, want := status.ParentCount, 2; got != want {
		t.Fatalf("parent count = %d, want %d", got, want)
	}
	if got, want := status.PeerCount, 1; got != want {
		t.Fatalf("peer count = %d, want %d", got, want)
	}
	wantStart := time.Date(2026, time.August, 7, 0, 54, 51, 0, time.UTC)
	if !status.TotalsStartTime.Equal(wantStart) {
		t.Fatalf("totals start time = %v, want %v", status.TotalsStartTime, wantStart)
	}
	if got, want := status.BytesStoredFromPeers, float64(80); got != want {
		t.Fatalf("bytes stored from peers = %v, want %v", got, want)
	}
}

func TestParseStatusWithoutEndpoint(t *testing.T) {
	for _, tt := range []struct {
		data        string
		portPresent bool
	}{
		{`{"result":{"Activated":false,"Active":false}}`, false},
		{`{"result":{"Activated":false,"Active":false,"Port":0,"PrivateAddresses":[]}}`, true},
	} {
		status, err := parseStatus([]byte(tt.data))
		if err != nil {
			t.Fatalf("parse status without endpoint: %v", err)
		}
		if (status.Port != nil) != tt.portPresent {
			t.Fatalf("port present = %v, want %v", status.Port != nil, tt.portPresent)
		}
		if status.Port != nil && *status.Port != 0 {
			t.Fatalf("port = %v, want absent or zero", *status.Port)
		}
		if len(status.PrivateAddresses) != 0 {
			t.Fatalf("unexpected addresses: %v", status.PrivateAddresses)
		}
	}
}

func TestParseStatusRejectsMissingResult(t *testing.T) {
	if _, err := parseStatus([]byte(`{"name":"status"}`)); err == nil {
		t.Fatal("expected missing result to fail")
	}
}

func TestParseStatusRejectsEmptyResult(t *testing.T) {
	if _, err := parseStatus([]byte(`{"result":{}}`)); err == nil {
		t.Fatal("expected empty result to fail")
	}
}

func TestParseStatusRejectsInvalidTotalsTimestamp(t *testing.T) {
	if _, err := parseStatus([]byte(`{"result":{"Activated":true,"Active":true,"TotalBytesAreSince":"not-a-time"}}`)); err == nil {
		t.Fatal("expected invalid totals timestamp to fail")
	}
}
