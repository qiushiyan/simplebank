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

To run it without docker

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

The project is organized into the following folders, each representing a
different layer of the application:

- [apps](https://github.com/qiushiyan/simplebank/blob/main/apps)

## Routing

### Authentication and Authorization

### Error Handling

## Asynchronous Processing

## Logging

`foundation/logger` exports a customized `SugaredLogger` using the
[`zap`](https://github.com/uber-go/zap) package.

`app/tooling/logfmt` is a simple program that converts the JSON logs
output by `foundation/logger` to human-readable logs.
