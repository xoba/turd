#!/bin/bash -e
./clean.sh
rm -rf lisp/gen
go run . -m compile
go run github.com/xoba/turd/lisp/gen
