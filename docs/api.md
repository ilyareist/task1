# API Documentation

## Table of Contents

<!-- TOC depthFrom:2 depthTo:6 updateOnSave:true withLinks:true -->
- [API Documentation](#api-documentation)
  * [Table of Contents](#table-of-contents)
  * [Accounts Collection `/api/accounts/v1/accounts`](#accounts-collection---api-accounts-v1-accounts-)
    + [List All Accounts](#list-all-accounts)
      - [Request](#request)
    + [Create a New Account](#create-a-new-account)
      - [Request](#request-1)
  * [Account `/api/accounts/v1/accounts/{account_id}`](#account---api-accounts-v1-accounts--account-id--)
    + [Get account by ID](#get-account-by-id)
      - [Request](#request-2)
  * [Payments Collection `/api/payments/v1/payments`](#payments-collection---api-payments-v1-payments-)
    + [List All Payments](#list-all-payments)
      - [Request](#request-3)
    + [Create a New Payment](#create-a-new-payment)
      - [Request](#request-4)
    + [Make a deposit](#make-a-deposit)
      - [Request](#request-5)
    + [Get currency rates to date](#get-currency-rates-to-date)
      - [Request](#request-6)
  * [Payments by Account `/api/payments/v1/payments/{account_id}`](#payments-by-account---api-payments-v1-payments--account-id--)
    + [Get Payments for Account](#get-payments-for-account)

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>
<!-- /TOC -->

## Accounts Collection `/api/accounts/v1/accounts`

### List All Accounts

Returns all accounts registered in the system.

#### Request

**URL**: `/api/accounts/v1/accounts`  
**Method**: `GET`

```bash
curl --include \
'http://0.0.0.0:8080/api/accounts/v1/accounts'
```


### Create a New Account

You may create new account using this action. It takes a JSON object containing an id, initial balance and currency.

#### Request

**URL**: `/api/accounts/v1/accounts`  
**Method**: `POST`

```bash
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
    \"id\": \"John\",
    \"country\": \"USA\",
    \"city\": \"Colorado\",
    \"id\": \"john789\",
    \"balance\": 55.00,
    \"currency\": \"USD\"
}" \
'http://0.0.0.0:8080/api/accounts/v1/accounts'
```


## Account `/api/accounts/v1/accounts/{account_id}`

### Get account by ID

Returns a read model of an account.

#### Request

**URL**: `/api/accounts/v1/accounts/{account_id}`  
**Method**: `GET`  

```bash
curl --include \
'http://0.0.0.0:8080/api/accounts/v1/accounts/John'
```


## Payments Collection `/api/payments/v1/payments`

### List All Payments

Returns all payments, registered in the system.

#### Request

**URL**: `/api/payments/v1/payments`  
**Method**: `GET`

```bash
curl --include \
'http://0.0.0.0:8080/api/payments/v1/payments'
```


### Create a New Payment

Creates a new payment.

#### Request

**URL**: `/api/payments/v1/payments`  
**Method**: `POST`

```bash
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
    \"from\": \"John\",
    \"amount\": 12.34,
    \"to\": \"Ivan\"
}" \
'http://0.0.0.0:8080/api/payments/v1/payments'
```

### Make a deposit

Deposit to account's balance 

#### Request

**URL**: `http://0.0.0.0:8080/api/payments/v1/payments/deposit`  
**Method**: `POST`

```bash
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
    \"account\": \"John\",
    \"amount\": 12.34,

}" \
'http://0.0.0.0:8080/api/payments/v1/payments/deposit'
```

### Get currency rates to date


#### Request

**URL**: `http://0.0.0.0:8080/api/payments/v1/payments/rates`  
**Method**: `POST`

```bash
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
    \"currency\": \"RUB\",
    \"date\": \"2019-1-12\"

}" \
'http://0.0.0.0:8080/api/payments/v1/payments/rates'
```
## Payments by Account `/api/payments/v1/payments/{account_id}`

### Get Payments for Account

Returns payments list for an account.

