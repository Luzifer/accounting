package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"git.luzifer.io/luzifer/accounting/pkg/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (a apiServer) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	var payload database.Transaction

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.errorResponse(w, err, "parsing body", http.StatusBadRequest)
		return
	}

	if payload.ID != uuid.Nil {
		a.errorResponse(w, errors.New("transaction id must be unset"), "validating request", http.StatusBadRequest)
		return
	}

	tx, err := a.dbc.CreateTransaction(payload)
	if err != nil {
		a.errorResponse(w, err, "creating transaction", http.StatusInternalServerError)
		return
	}

	u, err := a.router.Get("GetTransactionByID").URL("id", tx.ID.String())
	if err != nil {
		a.errorResponse(w, err, "getting redirect url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (a apiServer) handleDeleteTransaction(w http.ResponseWriter, r *http.Request) {
	txid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	if err = a.dbc.DeleteTransaction(txid); err != nil {
		a.errorResponse(w, err, "deleting transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a apiServer) handleGetTransactionByID(w http.ResponseWriter, r *http.Request) {
	txid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	tx, err := a.dbc.GetTransactionByID(txid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.errorResponse(w, err, "getting transaction", http.StatusNotFound)
			return
		}
		a.errorResponse(w, err, "getting transaction", http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, http.StatusOK, tx)
}

func (a apiServer) handleListTransactionsByAccount(w http.ResponseWriter, r *http.Request) {
	accid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	var since time.Time
	if v, err := time.Parse(time.RFC3339, r.URL.Query().Get("since")); err == nil {
		since = v
	}

	txs, err := a.dbc.ListTransactionsByAccount(accid, since)
	if err != nil {
		a.errorResponse(w, err, "getting transactions", http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, http.StatusOK, txs)
}

func (a apiServer) handleUpdateTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		txID uuid.UUID
		err  error
	)

	if txID, err = uuid.Parse(mux.Vars(r)["id"]); err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Has("cleared") {
		if err = a.dbc.UpdateTransactionCleared(txID, r.URL.Query().Get("cleared") == "true"); err != nil {
			a.errorResponse(w, err, "updating transaction cleared", http.StatusInternalServerError)
			return
		}
	}

	if r.URL.Query().Has("category") {
		cat, err := uuid.Parse(r.URL.Query().Get("category"))
		if err != nil {
			a.errorResponse(w, err, "parsing category id", http.StatusBadRequest)
			return
		}

		if err = a.dbc.UpdateTransactionCategory(txID, cat); err != nil {
			a.errorResponse(w, err, "updating transaction category", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
