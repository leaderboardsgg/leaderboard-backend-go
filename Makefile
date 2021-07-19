# Set up commands to be used on each environment.
ifeq ($(OS),Windows_NT)
    RMDEL := del
else
    RMDEL := rm
endif

.PHONY: app/graphql_server
gql:
	go build ./app/graphql_server/main.go

gql_run:
	go run ./app/graphql_server/main.go

# A temporary coverprofile file needs written in order to report coverage statistics.
test:
	go test ./... -v -race -coverprofile ./.tmpcover.out
	$(RMDEL) "./.tmpcover.out"

# Running a benchmark multiple times allows better comparison, especially with the benchstat tool.
bench:
	go test  -benchmem -bench=. ./... -run=^$ -v -count 5

# Build testdata database.
build_test_db:
	cd ./sql && docker build -t "speedrun-leaderboard-postgres-test" -f ./Dockerfile .

# Run database with testdata.
run_test_db:
	docker run -d -p 5432:5432 --name pg "speedrun-leaderboard-postgres-test"

