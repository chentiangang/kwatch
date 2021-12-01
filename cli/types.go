package cli

//
//import (
//	"k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/types"
//)
//
//type PodList struct {
//	Items []Pod `json:"items" protobuf:"bytes,2,rep,name=items"`
//}
//
//type Pod struct {
//	ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
//	Spec       v1.PodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
//	Status     v1.PodStatus  `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
//}
//
//type ObjectMeta struct {
//	Name        string            `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
//	Namespace   string            `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
//	UID         types.UID         `json:"uid,omitempty" protobuf:"bytes,5,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`
//	Labels      map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`
//	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
//}
//
//type PodStatus struct {
//	Phase             v1.PodPhase          `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
//	Conditions        []v1.PodCondition    `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
//	HostIP            string               `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
//	PodIP             string               `json:"podIP,omitempty" protobuf:"bytes,6,opt,name=podIP"`
//	PodIPs            []v1.PodIP           `json:"podIPs,omitempty" protobuf:"bytes,12,rep,name=podIPs" patchStrategy:"merge" patchMergeKey:"ip"`
//	StartTime         *metav1.Time         `json:"startTime,omitempty" protobuf:"bytes,7,opt,name=startTime"`
//	ContainerStatuses []v1.ContainerStatus `json:"containerStatuses,omitempty" protobuf:"bytes,8,rep,name=containerStatuses"`
//	QOSClass          v1.PodQOSClass       `json:"qosClass,omitempty" protobuf:"bytes,9,rep,name=qosClass"`
//}
