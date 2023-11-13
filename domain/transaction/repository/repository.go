package repository

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/labstack/echo/v4"
)

type TransactionRepository interface {
	Add(ctx echo.Context, transaction *entity.Transaction) error
	GetLastTransaction(ctx echo.Context, accountNumber string, numberOfTransaction int) ([]*entity.Transaction, error)
}
