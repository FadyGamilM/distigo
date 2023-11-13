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
	$(GO_RUN) cmd/main.go -bolt_db_location="$(CURRENT_WORKING_DIR)\egypt.db" -http_addr=127.0.0.1:8080 -shard_name=egypt -shards_configs="$(CURRENT_WORKING_DIR)/shards.yaml" 

run_egypt_shard:	
	$(GO_RUN) cmd/main.go -bolt_db_location="$(CURRENT_WORKING_DIR)\egypt.db" -http_addr=127.0.0.1:8080 -shard_name=egypt -shards_configs="$(CURRENT_WORKING_DIR)/shards.yaml" 
run_usa_shard:
	$(GO_RUN) cmd/main.go -bolt_db_location="$(CURRENT_WORKING_DIR)\usa.db" -http_addr=127.0.0.1:8081 -shard_name=usa -shards_configs="$(CURRENT_WORKING_DIR)/shards.yaml" 
run_italy_shard:
	$(GO_RUN) cmd/main.go -bolt_db_location="$(CURRENT_WORKING_DIR)\italy.db" -http_addr=127.0.0.1:8082 -shard_name=italy -shards_configs="$(CURRENT_WORKING_DIR)/shards.yaml" 

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
