package rest

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
)

type Rest struct {
	service *service.Service
}

type ResponseFormatter struct {
	IsError bool   `json:"isError"`
	Mesage  string `json:"message"`
}

func New(svc *service.Service) *Rest {
	return &Rest{service: svc}
}

/*
func (re *Rest) Register(e *echo.Echo) {

	g := e.Group("/api/v1/account")
	g.POST("/validate", re.PINValidation)
	g.POST("/withdraw", re.Withdraw)
	g.POST("/transfer", re.Transfer)
	g.GET("/balance", re.BalanceCheck)
	g.GET("/exit", re.Exit)
}

func (re *Rest) BalanceCheck(c echo.Context) error {
	acct, resp := re.service.BalanceCheck(c, cok.Values["acctNbr"].(string))
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acct)
	return
}

func (re *Rest) Exit(c echo.Context) error {
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

func (re *Rest) PINValidation(c echo.Context) error {
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

func (re *Rest) Withdraw(c echo.Context) error {
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

func (re *Rest) Transfer(c echo.Context) error {
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
*/
