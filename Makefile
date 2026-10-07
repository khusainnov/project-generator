BIN_DIR?=bin
BINARY?=project-gen
NAME?=
FLAGS?=

.PHONY: build gen test vet clean

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BINARY) .

gen: build
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is not set. Use 'make gen NAME=<project> [FLAGS=-go\ 1.24]'."; \
		exit 1; \
	fi
	./$(BIN_DIR)/$(BINARY) $(FLAGS) $(NAME)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf $(BIN_DIR)
