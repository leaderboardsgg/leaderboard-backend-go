FROM golang:1.16-alpine3.13 as builder

WORKDIR /app

# Only copy needed files
COPY app/               app/
COPY data/              data/
COPY graphql_server/    graphql_server/
COPY go.mod             go.mod
COPY go.sum             go.sum

RUN apk add --no-cache ca-certificates
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gql_server -ldflags="-w -s" ./app/graphql_server/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/gql_server /usr/bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["gql_server"]
