package jwtlib

import (
	"net/http"
	"time"

	"github.com/fazarmitrais/atm-simulation-stage-3/lib/envLib"
	jwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaim struct {
	AccountNumber string `json:"accountNumber"`
	Admin         bool   `json:"admin"`
	jwt.RegisteredClaims
}

func GenerateToken(c echo.Context, accountNumber string) (string, *echo.HTTPError) {
	claims := &JwtCustomClaim{
		AccountNumber: accountNumber,
		Admin:         true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(envLib.GetEnv("JWT_SECRET")))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Error when generating token")
	}
	return t, nil
}

func Required() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaim)
		},
		SigningKey: []byte(envLib.GetEnv("JWT_SECRET")),
	}
	return echojwt.WithConfig(config)
}
