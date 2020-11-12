#!/bin/bash -e
find . -name "*~" -exec rm \{} \;
find . -name "flymake_*.go" -exec rm \{} \;
