package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeInfo(dir string, name string) (string, error) {
	p := filepath.Join(dir, "index.json")
	v := map[string]any{"tag_name": name, "name": name}
	b, _ := json.Marshal(v)
	if err := os.WriteFile(p, b, 0o644); err != nil {
		return "", err
	}
	return p, nil
}

func TestCompareVersions(t *testing.T) {
	if compareVersions("v1.2.10", "1.2.3") <= 0 {
		t.Fatalf("expected 1.2.10 > 1.2.3")
	}
	if compareVersions("141000", "140900") <= 0 {
		t.Fatalf("expected 141000 > 140900")
	}
	if compareVersions("1.4.1.0", "1.4.0.9") <= 0 {
		t.Fatalf("expected 1.4.1.0 > 1.4.0.9")
	}
	if compareVersions("1.0.0", "1.0.0") != 0 {
		t.Fatalf("expected equal")
	}
}

func TestPickLatestStable(t *testing.T) {
	base := t.TempDir()
	s := NewState(base)
	v1Dir := filepath.Join(base, "fcl", "v1.2.3")
	os.MkdirAll(v1Dir, 0o755)
	p1, _ := writeInfo(v1Dir, "v1.2.3")
	s.UpdateIndex("fcl", "v1.2.3", p1)

	v2Dir := filepath.Join(base, "fcl", "v1.2.4-rc")
	os.MkdirAll(v2Dir, 0o755)
	p2, _ := writeInfo(v2Dir, "v1.2.4-rc")
	s.UpdateIndex("fcl", "v1.2.4-rc", p2)

	if s.latest["fcl"] != "v1.2.3" {
		t.Fatalf("expected latest stable v1.2.3, got %s", s.latest["fcl"])
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	s.handleStatus(rr, req)
	if got := rr.Header().Get("X-Latest-Versions"); got != "" {
		t.Fatalf("status should not include latest header, got %q", got)
	}
	var body map[string][]map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if len(body["fcl"]) == 0 {
		t.Fatalf("missing fcl entries")
	}
	if _, ok := body["fcl"][0]["latest"]; ok {
		t.Fatalf("status should not include latest field")
	}
}

func TestLatestEndpointAndRollback(t *testing.T) {
	base := t.TempDir()
	s := NewState(base)
	d1 := filepath.Join(base, "zl", "141000")
	os.MkdirAll(d1, 0o755)
	p1, _ := writeInfo(d1, "141000")
	s.UpdateIndex("zl", "141000", p1)
	d2 := filepath.Join(base, "zl", "140900")
	os.MkdirAll(d2, 0o755)
	p2, _ := writeInfo(d2, "140900")
	s.UpdateIndex("zl", "140900", p2)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/latest/zl", nil)
	s.handleLatestLauncher(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if got := rr.Header().Get("X-Latest-Version"); got != "141000" {
		t.Fatalf("unexpected header: %q", got)
	}

	s.RemoveVersion("zl", "141000")
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/latest/zl", nil)
	s.handleLatestLauncher(rr2, req2)
	if got := rr2.Header().Get("X-Latest-Version"); got != "140900" {
		t.Fatalf("rollback header: %q", got)
	}
}

func TestLatestAllEndpoint(t *testing.T) {
	base := t.TempDir()
	s := NewState(base)
	d1 := filepath.Join(base, "fcl", "v1.0.0")
	os.MkdirAll(d1, 0o755)
	p1, _ := writeInfo(d1, "v1.0.0")
	s.UpdateIndex("fcl", "v1.0.0", p1)
	d2 := filepath.Join(base, "zl", "140900")
	os.MkdirAll(d2, 0o755)
	p2, _ := writeInfo(d2, "140900")
	s.UpdateIndex("zl", "140900", p2)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/latest", nil)
	s.handleLatestAll(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if got := rr.Header().Get("X-Latest-Versions"); got == "" {
		t.Fatalf("missing header")
	}
}
