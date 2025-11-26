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
		t.Fatalf("预期 1.2.10 > 1.2.3")
	}
	if compareVersions("141000", "140900") <= 0 {
		t.Fatalf("预期 141000 > 140900")
	}
	if compareVersions("1.4.1.0", "1.4.0.9") <= 0 {
		t.Fatalf("预期 1.4.1.0 > 1.4.0.9")
	}
	if compareVersions("1.0.0", "1.0.0") != 0 {
		t.Fatalf("预期相等")
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
		t.Fatalf("预期最新稳定版 v1.2.3, 实际得到 %s", s.latest["fcl"])
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	s.handleStatus(rr, req)
	if got := rr.Header().Get("X-Latest-Versions"); got != "" {
		t.Fatalf("状态不应包含最新版本头, 实际得到 %q", got)
	}
	var body map[string][]map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("无效的 json: %v", err)
	}
	if len(body["fcl"]) == 0 {
		t.Fatalf("缺少 fcl 条目")
	}
	if _, ok := body["fcl"][0]["latest"]; ok {
		t.Fatalf("状态不应包含最新字段")
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
		t.Fatalf("状态 %d", rr.Code)
	}
	if got := rr.Header().Get("X-Latest-Version"); got != "141000" {
		t.Fatalf("非预期的头: %q", got)
	}

	s.RemoveVersion("zl", "141000")
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/api/latest/zl", nil)
	s.handleLatestLauncher(rr2, req2)
	if got := rr2.Header().Get("X-Latest-Version"); got != "140900" {
		t.Fatalf("回滚头: %q", got)
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
		t.Fatalf("状态 %d", rr.Code)
	}
	if got := rr.Header().Get("X-Latest-Versions"); got == "" {
		t.Fatalf("缺少头")
	}
}

func TestInitFromDiskLoadsAll(t *testing.T) {
	base := t.TempDir()
	os.MkdirAll(filepath.Join(base, "fcl", "v1.2.3"), 0o755)
	p1, _ := writeInfo(filepath.Join(base, "fcl", "v1.2.3"), "v1.2.3")
	_ = p1
	os.MkdirAll(filepath.Join(base, "zl", "141000"), 0o755)
	p2, _ := writeInfo(filepath.Join(base, "zl", "141000"), "141000")
	_ = p2
	s := NewState(base)
	if err := s.InitFromDisk(); err != nil {
		t.Fatalf("从磁盘初始化: %v", err)
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	s.handleStatus(rr, req)
	var body map[string][]map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("无效的 json: %v", err)
	}
	if _, ok := body["fcl"]; !ok {
		t.Fatalf("缺少 fcl")
	}
	if _, ok := body["zl"]; !ok {
		t.Fatalf("缺少 zl")
	}
}
