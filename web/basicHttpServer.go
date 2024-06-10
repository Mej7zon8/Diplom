package web

import (
	"crypto/tls"
	"github.com/k773/utils/io/filteredWriter"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// basicHttpServer is here for directing http requests to different server mux instances.
type basicHttpServer struct {
	defaultMux *http.ServeMux
}

func (m basicHttpServer) Run() {
	var logger = log.New(filteredWriter.NewFilteredWriter(os.Stdout, func(data []byte) bool {
		return !strings.Contains(string(data), "TLS handshake error") && !strings.Contains(string(data), "EOF") && !strings.Contains(string(data), "client requested unsupported application protocols")
	}, '\n', 1<<20), "", 1)
	var servers = []*http.Server{
		{
			ReadTimeout: 30 * time.Second,
			Addr:        ":80",
			Handler:     http.HandlerFunc(m.RedirectTls),
			ErrorLog:    nil,
		},
		{
			ReadTimeout:  30 * time.Second,
			Addr:         ":443",
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			ErrorLog:     logger,
			Handler:      http.HandlerFunc(m.HandleTls),
		},
	}

	for _, s := range servers {
		go func(s *http.Server) {
			println("basicHttpServer/" + s.Addr)
			switch s.Addr {
			case ":80":
				panic(s.ListenAndServe())
			case ":443":
				panic(s.ListenAndServeTLS(filepath.Join(certsFolder, "cert.crt"), filepath.Join(certsFolder, "cert.key")))
			}
		}(s)
	}
}

func (_ basicHttpServer) RedirectTls(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func (m basicHttpServer) HandleTls(w http.ResponseWriter, r *http.Request) {
	// Serve via the default mux
	m.defaultMux.ServeHTTP(w, r)
}
