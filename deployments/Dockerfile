FROM golang:1.14.2-buster as builder

WORKDIR /go/src/idempotencer

COPY ./deployments/Makefile .

RUN make docker-install-dependencies
RUN make docker-install-rd-kafka
# RUN make docker-install-libzmq

# export LD_LIBRARY_PATH="/usr/local/lib"
# export LD_RUN_PATH="/usr/local/lib"

# RUN make docker-install-czmq

COPY . .
RUN make vendor && make build

FROM c1rno/scratch:1.0
COPY --from=builder /go/src/idempotencer/idempotencer /bin/idempotencer
ENTRYPOINT ["/bin/idempotencer"]
