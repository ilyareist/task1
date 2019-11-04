package payment

import (
	"context"
	"github.com/shopspring/decimal"

	"github.com/ilyareist/task1/account"

	"github.com/go-kit/kit/endpoint"
)

type errorOnlyResponse struct {
	Err error `json:"error,omitempty"`
}

func (r errorOnlyResponse) ErrError() error { return r.Err }

type newPaymentRequest struct {
	FromAccountID account.ID      `json:"from" valid:"alphanum,required,stringlength(1|255)"`
	Amount        decimal.Decimal `json:"amount" valid:"decimal,required"`
	ToAccountID   account.ID      `json:"to" valid:"alphanum,required,stringlength(1|255)"`
}

func makeNewPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newPaymentRequest)
		err := s.New(req.FromAccountID, req.Amount, req.ToAccountID)
		return errorOnlyResponse{Err: err}, nil
	}
}

type newDepositRequest struct {
	AccountID account.ID      `json:"account" valid:"alphanum,required,stringlength(1|255)"`
	Amount        decimal.Decimal `json:"amount" valid:"decimal,required"`
}

func makeDepositEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newDepositRequest)
		err := s.Deposit(req.AccountID, req.Amount)
		return errorOnlyResponse{Err: err}, nil
	}
}

type loadPaymentsRequest struct {
	AccountID account.ID `json:"account"`
}
type ConvertCurrencyRequest struct {
	Amount   decimal.Decimal `json:"amount" valid:"decimal"`
	Currency Currency        `json:"currency" valid:"in(USD|RUB)"`
}

type ConvertCurrencyResponse struct {
	Amount decimal.Decimal `json:"amount"`
}

func makeConvertCurrencyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConvertCurrencyRequest)
		a := s.Convert(req.Amount, req.Currency)
		return ConvertCurrencyResponse{Amount: a}, nil
	}
}

func makeLoadPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loadPaymentsRequest)
		r := s.Load(req.AccountID)
		return r, nil
	}
}

func makeLoadAllPaymentsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := s.LoadAll()
		return r, nil
	}
}
