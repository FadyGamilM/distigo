# Makefile
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
BINARY_NAME=myapp
BINARY_UNIX=$(BINARY_NAME)_unix
CURRENT_WORKING_DIR=$(shell cd)

all: build
build: 
	$(GO_BUILD) -o $(BINARY_NAME) -v
run:
	$(GO_RUN) cmd/main.go -bolt_db_location=$(CURRENT_WORKING_DIR)\my.db -http_addr=127.0.0.1:8080 -shard_name=myShard
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
