#/bin/sh
./bin/client
sleep 1
echo "===METRICS==="
curl -s localhost:3000/_/metrics | jq .