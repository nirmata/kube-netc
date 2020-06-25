package cluster

import (
	"log"
	"os"
	"time"

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

type ObjectInfo struct {
	Name      string
	Kind      string
	Namespace string
	Node      string

	// Info from kubernetes labels
	LabelName      string
	LabelComponent string
	LabelInstance  string
	LabelVersion   string
	LabelPartOf    string
	LabelManagedBy string
}

type ClusterInfo struct {
	ObjectIPMap map[string]*ObjectInfo
}

func NewClusterInfo() *ClusterInfo {
	return &ClusterInfo{
		ObjectIPMap: make(map[string]*ObjectInfo),
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

	// Creating the informers for the different objects we want to track
	podInformer := factory.Core().V1().Pods().Informer()
	serviceInformer := factory.Core().V1().Services().Informer()
	nodeInformer := factory.Core().V1().Nodes().Informer()

	stopper := make(chan struct{})
	defer close(stopper)

	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleNewObject,
		UpdateFunc: c.handleUpdateObject,
		DeleteFunc: c.handleDeleteObject,
	})

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleNewObject,
		UpdateFunc: c.handleUpdateObject,
		DeleteFunc: c.handleDeleteObject,
	})

	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleNewObject,
		UpdateFunc: c.handleUpdateObject,
		DeleteFunc: c.handleDeleteObject,
	})

	go podInformer.Run(stopper)
	go serviceInformer.Run(stopper)
	nodeInformer.Run(stopper)
}
