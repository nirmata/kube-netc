package cluster

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func (c *ClusterInfo) check(err error) {
	if err != nil {
		c.Logger.Fatalw(err.Error(),
			"package", "cluster",
		)
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
	mux         sync.Mutex
	Logger      *zap.SugaredLogger
	objectIPMap map[string]*ObjectInfo
}

func (ci *ClusterInfo) Set(ip string, o *ObjectInfo) {
	ci.mux.Lock()
	ci.objectIPMap[ip] = o
	ci.mux.Unlock()
}

func (ci *ClusterInfo) Get(ip string) (*ObjectInfo, bool) {
	ci.mux.Lock()
	defer ci.mux.Unlock()
	val, ok := ci.objectIPMap[ip]
	return val, ok
}

func NewClusterInfo(logger *zap.SugaredLogger) *ClusterInfo {
	logger.Debugw("starting cluster mapping",
		"package", "cluster",
	)
	return &ClusterInfo{
		objectIPMap: make(map[string]*ObjectInfo),
		Logger:      logger,
	}
}

func (c *ClusterInfo) Run() {

	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	c.check(err)

	clientset, err := kubernetes.NewForConfig(config)
	c.check(err)

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

	c.Logger.Debugw("informers starting",
		"package", "cluster",
	)
	go podInformer.Run(stopper)
	go serviceInformer.Run(stopper)
	nodeInformer.Run(stopper)
}
