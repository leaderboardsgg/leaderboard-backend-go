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

# A temporary coverprofile file needs written in order to report coverage statistics.
test:
	go test ./... -v -coverprofile ./.tmpcover.out
	$(RMDEL) "./.tmpcover.out"

test_race:
	go test ./... -race

# Running a benchmark multiple times allows better comparison, especially with the benchstat tool.
bench:
	go test  -benchmem -bench=. ./... -run=^$ -v -count 5
