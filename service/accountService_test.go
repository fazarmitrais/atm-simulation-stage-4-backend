package service

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/fazarmitrais/atm-simulation-stage-3/domain/account/entity"
	trxEntity "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/entity"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccountRepositoryMock struct {
	mock.Mock
}

type TransactionRepositoryMock struct {
	mock.Mock
}

const (
	GetByAccountNumber       = "GetByAccountNumber"
	BatchInsert              = "BatchInsert"
	CreateTransactionHistory = "CreateTransactionHistory"
)

var Acct112233 = entity.Account{
	Name:          "John Doe",
	AccountNumber: "112233",
	PIN:           "012108",
	Balance:       100,
}

var Acct112244 = entity.Account{
	Name:          "Jane Doe",
	AccountNumber: "112244",
	PIN:           "932012",
	Balance:       100,
}

func (mt *TransactionRepositoryMock) Add(ctx echo.Context, tx *trxEntity.Transaction) error {
	mt.Called(ctx, tx)
	return nil
}

func (mt *TransactionRepositoryMock) GetLastTransaction(ctx echo.Context, accountNumber string, numOfLastTransaction int) ([]*trxEntity.Transaction, error) {
	mt.Called(ctx, accountNumber, numOfLastTransaction)
	return nil, nil
}

func (m *AccountRepositoryMock) GetByAccountNumber(ctx echo.Context, accountNumber string) (*entity.Account, error) {
	args := m.Called(ctx, accountNumber)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountRepositoryMock) GetAll(ctx echo.Context) ([]*entity.Account, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}
	return args.Get(0).([]*entity.Account), args.Error(1)
}

func (m *AccountRepositoryMock) BatchInsert(ctx echo.Context, accounts []*entity.Account) error {
	args := m.Called(ctx, accounts)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func TestPinValidationAccountNbrIsRequired(t *testing.T) {
	svc := NewService(nil, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		PIN: "456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Account Number is required", resp.Message)
}

func TestPinValidationPINIsRequired(t *testing.T) {
	svc := NewService(nil, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: "123",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "PIN is required", resp.Message)
}

// - Account Number should have 6 digits length. Display message `Account Number should have 6 digits length` for invalid Account Number.
func TestPinValidationAccountNumberMustSixDigitsLength(t *testing.T) {
	svc := NewService(nil, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: "123",
		PIN:           "456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Account Number should have 6 digits length", resp.Message)
}

//- PIN should have 6 digits length. Display message `PIN should have 6 digits length` for invalid PIN.

func TestPinValidationPINMustSixDigitsLength(t *testing.T) {
	svc := NewService(nil, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: "123456",
		PIN:           "456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "PIN should have 6 digits length", resp.Message)
}

// - Account Number should only contains numbers [0-9]. Display message `Account Number should only contains numbers` for invalid Account Number.
func TestPinValidationAccountNumberOnlyContainsNumber(t *testing.T) {
	svc := NewService(nil, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: "a123456",
		PIN:           "123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Account Number should only contains numbers", resp.Message)
}

// - PIN should only contains numbers [0-9]. Display message `PIN should only contains numbers` for invalid PIN.
func TestPinValidationPINOnlyContainsNumber(t *testing.T) {
	svc := NewService(nil, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: "123456",
		PIN:           "a123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "PIN should only contains numbers", resp.Message)
}

//- Check valid Acccount Number & PIN with ATM records. Display message `Invalid Account Number/PIN` if records is not exist.

func TestPinValidationInvalidAccountNumber(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, "123456").Return(nil, nil)
	svc := NewService(repo, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: "123456",
		PIN:           "1123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid Account Number/PIN", resp.Message)
}

// - Check valid Acccount Number & PIN with ATM records. Display message `Invalid Account Number/PIN` if records is not exist.
func TestPinValidationInvalidPIN(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	svc := NewService(repo, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: Acct112233.AccountNumber,
		PIN:           "1123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid Account Number/PIN", resp.Message)
}

func TestPinValidationSuccess(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	svc := NewService(repo, nil, nil)
	resp := svc.PINValidation(echo.New().AcquireContext(), entity.Account{
		AccountNumber: Acct112233.AccountNumber,
		PIN:           Acct112233.PIN,
	})
	assert.Nil(t, resp)
}

// - Maximum amount to withdraw is $1000. Display message `Maximum amount to withdraw is $1000` if withdraw amount is higher than $1000.
func TestWithdrawMaxAmount1000(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Withdraw(echo.New().AcquireContext(), Acct112233.AccountNumber, 1001)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Maximum amount to withdraw is $1000", resp.Message)
}

// - Display message `Invalid ammount` if withdraw amount is not multiple of $10.
func TestWithdrawAmountNotMultipleOf10(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Withdraw(echo.New().AcquireContext(), Acct112233.AccountNumber, 901)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid ammount", resp.Message)
}

// - Display message `Insufficient balance $10` for insufficient balance. `$10` is the withdraw amount
func TestWithdrawInsufficientBalance(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	svc := NewService(repo, nil, nil)
	var withdrawAmt float64 = 200
	_, resp := svc.Withdraw(echo.New().AcquireContext(), Acct112233.AccountNumber, withdrawAmt)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("Insufficient balance $%0.f", withdrawAmt), resp.Message)
}

// - Display message `Invalid account` if account is not numbers
func TestTransferAccountMustBeNumbers(t *testing.T) {
	svc := NewService(nil, nil, nil)
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: "a432214213",
		ToAccountNumber:   "a432214214",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid account", resp.Message)
}

// - Display message `Invalid account` if account is not found
func TestTransferFromAccountNumberMustBeCorrect(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, "432214213").Return(nil, nil)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: "432214213",
		ToAccountNumber:   Acct112233.AccountNumber,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid account", resp.Message)
}

// - Display message `Invalid account` if account is not found
func TestTransferToAccountNumberMustBeCorrect(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	repo.On(GetByAccountNumber, mock.Anything, "432214214").Return(nil, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "432214214",
		Amount:            10,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid account", resp.Message)
}

// - Maximum amount to transfer is $1000. Display message `Maximum amount to transfer is $1000` if transfer amount is higher than $1000.
func TestTransferMaxTransferAmountIs1000(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	repo.On(GetByAccountNumber, mock.Anything, Acct112244.AccountNumber).Return(&Acct112244, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: Acct112233.AccountNumber,
		ToAccountNumber:   Acct112244.AccountNumber,
		Amount:            1001,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Maximum amount to transfer is $1000", resp.Message)
}

// - Minimum amount to transfer is $1. Display message `Minimum amount to transfer is $1` if transfer amount is lower than $1.
func TestTransferMinTransferAmountIs1(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	repo.On(GetByAccountNumber, mock.Anything, Acct112244.AccountNumber).Return(&Acct112244, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: Acct112233.AccountNumber,
		ToAccountNumber:   Acct112244.AccountNumber,
		Amount:            0.5,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Minimum amount to transfer is $1", resp.Message)
}

// - Display message `Insufficient balance $300` for insufficient balance. `$300` is the transfer amount
func TestTransferInsufficientBalance(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	repo.On(GetByAccountNumber, mock.Anything, Acct112244.AccountNumber).Return(&Acct112244, nil)
	svc := NewService(repo, nil, nil)
	var trfAmount float64 = 200
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: Acct112233.AccountNumber,
		ToAccountNumber:   Acct112244.AccountNumber,
		Amount:            trfAmount,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("Insufficient balance $%0.f", trfAmount), resp.Message)
}

// - Display message `Invalid Reference Number` if reference number is not empty and not numbers
func TestTransferReferenceNumberMustBeNumber(t *testing.T) {
	repo := new(AccountRepositoryMock)
	repo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	repo.On(GetByAccountNumber, mock.Anything, Acct112244.AccountNumber).Return(&Acct112244, nil)
	svc := NewService(repo, nil, nil)
	_, resp := svc.Transfer(echo.New().AcquireContext(), entity.Transfer{
		FromAccountNumber: Acct112233.AccountNumber,
		ToAccountNumber:   Acct112244.AccountNumber,
		Amount:            20,
		ReferenceNumber:   "Ref 213342",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid Reference Number", resp.Message)
}

// - Valid amount will deduct the user balance with transfer amount and will add destination account with transfer amount. After that screen will
func TestTransferSuccess(t *testing.T) {
	accRepo := new(AccountRepositoryMock)
	trxRepo := new(TransactionRepositoryMock)
	accRepo.On(GetByAccountNumber, mock.Anything, Acct112233.AccountNumber).Return(&Acct112233, nil)
	accRepo.On(GetByAccountNumber, mock.Anything, Acct112244.AccountNumber).Return(&Acct112244, nil)
	trxRepo.On("Add", mock.Anything, mock.Anything).Return(nil, nil)
	svc := NewService(accRepo, nil, trxRepo)
	ctx := echo.New().AcquireContext()
	_, resp := svc.Transfer(ctx, entity.Transfer{
		FromAccountNumber: Acct112233.AccountNumber,
		ToAccountNumber:   Acct112244.AccountNumber,
		Amount:            20,
		ReferenceNumber:   "213342",
	})
	fromAcct, _ := svc.BalanceCheck(ctx, "112233")
	toAcct, _ := svc.BalanceCheck(ctx, "112244")
	assert.Nil(t, resp)
	assert.Equal(t, float64(80), fromAcct.Balance)
	assert.Equal(t, float64(120), toAcct.Balance)
}
