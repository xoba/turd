#!/bin/bash -e
./clean.sh
go run . -m lispcompile
go run . -m lisptest
