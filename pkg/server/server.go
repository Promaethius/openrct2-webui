package server

import (
	"context"
	"image/png"
	"net/http"

	"github.com/promaethius/openrct2-webui/pkg/screenshots"
)

type Screenshots func() []screenshots.Screenshot

type Server struct {
	screenshotsFn Screenshots

	server *http.Server
}

func (s *Server) handleScreenshots(w http.ResponseWriter, r *http.Request) {
	screenshots := s.screenshotsFn()

	if len(screenshots) == 0 {
		http.Error(w, "no screenshots are stored", http.StatusNoContent)
		return
	}

	screenshot := screenshots[len(screenshots)-1]
	if err := png.Encode(w, screenshot.Image); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func NewServer(addr string, screenshotsFn Screenshots) *Server {
	mux := http.NewServeMux()
	server := &Server{
		screenshotsFn: screenshotsFn,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}

	mux.HandleFunc("/screenshots", server.handleScreenshots)

	return server
}
