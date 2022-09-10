# Example Order book

## Install

```go
go mod download
```

## Run

```goa
go build -o ./bin/book; \
DB_PATH=bin/data.sqlite ./bin/book <userID> <new|cancel|flush> [input.csv]
```

## Test

```
go fmt
go vet
go test
```
