#!/bin/sh

go fmt main.go
go build -o sfagent main.go
./sfagent -ca-file=./safecomet_bundle.pem -token=s.9HoNJ15aBFvzrbFBKi74L32T -address=https://vault.safecomet.local  --roles=test1,test2
