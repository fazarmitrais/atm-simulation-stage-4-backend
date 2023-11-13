package mapdatastore

import (
	"context"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
)

type MapDatastore struct {
	Transaction map[string][]*entity.Transaction
}

func NewMapDatastore() *MapDatastore {
	m := &MapDatastore{}
	m.Transaction = make(map[string][]*entity.Transaction)
	return m
}

func (m *MapDatastore) Add(ctx context.Context, transaction *entity.Transaction) error {
	m.Transaction[transaction.AccountNumber] = append(m.Transaction[transaction.AccountNumber], transaction)
	return nil
}

func (m *MapDatastore) GetLastTransaction(ctx context.Context, accountNumber string, numberOfTransaction int) ([]*entity.Transaction, error) {
	if numberOfTransaction < 1 {
		numberOfTransaction = 10
	}
	startIndex := len(m.Transaction[accountNumber]) - numberOfTransaction
	if startIndex < 0 {
		startIndex = 0
	}
	return m.Transaction[accountNumber][startIndex:], nil
}
