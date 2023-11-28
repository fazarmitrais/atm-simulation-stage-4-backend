package Controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fazarmitrais/atm-simulation-stage-3/cookies"
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	jwtlib "github.com/fazarmitrais/atm-simulation-stage-3/lib/jwtLib"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	service *service.Service
}

type ResponseFormatter struct {
	IsError bool   `json:"isError"`
	Mesage  string `json:"message"`
}

func New(svc *service.Service) *Controller {
	return &Controller{service: svc}
}

func (re *Controller) Register(e *echo.Echo) {
	auth := e.Group("/api/v1/auth")
	auth.GET("/token-validation", re.TokenValidation, jwtlib.Required())
	auth.POST("/login", re.PINValidation)
	account := e.Group("/api/v1/account", jwtlib.Required())
	account.GET("", re.Accounts)
	account.POST("/import", re.Import)
	account.POST("/create", re.Create)
	withdraw := e.Group("/api/v1/withdraw", jwtlib.Required())
	withdraw.GET("", re.Withdraw)
	withdraw.POST("", re.Withdraw)
	withdraw.GET("/summary", re.WithdrawSummary)
	transfer := e.Group("/api/v1/transfer", jwtlib.Required())
	transfer.GET("", re.Transfer)
	transfer.POST("", re.Transfer)
	transfer.GET("/summary", re.TransferSummary)
}

var response = make(map[string]interface{})

func (re *Controller) TokenValidation(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (re *Controller) Transfer(c echo.Context) error {
	var statusCode = http.StatusOK
	if c.Request().Method == http.MethodPost {
		accNbr := cookies.GetAccountNumberFromCookie(c)
		var trf trxEntity.Transaction
		c.Bind(&trf)
		trf.AccountNumber = accNbr
		log.Println(trf)
		_, resp := re.service.Transfer(c, trf)
		if resp != nil {
			statusCode = resp.Code
			log.Println(resp.Message)
			response["message"] = resp.Message
		}
	}
	return c.Render(statusCode, "transfer.html", response)
}

func (re *Controller) TransferSummary(c echo.Context) error {
	accNbr := cookies.GetAccountNumberFromCookie(c)
	tp := string(trxEntity.TYPE_TRANSFER)
	var statusCode int = http.StatusOK
	trx, err := re.service.GetLastTransaction(c, accNbr, &tp, 1)
	if err != nil {
		statusCode = err.Code
		response["message"] = err.Error()
	} else if len(trx) > 0 {
		acc, err := re.service.GetByAccountNumber(c, accNbr)
		if err != nil {
			statusCode = err.Code
			response["message"] = err.Error()
		}
		response["DestAccNbr"] = trx[0].TransferToAccountNumber
		response["Amount"] = fmt.Sprintf("%0.f", trx[0].Amount)
		response["Balance"] = fmt.Sprintf("%0.f", acc.Balance)
		response["ReferenceNumber"] = trx[0].ReferenceNumber
	}
	return c.Render(statusCode, "transferSummary.html", response)
}

func (re *Controller) Withdraw(c echo.Context) error {
	var statusCode = http.StatusOK
	if c.Request().Method == http.MethodPost {
		accNbr := cookies.GetAccountNumberFromCookie(c)
		type transferAmount struct {
			Amount float64 `json:"amount"`
		}
		amt := transferAmount{}
		c.Bind(&amt)
		_, resp := re.service.Withdraw(c, accNbr, amt.Amount)
		if resp != nil {
			statusCode = resp.Code
			response["message"] = resp.Message
		}
	}
	return c.Render(statusCode, "withdraw.html", response)
}

func (re *Controller) WithdrawSummary(c echo.Context) error {
	accNbr := cookies.GetAccountNumberFromCookie(c)
	tp := string(trxEntity.TYPE_WITHDRAWAL)
	var statusCode int = http.StatusOK
	trx, err := re.service.GetLastTransaction(c, accNbr, &tp, 1)
	if err != nil {
		statusCode = err.Code
		response["message"] = err.Error()
	} else if len(trx) > 0 {
		acc, err := re.service.GetByAccountNumber(c, accNbr)
		if err != nil {
			statusCode = err.Code
			response["message"] = err.Error()
		}
		response["Date"] = trx[0].Date.Format("2006-01-02 15:04")
		response["Withdraw"] = fmt.Sprintf("%0.f", trx[0].Amount)
		response["Balance"] = fmt.Sprintf("%0.f", acc.Balance)
	}
	return c.Render(statusCode, "withdrawSummary.html", response)
}

func (re *Controller) Create(c echo.Context) error {
	var account entity.Account
	err := c.Bind(&account)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	errl := re.service.Insert(c, account)
	if errl != nil {
		return c.JSON(errl.Code, echo.Map{"error": errl.Message})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}

func (re *Controller) Accounts(c echo.Context) error {
	acc, err := re.service.GetAll(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, acc)
}

func (re *Controller) Import(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error parsing form"})
	}
	files := form.File["fileInput"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error retrieving file information"})
	}
	file := files[0]
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error opening file"})
	}
	defer src.Close()

	dst, err := os.Create("uploads/" + file.Filename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error creating file"})
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error copying file"})
	}

	errl := re.service.Import(c, "uploads/"+file.Filename)
	if errl != nil {
		return c.JSON(errl.Code, echo.Map{"error": errl.Message})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Success"})
}

func (re *Controller) PINValidation(c echo.Context) error {
	var acc entity.Account
	errl := c.Bind(&acc)
	if errl != nil {
		log.Println(errl.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{"error": errl.Error()})
	}
	token, err := re.service.PINValidation(c, acc)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(err.Code, err)
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

/*
func (re *Controller) Register(e *echo.Echo) {

	g := e.Group("/api/v1/account")
	g.POST("/validate", re.PINValidation)
	g.POST("/withdraw", re.Withdraw)
	g.POST("/transfer", re.Transfer)
	g.GET("/balance", re.BalanceCheck)
	g.GET("/exit", re.Exit)
}

func (re *Controller) BalanceCheck(c echo.Context) error {
	acct, resp := re.service.BalanceCheck(c, cok.Values["acctNbr"].(string))
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acct)
	return
}

func (re *Controller) Exit(c echo.Context) error {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	cookieStore.Values["authenticated"] = false
	cookieStore.Values["acctNbr"] = nil
	cookieStore.Save(r, w)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(echo.NewHTTPError(http.StatusOK, "Logout success", false))
}

func (re *Controller) Transfer(c echo.Context) error {
	cookieStore, err := re.cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	var transfer entity.Transfer
	err = json.Unmarshal(b, &transfer)
	if err != nil {
		echo.NewHTTPError(http.StatusBadRequest,
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
