#!/bin/sh

go fmt main.go
go build -o sfagent main.go
./sfagent -ca-file=./safecomet_bundle.pem -token=s.HQqjMn4kLUaDq8YG9imqvBnC -address=https://vault.safecomet.local  --roles=test1,test2
