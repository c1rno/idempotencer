include deployments/Makefile

NAME=idempotencer
IMAGE=c1rno/$(NAME):latest
CMD=debug
CC=go

image:
	docker build \
		-t $(IMAGE) \
		-f deployments/Dockerfile \
		.

dev-env:
	docker-compose -f deployments/docker-compose.yaml up -d
	docker exec -it idempotencer bash

dev-down:
	docker-compose -f deployments/docker-compose.yaml down --remove-orphans

vendor:
	$(CC) mod tidy && $(CC) mod download

check:
	goimports -w -e .
	go vet ./... # TODO: replace by golangci-lint

test:
	$(CC) test -v ./...

build:
	CGO_LDFLAGS="-lzmq -lczmq -luuid -lpthread -lsodium -lrt -lstdc++ -lm -lc -lgcc" \
	GODEBUG=asyncpreemptoff=1 \
	$(CC) build -v -o $(NAME) \
	-mod=readonly \
	-ldflags '-extldflags "-static"' \
	-tags 'netgo std static_all' \
	./main.go
