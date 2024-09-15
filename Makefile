.PHONY: build up run kill

build:
	@echo "Building Docker image..."
	podman build .

up:
	@echo "Running container..."
	podman-compose up

run: build up

kill:
	@echo "Stopping containers..."
	podman-compose down


