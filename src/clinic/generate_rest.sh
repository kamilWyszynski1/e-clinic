#!/bin/bash
# shellcheck disable=SC2164
SCRIPT_PATH="$( cd "$(dirname "$0")" ; pwd -P )"
go build "$SCRIPT_PATH"/../tools/rest_rpc/cmd/main.go
./main -prefix=/api/v1 -interfaces=SpecjalistAssistant,PatientAssistant "$SCRIPT_PATH"/clinic.go
rm main
echo "clinic rest generated"