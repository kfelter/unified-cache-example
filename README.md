
# Unified Cache Example
Unify N number of pull through caches with a ketama proxy

# Ketama hashing in the proxy
Consistent hashing via the ketama package is used by the proxy to choose the cache node that will pull and cache the content.

![](docs/diagram.png?raw=true)

# Code
[client](cmd/client/main.go)

[proxy](cmd/proxy/main.go)

[node](cmd/node/main.go)

[source](cmd/source/main.go)
