package middleware

import (
	"net/http"

	"github.com/fazarmitrais/atm-simulation-stage-3/cookie"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/envLib"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/responseFormatter"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Required(cookie *cookie.Cookie) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session, err := cookie.Store.Get(r, envLib.GetEnv("COOKIE_STORE_NAME"))
			if err != nil {
				responseFormatter.New(http.StatusInternalServerError, "Failed to get cookies", true).ReturnAsJson(w)
				return
			}
			if !session.Values["authenticated"].(bool) || session.Values["acctNbr"] == nil {
				responseFormatter.New(http.StatusForbidden, "Please login first", true).ReturnAsJson(w)
				return
			}
			f(w, r)
		}
	}
}

func Chain(f http.HandlerFunc, middleWares ...Middleware) http.HandlerFunc {
	for _, m := range middleWares {
		f = m(f)
	}
	return f
}
