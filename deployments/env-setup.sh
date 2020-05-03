#!/bin/bash

set -euxo pipefail

mkdir /etc/idempotencer
cp /go/src/app/deployments/default-config.yaml /etc/idempotencer/

make vendor

tail -F /dev/null

