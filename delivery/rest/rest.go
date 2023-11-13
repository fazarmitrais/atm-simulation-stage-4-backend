package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fazarmitrais/atm-simulation-stage-3/cookie"
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/envLib"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/responseFormatter"
	middleware "github.com/fazarmitrais/atm-simulation-stage-3/middleware"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	service *service.Service
	cookie  *cookie.Cookie
}

type ResponseFormatter struct {
	IsError bool   `json:"isError"`
	Mesage  string `json:"message"`
}

func New(svc *service.Service) *Rest {
	c := cookie.New()
	return &Rest{service: svc, cookie: c}
}

func (re *Rest) Register(m *mux.Router) {
	m = m.PathPrefix("/api/v1/account").Subrouter()
	m.HandleFunc("/validate", re.PINValidation).Methods(http.MethodPost)
	m.HandleFunc("/withdraw", middleware.Chain(re.Withdraw, middleware.Required(re.cookie))).Methods(http.MethodPost)
	m.HandleFunc("/transfer", middleware.Chain(re.Transfer, middleware.Required(re.cookie))).Methods(http.MethodPost)
	m.HandleFunc("/balance", middleware.Chain(re.BalanceCheck, middleware.Required(re.cookie))).Methods(http.MethodGet)
	m.HandleFunc("/exit", re.Exit).Methods(http.MethodGet)
}

func (re *Rest) BalanceCheck(w http.ResponseWriter, r *http.Request) {
	cok, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		responseFormatter.New(http.StatusBadRequest, "Error getting data from cookies", true)
		return
	}
	acct, resp := re.service.BalanceCheck(r.Context(), cok.Values["acctNbr"].(string))
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acct)
	return
}

func (re *Rest) Exit(w http.ResponseWriter, r *http.Request) {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	cookieStore.Values["authenticated"] = false
	cookieStore.Values["acctNbr"] = nil
	cookieStore.Save(r, w)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(responseFormatter.New(http.StatusOK, "Logout success", false))
}

func (re *Rest) PINValidation(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	var acc entity.Account
	err = json.Unmarshal(b, &acc)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	errl := re.service.PINValidation(r.Context(), acc)
	if errl != nil {
		errl.ReturnAsJson(w)
		return
	}
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		responseFormatter.New(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	cookieStore.Values["authenticated"] = true
	cookieStore.Values["acctNbr"] = acc.AccountNumber
	cookieStore.Save(r, w)
	errl.ReturnAsJson(w)
}

func (re *Rest) Withdraw(w http.ResponseWriter, r *http.Request) {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		responseFormatter.New(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	type transferAmount struct {
		Amount float64 `json:"amount"`
	}
	amt := transferAmount{}
	err = json.Unmarshal(b, &amt)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	acc, resp := re.service.Withdraw(r.Context(), fmt.Sprintf("%v", cookieStore.Values["acctNbr"]), amt.Amount)
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (re *Rest) Transfer(w http.ResponseWriter, r *http.Request) {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		responseFormatter.New(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	var transfer entity.Transfer
	err = json.Unmarshal(b, &transfer)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	transfer.FromAccountNumber = cookieStore.Values["acctNbr"].(string)
	acc, resp := re.service.Transfer(r.Context(), transfer)
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acc)
}
