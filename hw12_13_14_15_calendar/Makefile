BIN := "./bin/calendar"
BIN_SCHEDULER="./bin/calendar_scheduler"
BIN_SENDER="./bin/calendar_sender"
DOCKER_IMG="calendar:develop"
DOCKER_IMG="calendar:develop"
DOCKER_IMG_SH="scheduler:develop"
DOCKER_IMG_SENDER="sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

version: build
	$(BIN) version

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) --config ./configs/config.toml

run-scheduler: build
	$(BIN_SCHEDULER) --config ./configs/scheduler_config.toml

run-sender: build
	$(BIN_SENDER) --config ./configs/sender_config.toml

run-test-integration:
	go test -tags integration ./integration_tests/... -count 1 -v

generate:
	go generate ./...
		
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/calendar/Dockerfile .

build-img-sh:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SH) \
		-f build/scheduler/Dockerfile .

build-img-sender:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/sender/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

run-img-sh: build-img-sh
	docker run $(DOCKER_IMG_SH)

run-img-sender: build-img-sh
	docker run $(DOCKER_IMG_SENDER)	


test:
	go test -race ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run  --config=.golangci.yml ./...

up:
	docker compose -f ./deployments/docker-compose.yaml up -d 

down:
	docker compose -f deployments/docker-compose.yaml down

integration-tests:
	docker compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml up -d 
	docker compose -f deployments/docker-compose.yaml -f deployments/docker-compose.test.yaml down

.PHONY: build run run-scheduler run-sender build-img run-img version test lint up down integration-tests run-test-integration
