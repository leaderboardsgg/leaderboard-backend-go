.PHONY: app/graphql_server
gql:
	go build ./app/graphql_server/main.go

gql_run:
	go run ./app/graphql_server/main.go
