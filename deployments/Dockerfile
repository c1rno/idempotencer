FROM golang:1.14.2-buster as builder

ENV ZCXXFLAGS "-Wc,lp64 -Wl,lp64"
ENV LD_LIBRARY_PATH "/usr/local/lib"
ENV LD_RUN_PATH "/usr/local/lib"

WORKDIR /go/src/idempotencer

COPY ./deployments/Makefile ./deployments/Makefile

RUN make -C deployments docker-install-dependencies
RUN make -C deployments docker-install-rd-kafka
RUN make -C deployments docker-install-libzmq
RUN make -C deployments docker-install-czmq

COPY ./Makefile ./go.mod ./go.sum ./
RUN make vendor

COPY . .
RUN make build

FROM c1rno/scratch:1.0

COPY --from=builder /go/src/idempotencer/idempotencer /bin/
COPY --from=builder /go/src/idempotencer/deployments/default-config.yaml /etc/idempotencer/

EXPOSE 8888

ENTRYPOINT ["/bin/idempotencer"]
