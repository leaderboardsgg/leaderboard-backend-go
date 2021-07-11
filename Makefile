.PHONY: app/graphql_server
gql:
	go build ./app/graphql_server/main.go

gql_run:
	go run ./app/graphql_server/main.go

# A temporary coverprofile file needs written in order to report coverage statistics.
test:
	go test ./... -v -race -coverprofile ./.tmpcover.out
	del *.tmpcover.out

# Running a benchmark multiple times allows better comparison, especially with the benchstat tool.
bench:
	go test  -benchmem -bench=. ./... -run=^$ -v -count 5
