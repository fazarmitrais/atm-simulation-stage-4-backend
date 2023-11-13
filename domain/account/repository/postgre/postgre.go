package postgre

import (
	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Postgre struct {
	db *gorm.DB
}

func NewPostgre(db *gorm.DB) *Postgre {
	return &Postgre{db}
}

func (p *Postgre) BatchInsert(ctx echo.Context, accounts []*entity.Account) error {
	return p.db.Create(&accounts).Error
}

func (p *Postgre) GetByAccountNumber(ctx echo.Context, accountNumber string) (*entity.Account, error) {
	var account entity.Account
	tx := p.db.Where("account_number = ?", accountNumber).Find(&account)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &account, nil
}

func (p *Postgre) GetAll(ctx echo.Context) ([]*entity.Account, error) {
	var accounts []*entity.Account
	tx := p.db.Find(&accounts)
	if tx.Error == gorm.ErrEmptySlice {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return accounts, nil
}
