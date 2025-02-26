ifeq ($(OS), Windows_NT)
	BIN_FILENAME := my-grpc-go-client.exe
else
	BIN_FILENAME := my-grpc-go-client
endif

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf bin

.PHONY: build
build: clean
	go build -o ./bin/$(BIN_FILENAME) ./cmd

.PHONY: execute
execute: clean build
	./bin/$(BIN_FILENAME)