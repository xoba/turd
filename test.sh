#!/bin/bash -e
./clean.sh
git checkout lisp/gen.go
go run . -m eval # generates eval from template
go run . -m lispcompile # compiles defuns
go run . -m lispparse # tests parsing
go run . -m lisptest # tests eval
