# kube-netc
kube-netc (pronounced <i>kube-net-see</i>) is a Kubernetes network monitor built using eBPF

## Examples

Try the current working sample test programs:

#### recv.go: writing out the number of current active connections

```
make recv
sudo ./recv
```

#### promserv.go: exposing the number of current active connections to prometheus

```
make promserv
sudo ./promserv
```
The exposed metric is `active_connections`.

## Collector

[The collector](collector/) package is the prometheus collector. Currently it only exposes the trackers activeConnections for testing purposes and has been demonstrated to work with Grafana.

## Tracker

[The tracker](tracker/) package is being used to interface between DataDog's ebpf library and our collector. The tracker converts the connection data into a format usable to the collector and prepares to be read. Much work still needs to be done to design out this package to get the right information, and prepare it to be used when needed by the collector.

__Note that [a fork of DataDog's ebpf library](https://github.com/drewrip/datadog-agent) is currently being used until compilation errors can be resolved.__
