#SHELL:=/bin/bash

.PHONY:

build:
	go build -o articel-dispatcher

build_for_mac:
	env GOOS=darwin GOARCH=amd64 go build -o articel-dispatcher

run:
	./articel-dispatcher

unit_tests:
	go test -v -count=1 ./...

e2e_test:
	go test -tags=e2e ./e2e-test -v -count=1

lint:
	golangci-lint run -v

docker_build:
	docker build -t article-dispatcher .

docker_run:
	docker run -p 8888:8888 article-dispatcher