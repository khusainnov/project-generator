PROJECT_NAME?=project-gen

build:
	go build -o $(PROJECT_NAME) main.go

gen:
	@if [ -z "$(PROJECT_NAME)" ]; then \
		echo "Error: PROJECT_NAME is not set. Use 'make gen PROJECT_NAME=<name> <arg>'."; \
		exit 1; \
	fi
	@if [ -z "$(filter gen,$(MAKECMDGOALS))" ] || [ "$(word 2, $(MAKECMDGOALS))" = "" ]; then \
		echo "Error: Additional argument is required. Use 'make gen <arg>'."; \
		exit 1; \
	fi
	$(eval EXTRA_ARG := $(word 2, $(MAKECMDGOALS)))
	@echo "Running: ./$(PROJECT_NAME) $(EXTRA_ARG)"
	./$(PROJECT_NAME) $(EXTRA_ARG)