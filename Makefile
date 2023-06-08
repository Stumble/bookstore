GO := go
NAME := bookstore
CGO_ENABLED = 0

POSTGRES_DOCKER_NAME=$(NAME)-postgres
POSTGRES_PASSWORD=my-secret
POSTGRES_DB=$(NAME)_test_db
POSTGRES_PORT=5432

REDIS_DOCKER_NAME=$(NAME)-redis
REDIS_PORT=6379

.PHONY: sqlc sqlc-verify
sqlc:
	cd pkg/repos && sqlc generate
sqlc-verify:
	cd pkg/repos && sqlc diff

docker-postgres-start:
	docker run -d --name $(POSTGRES_DOCKER_NAME) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -e POSTGRES_DB=$(POSTGRES_DB) -p $(POSTGRES_PORT):5432 postgres:14.5
	sleep 1
docker-postgres-stop:
	docker stop $(POSTGRES_DOCKER_NAME)
	docker rm $(POSTGRES_DOCKER_NAME)

docker-redis-start:
	docker run -d=true --name $(REDIS_DOCKER_NAME) -p $(REDIS_PORT):6379 redis
docker-redis-stop:
	docker stop $(REDIS_DOCKER_NAME)
	docker rm $(REDIS_DOCKER_NAME)

test-start-all: docker-postgres-start docker-redis-start
test-stop-all: docker-postgres-stop docker-redis-stop

test-cmd:
	export ENV=test && \
	export POSTGRES_APPNAME=bookstore && \
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -p 1 ./... -test.v

test-update-golden-cmd: test-start-all
	export ENV=test && \
	export POSTGRES_APPNAME=bookstore && \
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test -p 1 ./... -test.v -update
	make test-stop-all

test: test-start-all
	make test-cmd && make test-stop-all || (make test-stop-all; exit 2)

.PHONY: lint lint-fix
lint:
	@echo "--> Running linter"
	@golangci-lint run

lint-fix:
	@echo "--> Running linter auto fix"
	@golangci-lint run --fix
