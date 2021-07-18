# leaderboard-backend
An open-source community-driven leaderboard backend for the gaming community.

## Links
- Website: https://speedrun.website/
- Other Repos: https://github.com/speedrun-website
- Discord: https://discord.gg/TZvfau25Vb

# Tech-Stack Information
- This repository only contains the backend, and not the UI for the website.
- GoLang is used for implementing the backend.
- Only a GraphQL API (with a JSON POST endpoint) is being implemented so far.

# Developing
## Requirements
- [Go](https://golang.org/doc/install) 1.11+.
- [Make](https://www.gnu.org/software/make/) to run build scripts.
- [gcc](https://gcc.gnu.org/) to run race detection on unit tests.
- [npm](https://www.npmjs.com/) to run newman.
- [newman](https://learning.postman.com/docs/running-collections/using-newman-cli/command-line-integration-with-newman/) to run integration tests.

## Useful links
- [VSCode](https://code.visualstudio.com/download) is a pretty good editor with helpful GoLang plugins.
- [A Tour of Go](https://tour.golang.org/welcome/1) is a great place to learn the basics of how to use GoLang.
- [Set up git](https://docs.github.com/en/get-started/quickstart/set-up-git) is Githubs guide on how to set up and begin using git.
- [How to Contribute to an Open Source Project](https://opensource.guide/how-to-contribute/#opening-a-pull-request) is a useful guide showing some of the steps involved in opening a Pull Request.

## How to run
To get an interactive API interface:
- `make gql_run`
- Go to `localhost:3030/graphiql`

To test HTTP endpoints:
- `make gql_run`
- Make POST requests to `localhost:3030/graphql/http`

Running tests:
- `make test`

Running benchmarks:
- `make bench`

## Running in a container
- `docker build . -t graphql_server`
- `docker run -p 3030:3030 -d graphql_server`
- Voila! Go to `localhost:3030/graphql`
