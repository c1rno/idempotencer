FROM golang:1.14.2-buster as builder

ENV ZCXXFLAGS "-Wc,lp64 -Wl,lp64"
ENV LD_LIBRARY_PATH "/usr/local/lib"
ENV LD_RUN_PATH "/usr/local/lib"

WORKDIR /go/src/idempotencer

COPY ./deployments/Makefile .

RUN make docker-install-dependencies
RUN make docker-install-rd-kafka
RUN make docker-install-libzmq
RUN make docker-install-czmq

COPY . .
RUN make vendor && make build

FROM c1rno/scratch:1.0
COPY --from=builder /go/src/idempotencer/idempotencer /bin/idempotencer
ENTRYPOINT ["/bin/idempotencer"]
