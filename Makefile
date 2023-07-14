#SHELL=/usr/bin/env bash

all:
	go mod tidy \
	&& go build -ldflags "-s -w" -o client example/client/main.go \
	&& go build -ldflags "-s -w" -o server example/server/main.go

client:
	go mod tidy && go build -ldflags "-s -w" -o client example/client/main.go

server:
	go mod tidy && go build -ldflags "-s -w" -o server example/server/main.go

.PHONY: all client server
BINS+=client
BINS+=server

docker:#env-GIT_USER env-GIT_PASSWORD
	docker build --tag gomicro-services -f ./Dockerfile .

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean
