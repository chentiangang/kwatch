package cli

import (
	"context"
	"flag"
	"path/filepath"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1 "k8s.io/api/apps/v1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeWatch struct {
	indexer    cache.Indexer
	queue      workqueue.RateLimitingInterface
	informer   cache.Controller
	clientSet  *kubernetes.Clientset
	Pods       []Pod
	Events     chan Events
	Deployment *appv1.Deployment
}

func NewClient() KubeWatch {
	var kubeconfig *string
	var master string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.StringVar(&master, "master", "", "master url")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags(master, *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	kubernetes.NewForConfig(config)

	// create the pod watcher
	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Pod than the version which was responsible for triggering the update.
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	//controller := NewController(queue, indexer, informer)

	return KubeWatch{
		clientSet: clientset,
		queue:     queue,
		indexer:   indexer,
		informer:  informer,
		Events:    make(chan Events, 1000),
	}
}

func (c *KubeWatch) GetDeployment() *appv1.Deployment {
	deploymentsClient := c.clientSet.AppsV1().Deployments(v1.NamespaceDefault)
	Deployment, _ := deploymentsClient.Get(context.TODO(), "nginx-deployment", meta_v1.GetOptions{})

	return Deployment
}
