# kube-netc's design


## Packages

There are three main packages that make up kube-netc: [tracker](pkg/tracker/tracker.go), [collector](pkg/collector/collector.go) and [cluster](pkg/cluster/cluster.go).

### pkg/tracker

Tracker manages providing network statistics using eBPF. It also calculates the transfer rates for each connection. 

### pkg/collector

The collector registers the Prometheus metrics and handles the updates provided by the other packages.

### pkg/cluster

Any necessary information from the Kubernetes cluster is aggregated through cluster. It creates a mapping between between IPs and Kubernetes objects to provide more detailed metrics.


## Organization

<pre>
+---------+           +-----------+           +---------+                                                  
| tracker |---------->| collector |<----------| cluster |                                                 
+---------+           +-----------+           +---------+                                                      
                            |                                                                            
                            |
                            |  
                            v
                         +------+                                                                     
                         | main |                                                       
                         +------+
                            |
                            |
                            v
                       +----------+
                       | /metrics |
                       +----------+
</pre>
