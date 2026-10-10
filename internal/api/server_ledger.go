package api

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const ledgerPrefix = "/v8/management/observability/ledger/"
const ledgerBodyLimit = 2 << 20

// The opt-in destination must be a literal loopback HTTP origin. It cannot use
// DNS, proxy environment variables, redirects, or a caller-provided destination.
func ledgerBridgeFromEnv() gin.HandlerFunc {
	return newLedgerBridge(os.Getenv("CLIPROXY_LEDGER_READER_URL"), 8*time.Second)
}

func newLedgerBridge(origin string, timeout time.Duration) gin.HandlerFunc {
	u, err := url.Parse(origin)
	valid := err == nil && u.Scheme == "http" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && (u.Path == "" || u.Path == "/")
	if valid {
		ip := net.ParseIP(u.Hostname())
		port, portErr := strconv.Atoi(u.Port())
		valid = ip != nil && ip.IsLoopback() && portErr == nil && port > 0 && port <= 65535
	}
	transport := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: timeout}).DialContext, MaxIdleConnsPerHost: 2, IdleConnTimeout: 90 * time.Second}
	client := &http.Client{Transport: transport, Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		if !valid {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "ledger_unavailable"})
			return
		}
		path := c.Request.URL.Path
		if c.Request.Method != http.MethodGet || (path != ledgerPrefix+"status" && path != ledgerPrefix+"summary" && path != ledgerPrefix+"events") {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		target := *u
		target.Path, target.RawPath, target.RawQuery = path, "", c.Request.URL.RawQuery
		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, target.String(), nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "ledger_unavailable"})
			return
		}
		// Middleware accepts either form; normalize it for the reader. Forward no
		// cookies, arbitrary headers or request bodies and never log authorization.
		auth := c.GetHeader("Authorization")
		if auth == "" && c.GetHeader("X-Management-Key") != "" {
			auth = "Bearer " + c.GetHeader("X-Management-Key")
		}
		req.Header.Set("Authorization", auth)
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "ledger_unavailable"})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			status, message := http.StatusServiceUnavailable, "ledger_unavailable"
			if resp.StatusCode == http.StatusBadRequest {
				status, message = http.StatusBadRequest, "invalid_or_oversized_query"
			}
			c.AbortWithStatusJSON(status, gin.H{"error": message})
			return
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, ledgerBodyLimit+1))
		if err != nil || len(body) > ledgerBodyLimit || !json.Valid(body) {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "ledger_unavailable"})
			return
		}
		c.Data(http.StatusOK, "application/json", body)
	}
}
