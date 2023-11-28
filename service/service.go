package service

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	trxRepository "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/repository"
	"github.com/labstack/echo/v4"
)

type Service struct {
	AccountRepository     repository.AccountRepository
	AccountCsvRepository  repository.AccountCsvRepository
	TransactionRepository trxRepository.TransactionRepository
}

func NewService(accountRepository repository.AccountRepository, accountCsvRepository repository.AccountCsvRepository, transactionRepository trxRepository.TransactionRepository) *Service {
	return &Service{accountRepository, accountCsvRepository, transactionRepository}
}

type ServiceInterface interface {
	Insert(ctx echo.Context, account entity.Account) *echo.HTTPError
	PINValidation(c echo.Context, account entity.Account) (*string, *echo.HTTPError)
	Transfer(ctx echo.Context, transfer trxEntity.Transaction) (*entity.AccountResponse, *echo.HTTPError)
	BalanceCheck(ctx echo.Context, acctNbr string) (*entity.Account, *echo.HTTPError)
	Import(path string) *echo.HTTPError
	GetAll() ([]*entity.Account, error)
	GetByAccountNumber(c echo.Context, acctNbr string) (*entity.Account, *echo.HTTPError)
	GetLastTransaction(c echo.Context, accountNumber string, transactionType *string, numOfLastTransaction int) ([]*trxEntity.Transaction, error)
}
