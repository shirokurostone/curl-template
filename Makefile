.PHONY: build

VERSION:=$(shell git rev-parse --short HEAD)

build:
	go build -o ct -ldflags "-X main.version=${VERSION}"