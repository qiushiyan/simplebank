# simplebank

## Logging

`foundation/logger` exports a customized `SugaredLogger` using the [`zap`](https://github.com/uber-go/zap) package.

`app/tooling/logfmt` is a simple program that converts the JSON logs output by `foundation/logger` to human-readable logs.

You can use it like this:

```bash
# make run-local
go run app/services/bank-api/main.go | go run app/tooling/logfmt/main.go
```
