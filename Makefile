.PHONY: deps build clean run tidy help

deps:
	go mod download
	@mkdir -p output

build: deps
	docker-compose build

clean:
	docker-compose down -v
	rm -rf ./output/*

run: build
	docker-compose up

tidy:
	go mod tidy

help:
	@echo "Available commands:"
	@echo "  deps   - Install dependencies"
	@echo "  build  - Build the project (requires deps)"
	@echo "  clean  - Clean up containers and output directory"
	@echo "  run    - Run the project using docker-compose"
	@echo "  tidy   - Tidy up the go.mod file"
	@echo "  help   - Show this help message"

# Set the default target to help
.DEFAULT_GOAL := help