package repository

import (
	"context"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
)

type AccountRepository interface {
	Store(ctx context.Context, accounts []*entity.Account) error
	GetByAccountNumber(ctx context.Context, accountNumber string) (*entity.Account, error)
	GetAll(ctx context.Context) ([]*entity.Account, error)
}
