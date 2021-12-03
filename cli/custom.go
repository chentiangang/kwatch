package cli

import (
	"context"
	"fmt"
	msg "kwatch/msgdiff"

	reg "github.com/AlexsJones/go-type-registry/core"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *KubeWatch) SetItems() {
	k.Items = k.GetPods().Items
}

func (k *KubeWatch) GetItems() []v1.Pod {
	return k.GetPods().Items
}

func (k *KubeWatch) GetPods() *v1.PodList {
	pods, _ := k.clientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	return pods
}

func generateRegistry(r *reg.Registry) error {
	r.Put(&v1.Pod{})
	return nil
}

func (k *KubeWatch) DiffItems() {

	diff, equal := msg.PrettyDiff(k.Items, k.GetItems())
	d, _ := msg.PrettyDiff(k.GetItems(), k.Items)
	if !equal {
		for _, i := range d {
			fmt.Printf(i)
		}
		for _, i := range diff {
			fmt.Printf(i)
		}
	}
	fmt.Println("----------")
	fmt.Println()
	k.SetItems()
}
