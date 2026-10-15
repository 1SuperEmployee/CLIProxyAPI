package helps

import (
	"bufio"
	"context"
	"crypto/rand"
	gotls "crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"testing/synctest"
	"time"
)

// The Chrome transport uses system roots. Isolate the synthetic fallback root
// in a child test process so no other test inherits a changed trust store.
func TestUtlsRoundTripperHandshakeTimeoutPreservesSuccessfulStream(t *testing.T) {
	const childFlag = "SE_SYNTHETIC_CHROME_TLS_CHILD"
	if os.Getenv(childFlag) != "1" {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestUtlsRoundTripperHandshakeTimeoutPreservesSuccessfulStream$", "-test.count=1")
		cmd.Env = append(os.Environ(), childFlag+"=1", "GODEBUG="+os.Getenv("GODEBUG")+",x509usefallbackroots=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("isolated Chrome TLS test: %v\n%s", err, output)
		}
		return
	}
	synctest.Test(t, func(t *testing.T) {
		certificate := newResumptionTestCertificateAt(t, time.Now())
		template := *certificate.Leaf
		template.DNSNames = []string{"chatgpt.com"}
		template.Subject.CommonName = "chatgpt.com"
		der, err := x509.CreateCertificate(rand.Reader, &template, &template, certificate.Leaf.PublicKey, certificate.PrivateKey)
		if err != nil {
			t.Fatal(err)
		}
		certificate.Certificate = [][]byte{der}
		certificate.Leaf, err = x509.ParseCertificate(der)
		if err != nil {
			t.Fatal(err)
		}
		roots := x509.NewCertPool()
		roots.AddCert(certificate.Leaf)
		x509.SetFallbackRoots(roots)
		clientConn, serverConn := net.Pipe()
		defer clientConn.Close()
		defer serverConn.Close()
		recorded := &deadlineRecordingConn{Conn: clientConn, deadlines: make(chan time.Time, 16)}
		transport := &utlsRoundTripper{dialer: contextDialerFunc(func(context.Context, string, string) (net.Conn, error) { return recorded, nil })}
		server := gotls.Server(serverConn, &gotls.Config{Certificates: []gotls.Certificate{certificate}, MinVersion: gotls.VersionTLS13, NextProtos: []string{"http/1.1"}})
		serverDone := make(chan error, 1)
		go func() {
			defer serverConn.Close()
			if err := server.Handshake(); err != nil {
				serverDone <- err
				return
			}
			if _, err := http.ReadRequest(bufio.NewReader(server)); err != nil {
				serverDone <- err
				return
			}
			if _, err := io.WriteString(server, "HTTP/1.1 200 OK\r\nContent-Type: text/event-stream\r\nContent-Length: 12\r\nConnection: close\r\n\r\n"); err != nil {
				serverDone <- err
				return
			}
			// Fake time advances after response headers, beyond the handshake budget.
			time.Sleep(time.Second)
			_, err := io.WriteString(server, "data: done\n\n")
			serverDone <- err
		}()
		const bound = 100 * time.Millisecond
		ctx := withTLSHandshakeTimeout(context.Background(), bound)
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://chatgpt.com/backend-api/codex/responses", nil)
		if err != nil {
			t.Fatal(err)
		}
		started := time.Now()
		response, err := transport.RoundTrip(request)
		if err != nil {
			t.Fatal(err)
		}
		// Check while the response body is still open. Transport cleanup may
		// legitimately set a close_notify deadline after the stream completes.
		select {
		case deadline := <-recorded.deadlines:
			t.Fatalf("residual connection deadline after handshake: %v", deadline)
		default:
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("stream body beyond handshake bound: %v", err)
		}
		if string(body) != "data: done\n\n" || response.StatusCode != http.StatusOK {
			t.Fatalf("unexpected response: %d %q", response.StatusCode, body)
		}
		if time.Since(started) <= bound {
			t.Fatal("response did not extend beyond handshake deadline")
		}
		if err := <-serverDone; err != nil {
			t.Fatal(err)
		}
	})
}
