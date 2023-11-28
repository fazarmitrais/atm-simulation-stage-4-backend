package repository

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/labstack/echo/v4"
)

type AccountCsvRepository interface {
	Import(path string) ([]*entity.Account, *echo.HTTPError)
}
