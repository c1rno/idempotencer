#!/bin/bash

set -euxo pipefail

go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &
go run . downstream 2>> tmp.consume &

wait $(jobs -p)
