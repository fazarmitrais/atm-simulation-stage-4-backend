package service

import (
	"context"
	"sort"
	"time"

	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
)

func (s *Service) CreateTransactionHistory(transaction trxEntity.Transaction) error {
	transaction.Date = time.Now()
	return s.TransactionRepository.Add(context.Background(), &transaction)
}

func (s *Service) GetLastTransaction(c context.Context, accountNumber string, numOfLastTransaction int) ([]*trxEntity.Transaction, error) {
	trxs, err := s.TransactionRepository.GetLastTransaction(c, accountNumber, numOfLastTransaction)
	if err != nil {
		return nil, err
	}
	sort.Slice(trxs, func(x, y int) bool {
		return trxs[x].Date.After(trxs[y].Date)
	})
	return trxs, nil
}
