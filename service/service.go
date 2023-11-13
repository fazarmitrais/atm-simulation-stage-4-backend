package service

import (
	"context"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository"
	trxRepository "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/repository"
	"github.com/fazarmitrais/atm-simulation-stage-3/lib/responseFormatter"
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
	PINValidation(c context.Context, account entity.Account) *responseFormatter.ResponseFormatter
	Transfer(ctx, transfer entity.Transfer) (*entity.Account, *responseFormatter.ResponseFormatter)
	BalanceCheck(ctx context.Context, acctNbr string) (*entity.Account, *responseFormatter.ResponseFormatter)
	Import() error
	GetAll() ([]*entity.Account, error)
}
