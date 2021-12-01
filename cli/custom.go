package cli

import (
	"context"
	"fmt"
	"strings"

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
	d, _ := msg.PrettyDiff(k.GetItems(), k.Items)
	if !equal {
		for _, i := range d {
			switch  {
			case strings.HasPrefix(i,"removed:"):
				s := strings.Split(i,"=")[1]
				fmt.Println( s.(v1.Pod).ObjectMeta.Name)
			case strings.HasPrefix(i,"added:"):
			case strings.HasPrefix(i,"modified:"):
			}
			}
		}
		for _, i := range diff {
			fmt.Printf(i)
		}
		fmt.Println("----------")
		fmt.Println()
		k.SetItems()
	}
	//k.deepDiff()

}

func (k *KubeWatch) deepDiff() {
	d, equal := msg.DeepDiff(k.Items, k.GetItems())

	if !equal {
		for path, added := range d.Added {

			fmt.Printf("added: %s = %#v\n", path.String(), added)
		}
		for path, removed := range d.Removed {
			fmt.Printf("removed: %s = %#v\n", path.String(), removed)
		}
		for path, modified := range d.Modified {
			fmt.Printf("modified: %s = %#v\n", path.String(), modified)
		}

	}
}

//
//func reSetItem(podList []v1.Pod) (list []Pod) {
//	for _, i := range podList {
//		var l Pod
//		// objectMeta
//		l.Name = i.Name
//		l.Namespace = i.Namespace
//		l.UID = i.UID
//		l.Annotations = i.Annotations
//		l.Labels = i.Labels
//
//		// spec
//		l.Spec = i.Spec
//
//		// status
//		l.Status.ContainerStatuses = i.Status.ContainerStatuses
//		l.Status.HostIP = i.Status.HostIP
//		l.Status.Conditions = i.Status.Conditions
//		l.Status.PodIP = i.Status.PodIP
//		l.Status.StartTime = i.Status.StartTime
//		l.Status.Phase = i.Status.Phase
//		l.Status.QOSClass = i.Status.QOSClass
//		l.Status.PodIPs = i.Status.PodIPs
//
//		list = append(list, l)
//	}
//	return list
//}
