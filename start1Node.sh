#/bin/sh
go build -o=./bin ./...
open -a Terminal ./run/node1.sh
open -a Terminal ./run/proxy1Node.sh
open -a Terminal ./run/source.sh