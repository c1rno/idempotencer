include deployments/Makefile

NAME=idempotencer
IMAGE=c1rno/$(NAME):latest
BASE_CMD=bash
CC=go

image:
	docker build \
		-t $(IMAGE) \
		-f deployments/Dockerfile \
		.

dev-shell: image
	docker run -it --rm \
		--name $(NAME) \
		$(IMAGE) $(BASE_CMD)

vendor:
	$(CC) mod tidy && $(CC) mod download

test:
	$(CC) test -v ./...

build:
	# CGO_LDFLAGS="-lzmq -lpthread -lsodium -lrt -lstdc++ -lm -lc -lgcc" \
	$(CC) build -v -o $(NAME) \
	-ldflags '-extldflags "-static"' \
	-tags 'netgo std static_all' \
	./main.go
