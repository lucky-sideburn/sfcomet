#!/bin/sh

go fmt main.go
go build -o sfagent main.go
./sfagent
