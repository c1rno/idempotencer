#!/bin/bash

set -euxo pipefail

mkdir /etc/idempotencer
ln -s /go/src/app/deployments/default-config.yaml /etc/idempotencer/

# make vendor

tail -F /dev/null

