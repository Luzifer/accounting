package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"git.luzifer.io/luzifer/accounting/pkg/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (a apiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Name string               `json:"name"`
		Type database.AccountType `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.errorResponse(w, err, "parsing body", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		a.errorResponse(w, errors.New("empty name"), "validating request", http.StatusBadRequest)
		return
	}

	if !payload.Type.IsValid() {
		a.errorResponse(w, errors.New("invalid account type"), "validating request", http.StatusBadRequest)
		return
	}

	acc, err := a.dbc.CreateAccount(payload.Name, payload.Type)
	if err != nil {
		a.errorResponse(w, err, "creating account", http.StatusInternalServerError)
		return
	}

	u, err := a.router.Get("GetAccount").URL("id", acc.ID.String())
	if err != nil {
		a.errorResponse(w, err, "getting redirect url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (a apiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accid, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	acc, err := a.dbc.GetAccount(accid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.errorResponse(w, err, "getting account", http.StatusNotFound)
			return
		}
		a.errorResponse(w, err, "getting account", http.StatusInternalServerError)
		return
	}

	a.jsonResponse(w, http.StatusOK, acc)
}

func (a apiServer) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	var payload any
	if r.URL.Query().Has("with-balances") {
		accs, err := a.dbc.ListAccountBalances()
		if err != nil {
			a.errorResponse(w, err, "getting account balances", http.StatusInternalServerError)
			return
		}
		payload = accs
	} else {
		at := database.AccountType(r.URL.Query().Get("account-type"))
		if at.IsValid() {
			accs, err := a.dbc.ListAccountsByType(at)
			if err != nil {
				a.errorResponse(w, err, "getting accounts", http.StatusInternalServerError)
				return
			}
			payload = accs
		} else {
			accs, err := a.dbc.ListAccounts()
			if err != nil {
				a.errorResponse(w, err, "getting accounts", http.StatusInternalServerError)
				return
			}
			payload = accs
		}
	}

	a.jsonResponse(w, http.StatusOK, payload)
}

func (a apiServer) handleTransferMoney(w http.ResponseWriter, r *http.Request) {
	var (
		amount             float64
		err                error
		from, to, category uuid.UUID
	)

	if from, err = uuid.Parse(mux.Vars(r)["id"]); err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	if to, err = uuid.Parse(mux.Vars(r)["to"]); err != nil {
		a.errorResponse(w, err, "parsing to", http.StatusBadRequest)
		return
	}

	if amount, err = strconv.ParseFloat(mux.Vars(r)["to"], 64); err != nil {
		a.errorResponse(w, err, "parsing amount", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Has("category") {
		if category, err = uuid.Parse(mux.Vars(r)["category"]); err != nil {
			a.errorResponse(w, err, "parsing category", http.StatusBadRequest)
			return
		}
	}

	if category == uuid.Nil {
		if err = a.dbc.TransferMoney(from, to, amount); err != nil {
			a.errorResponse(w, err, "transferring money", http.StatusInternalServerError)
			return
		}
	}

	if err = a.dbc.TransferMoneyWithCategory(from, to, amount, category); err != nil {
		a.errorResponse(w, err, "transferring money", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a apiServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	var (
		acctID uuid.UUID
		err    error
	)

	if acctID, err = uuid.Parse(mux.Vars(r)["id"]); err != nil {
		a.errorResponse(w, err, "parsing id", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Has("name") {
		if err = a.dbc.UpdateAccountName(acctID, r.URL.Query().Get("name")); err != nil {
			a.errorResponse(w, err, "renaming account", http.StatusInternalServerError)
			return
		}
	}

	if r.URL.Query().Has("hidden") {
		if err = a.dbc.UpdateAccountHidden(acctID, r.URL.Query().Get("hidden") == "true"); err != nil {
			a.errorResponse(w, err, "updating account visibility", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}