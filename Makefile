#SHELL=/usr/bin/env bash

gomicro:
	go mod tidy \
	&& go build -ldflags "-s -w" -o client example/client/main.go \
	&& go build -ldflags "-s -w" -o server example/server/main.go

.PHONY: gomicro
BINS+=client
BINS+=server

docker:#env-GIT_USER env-GIT_PASSWORD
	docker build --tag gomicro-services -f ./Dockerfile .

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean
