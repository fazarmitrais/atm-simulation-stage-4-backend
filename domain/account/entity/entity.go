package entity

type Account struct {
	Name          string  `json:"name"`
	AccountNumber string  `json:"accountNumber"`
	PIN           string  `json:"pin"`
	Balance       float64 `json:"balance"`
}

type AccountResponse struct {
	Name          string  `json:"name"`
	AccountNumber string  `json:"accountNumber"`
	Balance       float64 `json:"balance"`
}

type Transfer struct {
	FromAccountNumber string  `json:"fromAccountNumber"`
	ToAccountNumber   string  `json:"toAccountNumber"`
	ReferenceNumber   string  `json:"referenceNumber"`
	Amount            float64 `json:"amount"`
}

func (a *Account) ToAccountResponse() *AccountResponse {
	return &AccountResponse{
		Name:          a.Name,
		AccountNumber: a.AccountNumber,
		Balance:       a.Balance,
	}
}
