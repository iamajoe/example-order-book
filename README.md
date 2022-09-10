# Example order book

## Install dependencies

```sh
go mod download
```

## Run

### Locally

```sh
go build -o ./bin/book; \
DB_PATH=bin/data.sqlite ./bin/book <input.csv> <output.csv>
```

### Docker

```sh
docker build  -t orderbook .

# create the folder with the input
mkdir -p ./dockerdata
cp -rf interfaces/cli/fixtures/1_input_balanced_book.csv ./dockerdata/input.csv

docker run --name orderbooking -v $(pwd)/dockerdata:/data/user orderbook
```

## Test

```sh
go fmt ./...
go test ./...
```
