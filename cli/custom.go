package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *KubeWatch) SetPods() {

	for _, i := range k.GetItems() {
		var pod Pod
		pod.UID = fmt.Sprintf("%s", i.UID)
		pod.Labels = make(map[string]string, 1)
		pod.Labels = i.Labels
		pod.IP = i.Status.HostIP
		pod.Namespace = i.Namespace
		pod.Name = i.Name

		for _, j := range i.Status.ContainerStatuses {
			var container Container
			container.ID = j.ContainerID
			//container.State = j.State.String()
			container.Name = j.Name
			pod.Containers = append(pod.Containers, container)
		}
		//pod.Status = i.Status.String()
		k.Pods = append(k.Pods, pod)
	}
}

func (k *KubeWatch) GetPods() (pods []Pod) {
	items := k.GetItems()
	for _, i := range items {
		var pod Pod
		pod.UID = fmt.Sprintf("%s", i.UID)
		pod.Labels = make(map[string]string, 1)
		pod.Labels = i.Labels
		pod.Namespace = i.Namespace
		pod.IP = i.Status.HostIP
		pod.Name = i.Name

		for _, j := range i.Status.ContainerStatuses {
			var container Container
			container.ID = j.ContainerID
			container.Name = j.Name
			pod.Containers = append(pod.Containers, container)
		}
		pods = append(pods, pod)
	}
	return pods
}

type Pod struct {
	Name       string            `json:"podName"`
	Namespace  string            `testdiff:"ignore" json:"-"`
	Labels     map[string]string `testdiff:"ignore" json:"-"`
	IP         string            `json:"ip"`
	UID        string            `testdiff:"ignore" json:"-"`
	Containers []Container       `json:"containers"`
}

type Data struct {
	Action string
	Key    string
	Old    interface{}
	New    interface{}
}

type Events struct {
	Event     string      `json:"event"`
	EventTime string      `json:"event_time"`
	Message   interface{} `json:"message"`
}

type Container struct {
	ID   string `json:"id"`
	Pod  string `json:"podName,omitempty"`
	Name string `json:"name"`
}

func (k *KubeWatch) GetItems() []v1.Pod {
	pods, _ := k.clientSet.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	return pods.Items
}

//

func (k KubeWatch) Diff() {

	var addpods []Pod
	for _, i := range k.AddedPod() {
		if i == "" {
			continue
		}
		for _, pod := range k.GetPods() {
			if i == pod.Name {
				addpods = append(addpods, pod)
			}
		}
	}

	if addpods != nil {
		bs, _ := json.Marshal(&addpods)
		k.Log <- Data{
			Action: "added",
			Key:    string(bs),
		}

		k.Events <- Events{
			Event:     "added",
			EventTime: time.Now().String(),
			Message:   addpods,
		}
	}

	var pods []Pod
	for _, i := range k.RemovedPod() {
		for _, pod := range k.Pods {
			if i == pod.Name {
				pods = append(pods, pod)
			}
		}
	}

	if pods != nil {
		bs, _ := json.Marshal(pods)
		k.Log <- Data{
			Action: "removed",
			Key:    string(bs),
		}
		k.Events <- Events{
			Event:     "removed",
			EventTime: time.Now().String(),
			Message:   pods,
		}
	}

	var removedc []Container
	for _, i := range k.RemovedContainer() {
		if i == "" {
			continue
		}
		for _, p := range k.Pods {
			for _, c := range p.Containers {
				if i == c.ID {
					var container Container
					container.ID = i
					container.Pod = p.Name
					container.Name = c.Name
					removedc = append(removedc, container)
				}
			}
		}
	}
	if removedc != nil {
		bs, _ := json.Marshal(removedc)
		k.Log <- Data{
			Action: "removed container",
			Key:    string(bs),
		}
		k.Events <- Events{
			Event:     "removed-container",
			EventTime: time.Now().String(),
			Message:   pods,
		}
	}

	var addedc []Container
	for _, i := range k.AddedContainer() {
		if i == "" {
			continue
		}
		for _, p := range k.GetPods() {
			for _, c := range p.Containers {
				if i == c.ID {
					var container Container
					container.ID = i
					container.Pod = p.Name
					container.Name = c.Name
					addedc = append(addedc, container)
				}
			}
		}
	}

	if addedc != nil {
		bs, _ := json.Marshal(addedc)
		k.Log <- Data{
			Action: "added container",
			Key:    string(bs),
		}
		k.Events <- Events{
			Event:     "added container",
			EventTime: time.Now().String(),
			Message:   addedc,
		}
	}

}

func (k KubeWatch) PodsJson() string {
	js, _ := json.Marshal(k.Pods)
	return string(js)
}

func (k KubeWatch) GetPodsJson() string {
	js, _ := json.Marshal(k.GetPods())
	return string(js)
}
