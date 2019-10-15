# ddksm

POC of a process using kube-state-metrics as library with a custom `Builder` and `cache.Store`.
This process generate directly datadog metrics instead of exposing a Prometheus endpoint.

## build

```console
export GO111MODULE="on"
go build .
```

## Build docker image

```console
docker build -t cedriclamoriniere/ddksm:v0.0.1 .
```
