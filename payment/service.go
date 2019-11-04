// Package payment provides handlers for work with payments in the system.
package payment

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilyareist/task1/account"
	"github.com/ilyareist/task1/errs"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
)

// Direction of payment regarding account.
type Direction string
type Currency string

const (
	Incoming Direction = "incoming"
	Outgoing Direction = "outgoing"
)

// Payment holding a money transfer between two accounts in the system.
type Payment struct {
	ID          uuid.UUID       `json:"-" sql:"id,pk,type:varchar(36)"`
	Account     account.ID      `json:"account" sql:"type:varchar(255)" pg:"fk:base_account_id"`
	Amount      decimal.Decimal `json:"amount" sql:"amount,notnull,type:'decimal(16,4)'"`
	ToAccount   account.ID      `json:"to_account,omitempty" sql:"to_account,type:varchar(255)" pg:"fk:to_account_id"`
	FromAccount account.ID      `json:"from_account,omitempty" sql:"from_account,type:varchar(255)" pg:"fk:from_account_id"`
	Direction   Direction       `json:"direction" sql:"direction,notnull,type:varchar(16)"`
	Deleted     bool            `json:"-" sql:"deleted,notnull"`
}

type Rate struct {
	Currency string  `json:"currency" sql:"type:varchar(255)"`
	Date     string  `json:"date" sql:"type:varchar(255)"`
	Rate     float64 `json:"rate" sql:"type:float"`
}

// Service is the interface that provides payment methods.
type Service interface {
	// New registers a new payment in the system.
	New(fromAccountID account.ID, amount decimal.Decimal, toAccountID account.ID) error

	// Load returns payments list for an account.
	Load(accountID account.ID) []*Payment

	// LoadAll returns all payments, registered in the system.
	LoadAll() []*Payment

	// Show rate on the specific date
	Rates(currency string, date string) (Rate, error)

	// Update an account returns a read model of an account.
	Deposit(accountID account.ID, amount decimal.Decimal) error
}

type service struct {
	accounts account.Repository
	payments Repository
}

// New registers a new payment in the system.
func (s *service) New(fromAccountID account.ID, amount decimal.Decimal, toAccountID account.ID) error {
	if fromAccountID == toAccountID {
		return errs.ErrAccountsAreEqual
	}
	from, err := s.accounts.Find(fromAccountID)
	if err != nil {
		return errs.ErrUnknownSourceAccount
	}
	fmt.Println(from.Currency)

	var fromAmount decimal.Decimal
	if from.Currency != "USD" {
		date := "latest"
		currency := string(from.Currency)
		rate, _ := s.Rates(currency, date)
		fromAmount = amount.Mul(decimal.NewFromFloat(rate.Rate))
	} else {
		fromAmount = amount
	}

	if from.Balance.LessThan(fromAmount) {
		return errs.ErrInsufficientMoney
	}


	to, err := s.accounts.Find(toAccountID)
	//fmt.Println(to.Currency)
	if err != nil {
		return errs.ErrUnknownTargetAccount
	}

	var toAmount decimal.Decimal
	if to.Currency != "USD" {
		date := "latest"
		currency := string(to.Currency)
		rate, _ := s.Rates(currency, date)
		toAmount = amount.Mul(decimal.NewFromFloat(rate.Rate))
	} else {
		toAmount = amount
	}


	outgoingPayment := Payment{
		ID:        uuid.New(),
		Account:   fromAccountID,
		Amount:    fromAmount,
		ToAccount: toAccountID,
		Direction: Outgoing,
	}
	incomingPayment := Payment{
		ID:          uuid.New(),
		Account:     toAccountID,
		Amount:      toAmount,
		FromAccount: fromAccountID,
		Direction:   Incoming,
	}
	err = s.payments.Store(&outgoingPayment, &incomingPayment)
	if err != nil {
		return errs.ErrStorePayments
	}
	return nil
}

func (s *service) Deposit(accountID account.ID, amount decimal.Decimal) error {

	account, err := s.accounts.Find(accountID)
	if err != nil {
		return errs.ErrUnknownSourceAccount
	}
	var amountUSD decimal.Decimal
	fmt.Println(account.Currency)
	if account.Currency != "USD" {
		date := "latest"
		currency := string(account.Currency)
		rate, _ := s.Rates(currency, date)
		amountUSD = amount.Mul(decimal.NewFromFloat(rate.Rate))
	} else {
		amountUSD = amount
	}

	incomingPayment := Payment{
		ID:          uuid.New(),
		Account:     accountID,
		Amount:      amountUSD,
		FromAccount: accountID,
		Direction:   Incoming,
	}
	err = s.payments.Store(&incomingPayment)
	if err != nil {
		return errs.ErrStorePayments
	}
	return nil
}

// Load returns payments list for an account.
func (s *service) Load(accountID account.ID) []*Payment {
	return s.payments.Find(accountID)
}

// LoadAll returns all payments, registered in the system.
func (s *service) LoadAll() []*Payment {
	return s.payments.FindAll()
}

// Convert returns amount in a given currency
func (s *service) Rates(currency string, date string) (Rate, error) {
	var currentRate Rate
	var f map[string]interface{}
	url := "https://api.exchangeratesapi.io/"

	response, err := http.Get(url + date + "?base=USD&&symbols=" + currency)
	body, _ := ioutil.ReadAll(response.Body)

	if err := json.Unmarshal([]byte(body), &f); err != nil {
		return currentRate, err
	}

	k, _ := f["rates"].(map[string]interface{})[currency].(float64)

	currentRate.Currency = currency
	currentRate.Date = date
	currentRate.Rate = k
	return currentRate, err
}

// NewService creates a payment service with necessary dependencies.
func NewService(payments Repository, accounts account.Repository) Service {
	return &service{
		payments: payments,
		accounts: accounts,
	}
}

// Repository interface for payment storing and operations.
type Repository interface {
	// Store payments in the repository.
	Store(payment ...*Payment) error

	// Find payments list for an account.
	Find(id account.ID) []*Payment

	//Deposit(accountID account.ID, amount decimal.Decimal) error
	// FindAll returns all payments, registered in the system.
	FindAll() []*Payment

	// MarkDeleted is mark as deleted specified payment in the system
	MarkDeleted(id uuid.UUID) error
}
