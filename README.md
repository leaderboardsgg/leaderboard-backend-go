# leaderboard-backend
An open-source community-driven leaderboard backend for the gaming community.

## Links
- Website: https://speedrun.website/
- Other Repos: https://github.com/leaderboardsgg
- Discord: https://discord.gg/TZvfau25Vb

# Tech-Stack Information
- This repository only contains the backend, and not the UI for the website.
- GoLang is used for implementing the backend.
- JSON API with JWT Authentication

# Developing
## Requirements
- [Go](https://golang.org/doc/install) 1.16+.
- [Make](https://www.gnu.org/software/make/) to run build scripts.
- [Docker](https://hub.docker.com/search?q=&type=edition&offering=community) to run the database and admin interface.
## Optional
- [golangci-lint](https://golangci-lint.run/usage/install/) to run CI linting on your machine.
- [staticcheck](https://staticcheck.io/docs/install) for linting that will integrate with your editor.
- [gcc](https://gcc.gnu.org/) to run race detection on unit tests.

## Useful links
- [VSCode](https://code.visualstudio.com/download) is a pretty good editor with helpful GoLang plugins.
- [GoLand](https://www.jetbrains.com/go/) is JetBrains' Go offering and is very fully featured.
- [A Tour of Go](https://tour.golang.org/welcome/1) is a great place to learn the basics of how to use GoLang.
- [Effective Go](https://golang.org/doc/effective_go) is the best place to check to learn recommended best practices.
- [Set up git](https://docs.github.com/en/get-started/quickstart/set-up-git) is GitHub's guide on how to set up and begin using git.
- [How to Contribute to an Open Source Project](https://opensource.guide/how-to-contribute/#opening-a-pull-request) is a useful guide showing some of the steps involved in opening a Pull Request.

## How to run
To start the postgres docker container
- `docker-compose up -d`
- Go to `localhost:1337` for an Adminer interface

To test HTTP endpoints:
- `make run` or `make build` and run the binary
- Make requests to `localhost:3000/api/v1` (or whatever port from .env)

Running tests:
- `go test ./...`

Running tests with coverage:
- `make test`

Running tests with race detection (requires `gcc`):
- `make test_race`

Running benchmarks:
- `make bench`
