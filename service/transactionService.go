package service

import (
	"sort"
	"time"

	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/labstack/echo/v4"
)

func (s *Service) CreateTransactionHistory(c echo.Context, transaction trxEntity.Transaction) error {
	transaction.Date = time.Now()
	return s.TransactionRepository.Add(c, &transaction)
}

func (s *Service) GetLastTransaction(c echo.Context, accountNumber string, numOfLastTransaction int) ([]*trxEntity.Transaction, error) {
	trxs, err := s.TransactionRepository.GetLastTransaction(c, accountNumber, numOfLastTransaction)
	if err != nil {
		return nil, err
	}
	sort.Slice(trxs, func(x, y int) bool {
		return trxs[x].Date.After(trxs[y].Date)
	})
	return trxs, nil
}
