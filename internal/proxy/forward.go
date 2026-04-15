package proxy

import "net/http"

func (s *Server) handleHTTP(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("http request intercepted", "method", r.Method, "host", r.Host, "url", r.URL.String())
	http.Error(w, "http forwarding not implemented yet", http.StatusNotImplemented)
}
