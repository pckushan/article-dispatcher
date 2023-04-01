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
1. clone the project from the github 
2. run `go mod tidy` to resolve any dependency libraries
3. run `make unit_tests` to check unit tests`
4. run `make lint` to check linters 
5. run `make build` and `./article-dispatcher`
6. run `make run` to run the article dispatcher service on default port [8888] for the env `HTTP_SERVER_HOST`