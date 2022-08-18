#!/usr/bin/env bash
echo "Building localdev splice router"
go build -gcflags "all=-N -l" -o "/tmp/localdev" "./cmd"
