#!/usr/bin/env bash
go generate ./... || exit
go fmt ./...  || exit
go vet ./... || exit
