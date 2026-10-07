BINARY?=project-gen
NAME?=
FLAGS?=

.PHONY: build gen test vet

build:
	go build -o $(BINARY) .

gen: build
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is not set. Use 'make gen NAME=<project> [FLAGS=-go\ 1.24]'."; \
		exit 1; \
	fi
	./$(BINARY) $(FLAGS) $(NAME)

test:
	go test ./...

vet:
	go vet ./...
