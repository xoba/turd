#!/bin/bash -e
go run . -m trans -profile profile.dat
#go tool pprof -top turd profile.dat
go tool pprof -web turd profile.dat
