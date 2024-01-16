// Package frontend contains compiled frontend assets and a registration
// for the HTTP listeners
package frontend

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	frontendServer struct {
		router *mux.Router
		log    *logrus.Logger
	}
)

//go:embed assets/*
var assets embed.FS

// RegisterHandler takes a Router and registers the frontend onto that
// router
func RegisterHandler(router *mux.Router, logger *logrus.Logger) {
	srv := frontendServer{router, logger}

	router.
		PathPrefix("/assets").
		Handler(http.StripPrefix("/assets", http.HandlerFunc(srv.handleAsset))).
		Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(srv.handleIndex)
}

func (f frontendServer) handleAsset(w http.ResponseWriter, r *http.Request) {
	asset, err := assets.Open(path.Join("assets", r.URL.Path))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "that's not the file you're looking for", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("opening asset: %s", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := asset.Close(); err != nil {
			f.log.WithError(err).Error("closing assets file (leaked fd)")
		}
	}()

	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(r.URL.Path)))

	if _, err = io.Copy(w, asset); err != nil {
		f.log.WithError(err).Debug("copying index to browser")
	}
}

func (f frontendServer) handleIndex(w http.ResponseWriter, _ *http.Request) {
	index, err := assets.Open("assets/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("opening index: %s", err), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := index.Close(); err != nil {
			f.log.WithError(err).Error("closing assets file (leaked fd)")
		}
	}()

	if _, err = io.Copy(w, index); err != nil {
		f.log.WithError(err).Debug("copying index to browser")
	}
}
