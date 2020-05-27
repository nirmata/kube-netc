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

type ClusterInfo struct {
	PodIPMap map[string]string
}

func handleNewPod(obj interface{}) {
	mObj, ok := obj.(*v1.Pod)

	if !ok {
		check(errors.New("Cannot treat obj as this type"))
	}

	fmt.Printf("[NEW] Pod %s added\n", mObj.GetName())
	fmt.Printf("\tLabels: %v\n", mObj.GetLabels())
	fmt.Printf("\tIP: %v\n", mObj.Status.PodIP)
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
		AddFunc: handleNewPod,
	})

	informer.Run(stopper)

}


/*
pods, err := clientset.CoreV1().Pods("mpich-system").List(metav1.ListOptions{})
if err != nil {
	// handle error
}
for _, pod := range pods.Items {
	fmt.Println(pod.Name, pod.Status.PodIP)
}
*/
