#!/bin/bash -e
./clean.sh
go run . -m compile
#go run github.com/xoba/turd/lisp/gen
go run . -m test
