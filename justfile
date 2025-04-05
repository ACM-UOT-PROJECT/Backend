jet:
    jet -source="${DB}" -dsn="$GOOSE_DBSTRING" -path=./database/gen/

build:
    go build -o ./tmp/main .

all: jet build

test-all:
    go test ./...

run: all
    ./tmp/main
