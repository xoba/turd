#!/bin/bash -e
./clean.sh
git checkout lisp/gen.go

# generate all cxr's
go run . -m cxr

# generate eval and teval from template
go run . -m geneval

# compile defuns
go run . -m lispcompile

# test parsing
go run . -m lispparse

# test eval
go run . -m lisptest -debug
