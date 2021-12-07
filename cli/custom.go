package cli

import (
	"context"
	"fmt"
	msg "kwatch/msgdiff"
	"strings"

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
			k.Print(i)
		}
	}
	fmt.Println("----------")
	fmt.Println()
	k.SetItems()
}

func (k *KubeWatch) Print(i map[string]interface{}) {
	for key, value := range i {
		switch key {
		case "removed":
			fmt.Printf("\tremoved pod: %s\n", value.(v1.Pod).Name)
		case "added":
			fmt.Printf("\tadded pod: %s\n", value.(v1.Pod).Name)
		case "modified":
			vKey, vValue := GetModifiedKeyValue(value)
			d, _ := msg.PrettyDiff(k.GetItems(), k.Items)
			for _, j := range d {
				for Key, Value := range j {
					if Key == "modified" {
						dKey, dValue := GetModifiedKeyValue(Value)
						if dKey == vKey {
							fmt.Printf("\tmodified: %s %s ==> %s\n", vKey, dValue, vValue)
						}
					}
				}
			}
		}
	}
}

func TrimFunc(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return r == '\n' || r == ' '
	})
}

func GetModifiedKeyValue(s interface{}) (key, value string) {
	vSplit := strings.Split(s.(string), "=")
	return TrimFunc(vSplit[0]), TrimFunc(vSplit[1])
}
