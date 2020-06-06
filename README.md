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

To test the current capabilities of **kube-netc**, this guide will walk you through viewing the network statistics of your nodes.

First, install the daemon set using the install.yaml:

``` 
kubectl apply -f config/install.yaml
```

This will start the **kube-netc** DaemonSet on your cluster and setup the required roles. Then, we get the name of the kube-netc pod:

```
kubectl get pods | grep kube-netc
```

For example, my **kube-netc** pod is:

```
kube-netc-j56cx
```
In a new terminal, we port-forward the port of our pod so we can access it with *curl* outside the cluster with:

```
kubectl port-forward kube-netc-j56cx 2112:2112
```

2112 is the port we are going to access the Prometheus endpoint on. We can then curl the /metrics endpoint using curl on local host to show the Prometheus metrics:

```
curl localhost:2112/metrics | grep bytes
```

This is an example output of the query showing the total bytes received by this node from each given connection:

```
...
# HELP bytes_recv Total bytes received from a given connection
# TYPE bytes_recv gauge
bytes_recv{pod_address="10.244.0.1",pod_name="NOT_FOUND"} 117
bytes_recv{pod_address="10.244.0.2",pod_name="kube-netc-j56cx"} 420
bytes_recv{pod_address="10.244.0.3",pod_name="local-path-provisioner-bd4bb6b75-fjpmj"} 36494
bytes_recv{pod_address="10.244.0.4",pod_name="coredns-66bff467f8-csjhc"} 2346
bytes_recv{pod_address="10.244.0.5",pod_name="coredns-66bff467f8-69jxm"} 2346
bytes_recv{pod_address="10.96.0.1",pod_name="NOT_FOUND"} 71754
bytes_recv{pod_address="104.18.123.25",pod_name="NOT_FOUND"} 0
bytes_recv{pod_address="127.0.0.1",pod_name="NOT_FOUND"} 290
bytes_recv{pod_address="172.18.0.1",pod_name="NOT_FOUND"} 1672
bytes_recv{pod_address="172.18.0.2",pod_name="kindnet-jpbtd"} 168608
bytes_recv{pod_address="172.18.0.2",pod_name="kube-apiserver-drewcluster-control-plane"} 56160
bytes_recv{pod_address="172.18.0.2",pod_name="kube-controller-manager-drewcluster-control-plane"} 43609
bytes_recv{pod_address="172.18.0.2",pod_name="kube-proxy-s57cn"} 157359
bytes_recv{pod_address="192.168.1.1",pod_name="NOT_FOUND"} 132
bytes_recv{pod_address="3.223.220.229",pod_name="NOT_FOUND"} 0
bytes_recv{pod_address="34.195.246.183",pod_name="NOT_FOUND"} 0
bytes_recv{pod_address="52.1.121.53",pod_name="NOT_FOUND"} 0
bytes_recv{pod_address="52.5.11.128",pod_name="NOT_FOUND"} 0
bytes_recv{pod_address="52.72.232.213",pod_name="NOT_FOUND"} 0
bytes_recv{pod_address="::1",pod_name="NOT_FOUND"} 1939
...
```

As we see the bytes received by each connection is shown and the source IP is given. If there is a known pod with the same IP the *pod_name* is also given, if not it defaults to "NOT FOUND".
