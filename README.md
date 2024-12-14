# Dragon ball

This app search and save characters of Dragon ball from external API https://web.dragonball-api.com/

## How to run application

Run docker mode
```sh
docker-compose up -d --build
```

Run unit test:
```sh
# Execute test
go test -timeout 30m -coverprofile=coverage.out -coverpkg=./... ./src/...
# Generate html of coverage
go tool cover -html coverage.out -o coverage.html
# Getting % of coverage
go tool cover -func=coverage.out | awk 'NR>1{print $3}' | awk '{sum+=$1} END {print sum/NR}'
```

Run lint:
```sh
go vet ./...
golangci-lint run --timeout=30m
# remember install libraries of linter in local
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
