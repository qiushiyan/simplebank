# simplebank


## About

This is a fullstack application that implements a simple banking system,
where users can create accounts, deposit and withdraw money, and
transfer money between accounts.

## Usage

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

## Folder structure

Each folder represents a different layer of concerns

- [`app`](https://github.com/qiushiyan/simplebank/tree/main/app)
  application-level code. Consists of the the Next.js frontend app and
  the backend RESTful API.

- [`zarf`](https://github.com/qiushiyan/simplebank/tree/main/zarf)
  configuration files and static assets

## Routing

### Authentication and Authorization

### Error Handling

## Asynchronous Processing

## Logging

[`foundation/logger`](https://github.com/qiushiyan/simplebank/tree/main/foundation/logger)
exports a customized slog-based logger.

[`app/tooling/logfmt`](https://github.com/qiushiyan/simplebank/tree/main/app/tooling/logfmt)
is a simple program that converts the JSON logs output by
`foundation/logger` to human-readable logs.
