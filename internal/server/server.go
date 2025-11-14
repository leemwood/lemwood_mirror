package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"lemwood_mirror/internal/storage"
)

type State struct {
	BasePath string
	// cached status: map[launcher]map[version]infoPath
	mu     sync.RWMutex
	index  map[string]map[string]string
	latest map[string]string
}

func NewState(base string) *State {
	return &State{BasePath: base, index: make(map[string]map[string]string), latest: make(map[string]string)}
}

func (s *State) UpdateIndex(launcher string, version string, infoPath string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.index[launcher] == nil {
		s.index[launcher] = make(map[string]string)
	}
	s.index[launcher][version] = infoPath
	s.latest[launcher] = s.pickLatest(s.index[launcher])
}

func (s *State) RemoveVersion(launcher string, version string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.index[launcher] == nil {
		return
	}
	delete(s.index[launcher], version)
	s.latest[launcher] = s.pickLatest(s.index[launcher])
}

func (s *State) Routes(mux *http.ServeMux) {
	// Static UI
	staticDir := filepath.Join("web", "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})
	// Downloads
	mux.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(s.BasePath))))
	// API endpoints
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/status/", s.handleLauncherStatus)
	mux.HandleFunc("/api/files", s.handleFiles)
	mux.HandleFunc("/api/latest", s.handleLatestAll)
	mux.HandleFunc("/api/latest/", s.handleLatestLauncher)
}

// RoutesWithScan adds /api/scan endpoint to trigger a scan callback.
func (s *State) RoutesWithScan(mux *http.ServeMux, scan func()) {
	s.Routes(mux)
	mux.HandleFunc("/api/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		go func() {
			// run scan asynchronously to avoid blocking request
			defer func() { recover() }()
			scan()
		}()
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "scan started"})
	})
}

