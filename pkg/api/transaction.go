package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"git.luzifer.io/luzifer/accounting/pkg/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	jsonpatch "gopkg.in/evanphx/json-patch.v5"
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

	http.Redirect(w, r, u.String(), http.StatusFound)
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

func (a apiServer) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	var (
		since time.Time
		until = time.Now()
	)
	if v, err := time.Parse(time.RFC3339, r.URL.Query().Get("since")); err == nil {
		since = v
	}
	if v, err := time.Parse(time.RFC3339, r.URL.Query().Get("until")); err == nil {
		until = v
	}

	txs, err := a.dbc.ListTransactions(since, until)
	if err != nil {
		a.errorResponse(w, err, "getting transactions", http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, http.StatusOK, txs)
}

func (a apiServer) handleListTransactionsByAccount(w http.ResponseWriter, r *http.Request) {
	accid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	var (
		since time.Time
		until = time.Now()
	)
	if v, err := time.Parse(time.RFC3339, r.URL.Query().Get("since")); err == nil {
		since = v
	}
	if v, err := time.Parse(time.RFC3339, r.URL.Query().Get("until")); err == nil {
		until = v
	}

	txs, err := a.dbc.ListTransactionsByAccount(accid, since, until)
	if err != nil {
		a.errorResponse(w, err, "getting transactions", http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, http.StatusOK, txs)
}

func (a apiServer) handleOverwriteTransaction(w http.ResponseWriter, r *http.Request) {
	var (
		tx   database.Transaction
		txID uuid.UUID
		err  error
	)

	if txID, err = uuid.Parse(mux.Vars(r)["id"]); err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&tx); err != nil {
		a.errorResponse(w, err, "parsing body", http.StatusBadRequest)
		return
	}

	if err = a.dbc.UpdateTransaction(txID, tx); err != nil {
		a.errorResponse(w, err, "updating transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	a.handleTransactionJSONPatch(txID, w, r)
}

func (a apiServer) handleTransactionJSONPatch(txID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		reqBody = new(bytes.Buffer)
	)

	if _, err = io.Copy(reqBody, r.Body); err != nil {
		a.errorResponse(w, err, "reading request body", http.StatusBadRequest)
		return
	}

	if reqBody.Len() < 2 { //nolint:gomnd
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var patch jsonpatch.Patch
	if err = json.NewDecoder(reqBody).Decode(&patch); err != nil {
		a.errorResponse(w, err, "parsing json-patch body", http.StatusBadRequest)
		return
	}

	if len(patch) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	tx, err := a.dbc.GetTransactionByID(txID)
	if err != nil {
		a.errorResponse(w, err, "getting transaction", http.StatusInternalServerError)
		return
	}

	txdoc, err := json.Marshal(tx)
	if err != nil {
		a.errorResponse(w, err, "marshalling transaction", http.StatusInternalServerError)
		return
	}

	if txdoc, err = patch.Apply(txdoc); err != nil {
		a.errorResponse(w, err, "applying patch", http.StatusInternalServerError)
		return
	}

	var updTx database.Transaction
	if err = json.Unmarshal(txdoc, &updTx); err != nil {
		a.errorResponse(w, err, "unmarshalling transaction", http.StatusInternalServerError)
		return
	}

	if err = a.dbc.UpdateTransaction(txID, updTx); err != nil {
		a.errorResponse(w, err, "updating transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
