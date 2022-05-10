#/bin/sh
go build -o=./bin ./...
open -a Terminal ./run/node1.sh
open -a Terminal ./run/node2.sh
open -a Terminal ./run/node3.sh
open -a Terminal ./run/proxy.sh
open -a Terminal ./run/source.sh