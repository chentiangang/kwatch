package cli

import (
	"context"
	"fmt"
	msg "kwatch/msgdiff"

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

func (k *KubeWatch) DiffItems() {

	diff, equal := msg.PrettyDiff(k.Items, k.GetItems())

	if !equal {
		for _, i := range diff {
			Output(i)
		}
	}
	fmt.Println("----------")
	fmt.Println()
	k.SetItems()
}

func Output(i map[string]interface{}) {
	for k, v := range i {
		switch k {
		case "removed":
			fmt.Printf("\tremoved pod: %s\n", v.(v1.Pod).Name)
		case "added":
			fmt.Printf("\tadded pod: %s\n", v.(v1.Pod).Name)
		case "modified":
			fmt.Printf("\tmodified : %+v\n", v)
		}
	}
}
