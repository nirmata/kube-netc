# kube-netc: A Kubernetes eBPF network monitor

[![Build Status](https://travis-ci.org/nirmata/kube-netc.svg?branch=master)](https://travis-ci.org/nirmata/kube-netc) [![Go Report Card](https://goreportcard.com/badge/github.com/nirmata/kube-netc)](https://goreportcard.com/report/github.com/nirmata/kube-netc)


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

## Demo

Test out the current state of kube-netc. Will guide you through building the container, running kube-netc and querying the prometheus /metrics endpoint.

First, build and run the docker container.

``` 
sudo make run 
```

Now with the container started we can query the server's /metrics endpoint and see the number of bytes sent per second through each connection.

```
curl http://$(sudo docker inspect -f "{{ .NetworkSettings.IPAddress }}" kube-netc-server):2112/metrics | grep bytes_sent_per_second
```

The id is in the form `src-dst`. The `src` should always be local and the `dst` could be or could not be local. Then the number appended at the end is the bytes being sent per second.

```
bytes_sent_per_second{id="127.0.0.1:34098-127.0.0.1:12379"} 0
bytes_sent_per_second{id="127.0.0.1:34100-127.0.0.1:12379"} 211
bytes_sent_per_second{id="127.0.0.1:34102-127.0.0.1:12379"} 0
bytes_sent_per_second{id="127.0.0.1:38008-127.0.0.1:16443"} 0
bytes_sent_per_second{id="127.0.0.1:38074-127.0.0.1:16443"} 0
bytes_sent_per_second{id="127.0.0.1:38168-127.0.0.1:16443"} 191
bytes_sent_per_second{id="127.0.0.1:38178-127.0.0.1:16443"} 0
```

__Note that [a fork of DataDog's ebpf library](https://github.com/drewrip/datadog-agent) is currently being used until compilation errors can be resolved.__
