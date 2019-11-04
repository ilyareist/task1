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

// Service is the interface that provides payment methods.
type Service interface {
	// New registers a new payment in the system.
	New(fromAccountID account.ID, amount decimal.Decimal, toAccountID account.ID) error

	// Load returns payments list for an account.
	Load(accountID account.ID) []*Payment

	// LoadAll returns all payments, registered in the system.
	LoadAll() []*Payment

	// Update an account returns a read model of an account.
	Deposit(accountID account.ID, amount decimal.Decimal) error

	// Converts USD to the currency
	Convert(amount decimal.Decimal, currency Currency) (decimal.Decimal)
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

	if from.Balance.LessThan(amount) {
		return errs.ErrInsufficientMoney
	}
	to, err := s.accounts.Find(toAccountID)
	fmt.Println(to.Currency)
	if err != nil {
		return errs.ErrUnknownTargetAccount
	}

	outgoingPayment := Payment{
		ID:        uuid.New(),
		Account:   fromAccountID,
		Amount:    amount,
		ToAccount: toAccountID,
		Direction: Outgoing,
	}
	incomingPayment := Payment{
		ID:          uuid.New(),
		Account:     toAccountID,
		Amount:      amount,
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

	_, err := s.accounts.Find(accountID)
	if err != nil {
		return errs.ErrUnknownSourceAccount
	}

	incomingPayment := Payment{
		ID:          uuid.New(),
		Account:     accountID,
		Amount:      amount,
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
func (s *service) Convert(amount decimal.Decimal, currency Currency) decimal.Decimal {
	res, _ := http.Get("https://api.exchangeratesapi.io/latest?symbols=GBP&&base=USD")
	type people struct {
		Number int `json:"rates"`
	}
	body, _ := ioutil.ReadAll(res.Body)

	people1 := people{}
	_ = json.Unmarshal(body, &people1)


	fmt.Println(people1.Number)
	return amount
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
