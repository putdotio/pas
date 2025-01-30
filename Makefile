.PHONY: build up test lint

NAME := pas

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(NAME)

MAIN_COMPOSE := docker compose -f ./build/docker-compose.yml
build-docker:
	$(MAIN_COMPOSE) build pas
	mkdir -p ./bin
	docker create --name pas-builder pas
	docker cp pas-builder:/usr/local/bin/pas - > ./bin/pas
	docker rm -v pas-builder

MAIN_COMPOSE := docker compose -f ./build/docker-compose.yml
up:
	$(MAIN_COMPOSE) build pas
	$(MAIN_COMPOSE) rm -fsv
	$(MAIN_COMPOSE) up pas --build --abort-on-container-exit

TEST_COMPOSE := docker compose -f ./build/docker-compose-test.yml
test:
	$(TEST_COMPOSE) rm -fsv
	mkdir -p ./coverage
	$(TEST_COMPOSE) up pas-test --build --exit-code-from pas-test --abort-on-container-exit

lint:
	golangci-lint run
