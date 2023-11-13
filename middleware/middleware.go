package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Required() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
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
