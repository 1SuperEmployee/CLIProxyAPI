package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/api/handlers/management"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/config"
)

func TestLedgerBridgeAccessAndForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	calls := 0
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != "GET" || r.URL.Path != ledgerPrefix+"events" || r.URL.RawQuery != "model=a%2Fb&limit=1" || r.Header.Get("Authorization") != "Bearer test-password" || r.Header.Get("Cookie") != "" {
			t.Error("unexpected forwarded request")
		}
		w.Header().Set("Set-Cookie", "must-not-forward=yes")
		w.Write([]byte(`{"events":[],"next_cursor":null}`))
	}))
	defer upstream.Close()
	t.Setenv("CLIPROXY_LEDGER_READER_URL", upstream.URL)
	cfg := &config.Config{}
	cfg.RemoteManagement.SecretKey = "synthetic-configured-hash"
	h := management.NewHandler(cfg, "", nil)
	h.SetLocalPassword("test-password")
	s := &Server{cfg: cfg, mgmt: h, engine: gin.New()}
	s.managementRoutesEnabled.Store(true)
	s.registerManagementRoutes()
	for _, tc := range []struct {
		method, path, header string
		want                 int
	}{
		{"GET", ledgerPrefix + "events?model=a%2Fb&limit=1", "", 401},
		{"GET", ledgerPrefix + "events?model=a%2Fb&limit=1", "Authorization", 200},
		{"GET", ledgerPrefix + "events?model=a%2Fb&limit=1", "X-Management-Key", 200},
		{"POST", ledgerPrefix + "events", "Authorization", 404},
		{"GET", ledgerPrefix + "../config", "Authorization", 404},
		{"GET", ledgerPrefix + "other", "Authorization", 404},
	} {
		r := httptest.NewRequest(tc.method, tc.path, nil)
		r.RemoteAddr = "127.0.0.1:1234"
		r.Header.Set("Cookie", "not-forwarded=yes")
		if tc.header == "Authorization" {
			r.Header.Set(tc.header, "Bearer test-password")
		}
		if tc.header == "X-Management-Key" {
			r.Header.Set(tc.header, "test-password")
		}
		w := httptest.NewRecorder()
		s.engine.ServeHTTP(w, r)
		if w.Code != tc.want || w.Header().Get("Set-Cookie") != "" {
			t.Fatalf("%s %s: status=%d", tc.method, tc.path, w.Code)
		}
	}
	if calls != 2 {
		t.Fatalf("upstream calls=%d", calls)
	}
	s.managementRoutesEnabled.Store(false)
	r := httptest.NewRequest("GET", ledgerPrefix+"status", nil)
	r.Header.Set("Authorization", "Bearer test-password")
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, r)
	if w.Code != 404 || calls != 2 {
		t.Fatal("disabled management reached reader")
	}
}

func TestLedgerBridgeDestinationAndFailures(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, origin := range []string{"", "http://example.com:8320", "http://localhost:8320", "http://127.0.0.1", "https://127.0.0.1:8320", "http://user@127.0.0.1:8320", "http://127.0.0.1:8320/path", "http://127.0.0.1:8320?x=1"} {
		e := gin.New()
		e.GET(ledgerPrefix+"status", newLedgerBridge(origin, time.Second))
		w := httptest.NewRecorder()
		e.ServeHTTP(w, httptest.NewRequest("GET", ledgerPrefix+"status", nil))
		if w.Code != 503 {
			t.Fatalf("destination %q accepted", origin)
		}
	}
	for _, tc := range []struct {
		name   string
		status int
		body   string
		delay  time.Duration
		want   int
	}{
		{"query", 400, "private raw failure", 0, 400},
		{"redirect", 302, "private raw failure", 0, 503},
		{"auth failure", 401, "private raw failure", 0, 503},
		{"server failure", 500, "private raw failure", 0, 503},
		{"invalid body", 200, "private raw failure", 0, 503},
		{"large body", 200, strings.Repeat(" ", ledgerBodyLimit+1), 0, 503},
		{"timeout", 200, `{}`, 100 * time.Millisecond, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.delay)
				w.Header().Set("Location", "http://example.com/")
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer u.Close()
			e := gin.New()
			e.GET(ledgerPrefix+"status", newLedgerBridge(u.URL, 30*time.Millisecond))
			w := httptest.NewRecorder()
			e.ServeHTTP(w, httptest.NewRequest("GET", ledgerPrefix+"status", nil))
			if w.Code != tc.want || strings.Contains(w.Body.String(), "private raw failure") || w.Header().Get("Location") != "" {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
		})
	}
	u := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	origin := u.URL
	u.Close()
	e := gin.New()
	e.GET(ledgerPrefix+"status", newLedgerBridge(origin, time.Second))
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest("GET", ledgerPrefix+"status", nil))
	if w.Code != 503 {
		t.Fatal("reader outage not reported")
	}
}
