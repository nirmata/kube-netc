package cluster

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func check(err error) {
	if err != nil {
		log.Fatalf("[ERR] %s", err)
	}
}

// Message to be sent via MapUpdateChan to indicate to Collector that an update to the mapping has
// been made

type ClusterInfoChange struct {
	IP   string
	Info *PodInfo
}

type PodInfo struct {
	Name   string
	Labels map[string]string
}

type ClusterInfo struct {
	//IP->Name
	PodIPMap      map[string]*PodInfo
	MapUpdateChan chan ClusterInfoChange
}

func NewClusterInfo() *ClusterInfo {
	return &ClusterInfo{
		PodIPMap:      make(map[string]*PodInfo),
		MapUpdateChan: make(chan ClusterInfoChange),
	}
}

func (c *ClusterInfo) handleNewPod(obj interface{}) {
	mObj, ok := obj.(*v1.Pod)

	if !ok {
		check(errors.New("Cannot treat obj as v1.Pod"))
	}

	fmt.Printf("[NEW] Pod %s added\n", mObj.GetName())
	fmt.Printf("\tLabels: %v\n", mObj.GetLabels())
	fmt.Printf("\tIP: %v\n", mObj.Status.PodIP)

	ip := mObj.Status.PodIP
	name := mObj.GetName()
	labels := mObj.GetLabels()
	pinfo := &PodInfo{
		Name:   name,
		Labels: labels,
	}

	// Adding to the map
	c.PodIPMap[ip] = pinfo

	// Sending new update through channel to collector
	c.MapUpdateChan <- ClusterInfoChange{
		IP:   mObj.Status.PodIP,
		Info: pinfo,
	}
}

func (c *ClusterInfo) handleUpdatePod(oldObj interface{}, newObj interface{}) {
	mObj, mOk := newObj.(*v1.Pod)

	if !mOk {
		check(errors.New("Cannot treat obj as v1.Pod"))
	}

	_, oldOk := oldObj.(*v1.Pod)

	if !oldOk {
		check(errors.New("Cannot treat obj as v1.Pod"))
	}

	fmt.Printf("[UPDATE] Pod %s\n", mObj.GetName())
	fmt.Printf("\tLabels: %v\n", mObj.GetLabels())
	fmt.Printf("\tIP: %v\n", mObj.Status.PodIP)

	ip := mObj.Status.PodIP
	name := mObj.GetName()
	labels := mObj.GetLabels()
	pinfo := &PodInfo{
		Name:   name,
		Labels: labels,
	}

	// Adding to the map
	c.PodIPMap[ip] = pinfo

	// Sending new update through channel to collector
	c.MapUpdateChan <- ClusterInfoChange{
		IP:   mObj.Status.PodIP,
		Info: pinfo,
	}
}

func (c *ClusterInfo) handleDeletePod(obj interface{}) {
	mObj, ok := obj.(*v1.Pod)

	if !ok {
		check(errors.New("Cannot treat obj as v1.Pod"))
	}

	fmt.Printf("[DELETE] Pod %s deleted\n", mObj.GetName())
	fmt.Printf("\tLabels: %v\n", mObj.GetLabels())
	fmt.Printf("\tIP: %v\n", mObj.Status.PodIP)

	// Adding to the map
	c.PodIPMap[mObj.Status.PodIP] = nil

	// Sending new update through channel to collector
	c.MapUpdateChan <- ClusterInfoChange{
		IP:   mObj.Status.PodIP,
		Info: nil,
	}
}

func (c *ClusterInfo) Run() {

	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 5*time.Second)
	informer := factory.Core().V1().Pods().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleNewPod,
		UpdateFunc: c.handleUpdatePod,
		DeleteFunc: c.handleDeletePod,
	})

	informer.Run(stopper)

}
