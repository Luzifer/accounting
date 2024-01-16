package main

import (
	"net/http"
	"os"
	"time"

	"git.luzifer.io/luzifer/accounting/pkg/api"
	"git.luzifer.io/luzifer/accounting/pkg/database"
	"git.luzifer.io/luzifer/accounting/pkg/frontend"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		DatabaseConnection string `flag:"database-connection" default:"file::memory:?cache=shared" description:"Connection string for the selected database type"`
		DatabaseType       string `flag:"database-type" default:"sqlite" description:"Type of the database to connect to (postgres, sqlite)"`
		Listen             string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel           string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		VersionAndExit     bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		logrus.WithField("version", version).Info("accounting")
		os.Exit(0)
	}

	dbc, err := database.New(cfg.DatabaseType, cfg.DatabaseConnection)
	if err != nil {
		logrus.WithError(err).Fatal("connecting to database")
	}

	router := mux.NewRouter()
	api.RegisterHandler(router.PathPrefix("/api").Subrouter(), dbc, logrus.StandardLogger())
	frontend.RegisterHandler(router, logrus.StandardLogger())

	var hdl http.Handler = router
	hdl = httpHelper.GzipHandler(hdl)
	hdl = httpHelper.NewHTTPLogHandlerWithLogger(hdl, logrus.StandardLogger())

	server := &http.Server{
		Addr:              cfg.Listen,
		Handler:           hdl,
		ReadHeaderTimeout: time.Second,
	}

	logrus.
		WithField("version", version).
		WithField("addr", cfg.Listen).
		Info("accounting starting")

	if err = server.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("running HTTP server")
	}
}
