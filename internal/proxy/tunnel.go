package proxy

import "net/http"

func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("connect request intercepted", "host", r.Host)
	http.Error(w, "connect tunneling not implemented yet", http.StatusNotImplemented)
}
