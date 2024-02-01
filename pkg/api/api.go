// Package api defines an HTTP API for the database interface
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"git.luzifer.io/luzifer/accounting/pkg/database"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	apiServer struct {
		router *mux.Router
		dbc    *database.Client
		log    *logrus.Logger
	}
)

// RegisterHandler takes a (Sub)Router and registers the API onto that
// router
func RegisterHandler(apiRouter *mux.Router, dbc *database.Client, logger *logrus.Logger) {
	as := apiServer{apiRouter, dbc, logger}

	apiRouter.
		HandleFunc("/accounts", as.handleListAccounts).
		Methods(http.MethodGet)
	apiRouter.
		HandleFunc("/accounts", as.handleCreateAccount).
		Methods(http.MethodPost)
	apiRouter.
		HandleFunc("/accounts/{id}", as.handleGetAccount).
		Methods(http.MethodGet).
		Name("GetAccount")
	apiRouter.
		HandleFunc("/accounts/{id}", as.handleUpdateAccount).
		Methods(http.MethodPatch)
	apiRouter.
		HandleFunc("/accounts/{id}/reconcile", as.handleAccountReconcile).
		Methods(http.MethodPut)
	apiRouter.
		HandleFunc("/accounts/{id}/transactions", as.handleListTransactionsByAccount).
		Methods(http.MethodGet)
	apiRouter.
		HandleFunc("/accounts/{id}/transfer/{to}", as.handleTransferMoney).
		Methods(http.MethodPut)

	apiRouter.
		HandleFunc("/transactions", as.handleListTransactions).
		Methods(http.MethodGet)
	apiRouter.
		HandleFunc("/transactions", as.handleCreateTransaction).
		Methods(http.MethodPost)
	apiRouter.
		HandleFunc("/transactions/{id}", as.handleDeleteTransaction).
		Methods(http.MethodDelete)
	apiRouter.
		HandleFunc("/transactions/{id}", as.handleGetTransactionByID).
		Methods(http.MethodGet).
		Name("GetTransactionByID")
	apiRouter.
		HandleFunc("/transactions/{id}", as.handleUpdateTransaction).
		Methods(http.MethodPatch)
	apiRouter.
		HandleFunc("/transactions/{id}", as.handleOverwriteTransaction).
		Methods(http.MethodPut)
}

func (a apiServer) errorResponse(w http.ResponseWriter, err error, desc string, status int) {
	switch status {
	case http.StatusBadRequest:
		a.log.WithError(err).Debug(desc)

	case http.StatusNotFound:
		// No need to log that

	default:
		a.log.WithError(err).Error(desc)
	}

	a.jsonResponse(w, status, struct {
		Error string `json:"error"`
	}{fmt.Sprintf("%s: %s", desc, err)})
}

func (apiServer) jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("encoding response: %s", err), http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	body.WriteTo(w) //nolint:errcheck,gosec,revive
}