func (s *State) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	indexCopy := make(map[string]map[string]string)
	for k, v := range s.index {
		inner := make(map[string]string)
		for vk, vv := range v {
			inner[vk] = vv
		}
		indexCopy[k] = inner
	}
	s.mu.RUnlock()
	// Build response reading info.json content
	resp := make(map[string][]map[string]any)
	for launcher, versions := range indexCopy {
		for version, infoPath := range versions {
			v, err := storage.ReadInfoJSON(infoPath)
			if err != nil {
				log.Printf("read info.json failed for %s %s: %v", launcher, version, err)
				continue
			}
			relPath, err := filepath.Rel(s.BasePath, filepath.Dir(infoPath))
			if err != nil {
				log.Printf("could not get relative path for %s: %v", infoPath, err)
			} else {
				v["download_path"] = filepath.ToSlash(filepath.Join("download", relPath))
			}
			resp[launcher] = append(resp[launcher], v)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *State) handleLauncherStatus(w http.ResponseWriter, r *http.Request) {
	launcherID := filepath.Base(r.URL.Path)
	s.mu.RLock()
	versions, ok := s.index[launcherID]
	s.mu.RUnlock()

	if !ok {
		http.NotFound(w, r)
		return
	}

	resp := make([]map[string]any, 0, len(versions))
	for version, infoPath := range versions {
		v, err := storage.ReadInfoJSON(infoPath)
		if err != nil {
			log.Printf("read info.json failed for %s %s: %v", launcherID, version, err)
			continue
		}
		relPath, err := filepath.Rel(s.BasePath, filepath.Dir(infoPath))
		if err != nil {
			log.Printf("could not get relative path for %s: %v", infoPath, err)
		} else {
			v["download_path"] = filepath.ToSlash(filepath.Join("download", relPath))
		}
		resp = append(resp, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *State) handleFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}
	n, err := storage.ListTree(s.BasePath, path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

func (s *State) handleLatestAll(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	indexCopy := make(map[string]map[string]string)
	latestCopy := make(map[string]string)
	for k, v := range s.index {
		inner := make(map[string]string)
		for vk, vv := range v {
			inner[vk] = vv
		}
		indexCopy[k] = inner
	}
	for k, v := range s.latest {
		latestCopy[k] = v
	}
	s.mu.RUnlock()
	resp := make(map[string]map[string]any)
	for launcher, lver := range latestCopy {
		if lver == "" {
			continue
		}
		infoPath := indexCopy[launcher][lver]
		v, err := storage.ReadInfoJSON(infoPath)
		if err != nil {
			continue
		}
		relPath, err := filepath.Rel(s.BasePath, filepath.Dir(infoPath))
		if err == nil {
			v["download_path"] = filepath.ToSlash(filepath.Join("download", relPath))
		}
		v["latest"] = true
		resp[launcher] = v
	}
	pairs := make([]string, 0, len(latestCopy))
	for k, v := range latestCopy {
		if v != "" {
			pairs = append(pairs, k+"="+v)
		}
	}
	sort.Strings(pairs)
	w.Header().Set("Content-Type", "application/json")
	if len(pairs) > 0 {
		w.Header().Set("X-Latest-Versions", strings.Join(pairs, ","))
	}
	json.NewEncoder(w).Encode(resp)
}

func (s *State) handleLatestLauncher(w http.ResponseWriter, r *http.Request) {
	launcherID := filepath.Base(r.URL.Path)
	s.mu.RLock()
	lver := s.latest[launcherID]
	infoPath := ""
	if idx := s.index[launcherID]; idx != nil {
		infoPath = idx[lver]
	}
	s.mu.RUnlock()
	if lver == "" || infoPath == "" {
		http.NotFound(w, r)
		return
	}
	v, err := storage.ReadInfoJSON(infoPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
		return
	}
	relPath, err := filepath.Rel(s.BasePath, filepath.Dir(infoPath))
	if err == nil {
		v["download_path"] = filepath.ToSlash(filepath.Join("download", relPath))
	}
	v["latest"] = true
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Latest-Version", lver)
	json.NewEncoder(w).Encode(v)
}

func (s *State) pickLatest(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	stable := make([]string, 0, len(m))
	all := make([]string, 0, len(m))
	for v := range m {
		all = append(all, v)
		if isStableVersion(v) {
			stable = append(stable, v)
		}
	}
	if len(stable) == 0 {
		return maxVersion(all)
	}
	return maxVersion(stable)
}

func maxVersion(list []string) string {
	if len(list) == 0 {
		return ""
	}
	sort.Slice(list, func(i, j int) bool { return compareVersions(list[i], list[j]) < 0 })
	return list[len(list)-1]
}

func isStableVersion(v string) bool {
	s := strings.ToLower(v)
	if strings.Contains(s, "alpha") {
		return false
	}
	if strings.Contains(s, "beta") {
		return false
	}
	if strings.Contains(s, "rc") {
		return false
	}
	if strings.Contains(s, "snapshot") {
		return false
	}
	if strings.Contains(s, "pre") {
		return false
	}
	if strings.Contains(s, "dev") {
		return false
	}
	return true
}

func compareVersions(a, b string) int {
	sa := normalizeVersion(a)
	sb := normalizeVersion(b)
	aa := strings.Split(sa, ".")
	bb := strings.Split(sb, ".")
	n := len(aa)
	if len(bb) > n {
		n = len(bb)
	}
	for i := 0; i < n; i++ {
		var ai, bi int64
		if i < len(aa) {
			ai = parseIntSafe(aa[i])
		}
		if i < len(bb) {
			bi = parseIntSafe(bb[i])
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}
	return 0
}

func normalizeVersion(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "v")
	s = strings.ReplaceAll(s, "_", ".")
	s = strings.ReplaceAll(s, "-", ".")
	s = strings.Trim(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

func parseIntSafe(s string) int64 {
	var n int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int64(c-'0')
	}
	return n
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func StartHTTP(addr string, s *State) error {
	mux := http.NewServeMux()
	s.Routes(mux)
	log.Printf("HTTP server listening on %s", addr)
	return http.ListenAndServe(addr, corsMiddleware(mux))
}

func StartHTTPWithScan(addr string, s *State, scan func()) error {
	mux := http.NewServeMux()
	s.RoutesWithScan(mux, scan)
	log.Printf("HTTP server listening on %s", addr)
	return http.ListenAndServe(addr, corsMiddleware(mux))
}

// Ensure directories exist on startup
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
