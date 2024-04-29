# simplebank


## About

This is a fullstack application that implements a simple banking system,
where users can create accounts, deposit and withdraw money, and
transfer money between accounts.

## Features

### CRUD operations

CRUD operations are implemented for the following domains:

- **user**:
  - `POST /signup`: create user
- **account**:
  - `GET /accounts`, `GET /accounts/{id}` and `GET /accounts/search`:
    query account by username (from bearer token), id or account name.

  - `POST /accounts/create`: create accounts

  - `PATCH /accounts/{id}`: update account
- **transfer**:
  - `POST /trasnfer`: transfer money between accounts, update account
    balance and create entry records
- **entry**:
  - `GET /entries`: list entries
- **friendship**:
  - `GET /friend/list`: query incoming and outgoing friendship requests
    for account
  - `POST /friend/create`: create friendship requests
  - `PATCH /friend/{id}`: update friendship request status

### Authentication and authorization

Routes for querying and mutating resources require a bearer token, which
can be obtained from the `POST /signin` route. The token is a
[paseto](https://paseto.io/) token that carries the username and a role.
User can either be a regular end user or an admin. Currently, there is
only one route `GET /accounts/all` that requires admin access.

### Asynchronous processing

The API provides two options for asynchronous task processing that is
configured via the `TASK_MANAGER` env variable:

- `"simple"`: using a simple goroutine and channel based task queue

- `"asynq"`: using the [asynq](https://github.com/hibiken/asynq) library
  which is backed by redis.

### Email Sending

The `/send-email` route supports sending emails to an authenticated
user. You can specify `subject=verify` `subject=welcome`, and
`subject=report` in the request body to send emails of different types.

## Development

If you have docker-compose installed, simply run

``` bash
docker compose up
```

To run it without docker, first install the go dependencies

``` bash
go mod tidy
```

Start the backend with

``` bash
make run-local
# or
go run app/services/bank-api/main.go | go run app/tooling/logfmt/main.go
```

### Documentation

The backend serves a swagger documentation at `/swagger/index.html`

To update the swagger docs, make sure you have `swag` installed

``` bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Then run

``` bash
make docs
```

## Design Details

### Folder structure

Each folder represents a different layer of concerns

- [`app`](https://github.com/qiushiyan/simplebank/tree/main/app)
  application-level code. Consists of the the Next.js frontend app and
  the backend RESTful API.

- [`zarf`](https://github.com/qiushiyan/simplebank/tree/main/zarf)
  configuration files and static assets

### Routing

### Authentication and Authorization

### Error Handling

### Asynchronous Processing

### Logging

[`foundation/logger`](https://github.com/qiushiyan/simplebank/tree/main/foundation/logger)
exports a customized slog-based logger.

[`app/tooling/logfmt`](https://github.com/qiushiyan/simplebank/tree/main/app/tooling/logfmt)
is a simple program that converts the JSON logs output by
`foundation/logger` to human-readable logs.
