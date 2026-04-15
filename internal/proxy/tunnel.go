package proxy

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	decision := s.evaluateRequestDecision(r)

	if shouldBlockDecision(decision) {
		s.logger.Info("connect request blocked by rule", "rule", decision.Rule.Name, "host", requestHost(r))
		http.Error(w, "blocked by proxy rule", http.StatusForbidden)
		return
	}

	if shouldRedirectHTTP(decision) {
		target := strings.TrimSpace(decision.Action.Target)
		s.logger.Warn("redirect action is unsupported for CONNECT", "rule", decision.Rule.Name, "target", target)
		http.Error(w, "redirect action is unsupported for CONNECT", http.StatusBadRequest)
		return
	}

	if r.Host == "" {
		http.Error(w, "missing CONNECT host", http.StatusBadRequest)
		return
	}

	targetAddr := normalizeConnectAddress(r.Host)
	targetConn, err := s.dialer.DialContext(r.Context(), "tcp", targetAddr)
	if err != nil {
		status := http.StatusBadGateway
		message := "failed to connect to target"

		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			status = http.StatusGatewayTimeout
			message = "target connection timed out"
		}

		s.logger.Error("failed to establish CONNECT target", "target", targetAddr, "error", err)
		http.Error(w, message, status)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		targetConn.Close()
		http.Error(w, "proxy does not support connection hijacking", http.StatusInternalServerError)
		return
	}

	clientConn, rw, err := hijacker.Hijack()
	if err != nil {
		targetConn.Close()
		s.logger.Error("failed to hijack client connection", "target", targetAddr, "error", err)
		return
	}

	defer clientConn.Close()
	defer targetConn.Close()

	if _, err := rw.WriteString("HTTP/1.1 200 Connection Established\r\n\r\n"); err != nil {
		s.logger.Warn("failed to write CONNECT response", "target", targetAddr, "error", err)
		return
	}

	if err := rw.Flush(); err != nil {
		s.logger.Warn("failed to flush CONNECT response", "target", targetAddr, "error", err)
		return
	}

	if buffered := rw.Reader.Buffered(); buffered > 0 {
		if _, err := io.CopyN(targetConn, rw, int64(buffered)); err != nil {
			s.logger.Warn("failed to forward buffered CONNECT bytes", "target", targetAddr, "error", err)
			return
		}
	}

	errCh := make(chan error, 2)
	go tunnelCopy(targetConn, clientConn, errCh)
	go tunnelCopy(clientConn, targetConn, errCh)

	firstErr := <-errCh
	secondErr := <-errCh

	if firstErr != nil && !errors.Is(firstErr, io.EOF) {
		s.logger.Warn("CONNECT tunnel closed with downstream error", "target", targetAddr, "error", firstErr)
	}

	if secondErr != nil && !errors.Is(secondErr, io.EOF) {
		s.logger.Warn("CONNECT tunnel closed with upstream error", "target", targetAddr, "error", secondErr)
	}

	s.logger.Info("CONNECT tunnel closed", "target", targetAddr, "duration_ms", time.Since(start).Milliseconds())
}

func normalizeConnectAddress(host string) string {
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host
	}

	return net.JoinHostPort(host, "443")
}

func tunnelCopy(dst, src net.Conn, errCh chan<- error) {
	_, err := io.Copy(dst, src)

	if closeWriter, ok := dst.(interface{ CloseWrite() error }); ok {
		_ = closeWriter.CloseWrite()
	}

	errCh <- err
}
