package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chentiangang/xlog"
	"github.com/d4l3k/messagediff"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	klog "k8s.io/klog/v2"
)

func (c *KubeWatch) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.syncToStdout(key.(string))

	//c.Diff()
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the pod to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *KubeWatch) syncToStdout(key string) error {
	_, exists, err := c.indexer.GetByKey(key)

	if err != nil {
		klog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Pod, so that we will see a delete for one pod
		//fmt.Printf("Pod %s does not exist anymore\n", key)
		c.Diff()
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a Pod was recreated with the same name
		//fmt.Printf("Sync/Add/Update for Pod %s\n", obj.(*v1.Pod).GetName())
		c.Diff()

	}
	//c.SetPods()
	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *KubeWatch) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	klog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

// Run begins watching and syncing.
func (c *KubeWatch) Run(workers int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	klog.Info("Starting Pod controller")

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Info("Stopping Pod controller")
}

func (c *KubeWatch) runWorker() {
	for c.processNextItem() {
	}
}

func (c *KubeWatch) Parse() {
	for i := range c.Events {
		c.Pods = c.GetPods()
		_, equal := messagediff.PrettyDiff(i.Spec, c.Spec)
		i.ConfigChanged = !equal
		if !equal {
			if c.IsRunning() {
				c.Spec = i.Spec
			}

		}

		bs, err := json.Marshal(&i)
		if err != nil {
			xlog.Error("%s", err)
		}
		xlog.Debug("%s", string(bs))
	}
}

//
func (c *KubeWatch) IsRunning() bool {
	for _, i := range c.GetPods() {
		for _, j := range i.Containers {
			if j.State == "nil" || j.State == "" {
				return false
			}
		}

	}
	return true
}
