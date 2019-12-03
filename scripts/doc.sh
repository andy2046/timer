#!/usr/bin/env bash

set -euo pipefail

godoc2md github.com/andy2046/timer \
 > $GOPATH/src/github.com/andy2046/timer/docs.md
