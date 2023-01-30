GO := go
NAME := bookstore
CGO_ENABLED = 0

POSTGRES_DOCKER_NAME=$(NAME)-postgres
POSTGRES_PASSWORD=my-secret
POSTGRES_DB=$(NAME)_test_db
POSTGRES_PORT=5432

sqlc:
	cd pkg/repos && sqlc generate

docker-postgres-start:
	docker run -d --name $(POSTGRES_DOCKER_NAME) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) -p $(POSTGRES_PORT):5432 postgres:14.5

docker-postgres-stop:
	docker stop $(POSTGRES_DOCKER_NAME)
	docker rm $(POSTGRES_DOCKER_NAME)

test-start-all:
	make docker-postgres-start

test-stop-all:
	make docker-postgres-stop

test-cmd:
	export ENV=test && \
	export POSTGRES_APPNAME=bookstore && \
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -p 1 ./... -test.v

test-update-golden-cmd:
	export ENV=test && \
	export POSTGRES_APPNAME=bookstore && \
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -p 1 ./... -test.v -update


test: test-start-all
	make test-cmd && make test-stop-all || (make test-stop-all; exit 2)

.PHONY: lint lint-fix
lint:
	@echo "--> Running linter"
	@golangci-lint run

lint-fix:
	@echo "--> Running linter auto fix"
	@golangci-lint run --fix
