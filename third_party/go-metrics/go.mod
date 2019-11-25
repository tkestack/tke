module github.com/docker/go-metrics

go 1.11

// Reduce the version dependency on the promethues golang client, because the
// kubernetes 1.16 version requires version 0.9.2, the interface is not compatible
require github.com/prometheus/client_golang v0.9.2
