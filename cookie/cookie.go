package cookie

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/envLib"
	"github.com/gorilla/sessions"
)

type Cookie struct {
	Store *sessions.CookieStore
}

func New() *Cookie {
	key := []byte(envLib.GetEnv("COOKIE_SECRET_KEY"))
	store := sessions.NewCookieStore(key)
	return &Cookie{Store: store}
}
