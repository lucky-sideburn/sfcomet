#!/bin/sh

go fmt main.go
go build -o sfagent main.go
./sfagent -ca-file=./safecomet_bundle.pem -token=s.TpfHLe72M1DTThvmfcVRSVG5 -address=https://vault.safecomet.local 
