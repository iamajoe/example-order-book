# Example order book

## Run

### Locally

```sh
go mod download
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

## TODO

- "top of book" should be saved on the repository (easier / faster fetch)
- match and trade
- test domain order
- more tests overall

## Notes

- Used a pattern based on DDD so that the code is more readable and manageable