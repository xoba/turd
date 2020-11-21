#!/bin/bash -e
./clean.sh
git checkout lisp/gen.go
go run . -m lispcompile
go run . -m lisptest
