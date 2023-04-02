# article-dispatcher
Article dispatchcer service for channel nine

## TODO
1. create code structure
2. add packages (logs,configs,metrics)
3. define interfaces
4. implement interfaces 
5. routers
6. documentation
7. e2e testing

### Quick run
#### Locally
Prerequisite 
1. [Install Go](https://go.dev/doc/install) on the local machine version 1.17 or higher
2. Install make in order to easy run with the included Makefile - [reference](https://makefiletutorial.com/)

Steps

1. clone the project from the github using below command.  
    `git clone https://github.com/pckushan/article-dispatcher.git`
    This will clone the project into a new folder name `article-dispatcher` 
2. run `go mod download` to download any dependency libraries
3. run `make unit_tests` or `go test -v -count=1 ./...` to check unit tests`results
4. run `make lint` or `golangci-lint run -v `to check linters results
5. run `make build` or `go build -o articel-dispatcher` 
    NOTE: to run on MAC
    run `make build_for_mac` or `env GOOS=darwin GOARCH=amd64 go build -o articel-dispatcher`
6. run `make run` or `./article-dispatcher` to run the article dispatcher service on default port [8888] for the env `HTTP_SERVER_HOST`