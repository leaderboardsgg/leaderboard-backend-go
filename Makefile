# Set up commands to be used on each environment.
ifeq ($(OS),Windows_NT)
    RMDEL := del
else
    RMDEL := rm
endif

.PHONY: build
build:
	go build ./main.go

run:
	go run ./main.go

# generate models for graphql
generate:
	go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate

# A temporary coverprofile file needs written in order to report coverage statistics.
test:
	go test ./... -v -race -coverprofile ./.tmpcover.out
	$(RMDEL) "./.tmpcover.out"

# Running a benchmark multiple times allows better comparison, especially with the benchstat tool.
bench:
	go test  -benchmem -bench=. ./... -run=^$ -v -count 5