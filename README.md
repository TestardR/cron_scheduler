## To run the application

Here is an example of commands you could run
```
go run cmd/main.go -schedule="$(cat input.txt)" -offset=16:10
```

## To run tests

```
make test-fast
```

cron.go holds the core business logic, its coverage reaches 90.7% of statements.
```
github.com/TestardR/cron_scheduler/internal/cron        0.028s  coverage: 90.7% of statements
```
