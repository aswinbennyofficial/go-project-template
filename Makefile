.PHONY: build run

build:
    podman-compose build

run: build
    podman-compose up


