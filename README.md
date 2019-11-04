Payments
--------

## Table of Contents

<!-- TOC depthFrom:2 depthTo:6 updateOnSave:true withLinks:true -->

- [Table of Contents](#table-of-contents)
- [Project purpose](#project-purpose)
- [Usage](#usage)
    - [Command-line flags](#command-line-flags)
- [Dependencies](#dependencies)
- [How to set up](#how-to-set-up)
    - [Step 1. Build docker image](#step-1-build-docker-image)
    - [Step 2. Run it](#step-2-run-it)
- [How to run tests](#how-to-run-tests)
- [How to run code linting](#how-to-run-code-linting)
- [How to Contribute](#how-to-contribute)
    - [Support](#support)
    - [What to contribute](#what-to-contribute)
    - [Advice](#advice)

<!-- /TOC -->

## Project purpose

Payment system, provides ability to transfer money between accounts. 

System also provide reports: 
 - all registered accounts; 
 - all registered payments (transfers).

API Documentation see [here](./docs/api.md)

## Usage

### Command-line flags

#### Running locally
To run project locally with docker-compose use:

```bash
docker-compose up
```
This command will create .env file from .env.dist and start Docker cluster with following components:

- Postgresql Database
- Backend http://127.0.0.1:8080 
## Dependencies

- [go-kit](http://github.com/go-kit/kit) -- toolkit for building microservices, recommended by design;
- [govalidator](http://github.com/asaskevich/govalidator) -- package of validators and sanitizers for strings, 
numerical, slices and structures;
- [decimal](http://github.com/shopspring/decimal) -- arbitrary-precision fixed-point decimal numbers in go; 
- [uuid](http://github.com/google/uuid) -- go package for UUIDs based on RFC 4122 and DCE 1.1;
- [gorilla/mux](http://github.com/gorilla/mux) -- a powerful HTTP router and URL matcher for building Go web servers;
- [prometheus client](http://github.com/prometheus/client_golang) -- prometheus instrumentation library for Go
applications;
- [go-cmp](https://github.com/google/go-cmp) -- package for comparing Go values in tests;
- [go-pg](https://github.com/go-pg/pg) -- golang ORM with focus on PostgreSQL features and performance.

## How to set up

### Step 1. Build docker image
```bash
 docker build -t payments-app .
```

### Step 2. Run it

```bash
docker run --rm -p 8080:8080 payments-app --db_address=${DB_ADDR} --db_password=${DB_PASSWORD}
```
