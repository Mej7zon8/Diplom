package web

import (
	"messenger/web/features"
	"messenger/web/static/certificates"
	"messenger/web/static/router"
	"net/http"
	"os"
	"path/filepath"
)

const certsFolder = "program-data/certificates"

func init() {
	setup()
	registerRoutes()
}

func setup() {
	_ = os.MkdirAll(certsFolder, 0600)
	certificates.GenerateCertificatesIfNotFound(filepath.Join(certsFolder, "cert.crt"), filepath.Join(certsFolder, "cert.key"))
}

func registerRoutes() {
	var staticHandler = router.New()
	staticHandler.RegisterDirectory("html/messenger/dist/messenger/browser")

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseMultipartForm(1 << 30)
		features.Handle(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.HttpHandle(w, r, true)
	})
}

func Run() {
	basicHttpServer{defaultMux: http.DefaultServeMux}.Run()
}
