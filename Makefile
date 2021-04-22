Version := $(shell date "+%Y%m%d%H%M")
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := -s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)

.PHONY: run
run: build
	./bin/tech-share --debug

.PHONY: build
build: build-orm
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tech-share cmd/*.go

.PHONY: build-orm
build-orm:
	orm 'internal/service/model/*.yml'
	gofmt -s -w internal/service/model/*.go

.PHONY: dist
dist: build-orm
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/tech-share-linux cmd/*.go
