APP ?= $(shell basename ${PWD})
TAG ?= $(shell [ -d .git ] && git log --pretty=format:"%cd.%h" --date=short -1)

GRPC_HOST  = $(shell dasel -f config.toml grpc.host)
GRPC_PORT  = $(shell dasel -f config.toml grpc.port)
REST_HOST  = $(shell dasel -f config.toml rest.host)
REST_PORT  = $(shell dasel -f config.toml rest.port)
DEBUG_HOST = $(shell dasel -f config.toml debug.host)
DEBUG_PORT = $(shell dasel -f config.toml debug.port)

.PHONY: tool
tool:
	@awk '$$1 == "_" { print $$2 | "xargs go install" }' ./tool/tool.go

.PHONY: grpc
grpc:
	@cd api/grpc && buf mod update
	@buf generate

.PHONY: build
build:
	@go build -ldflags "-X 'main.version=${TAG}'" -o ./tmp/${APP}

.PHONY: debug
debug:
	@go tool pprof $${OPTIONS:--seconds 10} \
		http://${DEBUG_HOST}:${DEBUG_PORT}/debug/pprof/$${PROFILE:-profile}

.PHONY: env
env: .git
	@echo "APP=${APP}" > .env
	@echo "TAG=${TAG}" >> .env
	@echo "GRPC_HOST=${GRPC_HOST}" >> .env
	@echo "GRPC_PORT=${GRPC_PORT}" >> .env
	@echo "REST_HOST=${REST_HOST}" >> .env
	@echo "REST_PORT=${REST_PORT}" >> .env
	@echo "DEBUG_HOST=${DEBUG_HOST}" >> .env
	@echo "DEBUG_PORT=${DEBUG_PORT}" >> .env

.PHONY: docker-build
docker-build: env
	@docker-compose build

.PHONY: docker-push
docker-push: env
	@docker-compose push

.PHONY: docker-up
docker-up: env
	@docker-compose up -d

.PHONY: docker-down
docker-down: env
	@docker-compose down --remove-orphans
