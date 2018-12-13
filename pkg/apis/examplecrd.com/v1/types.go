
package v1

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//CronTab Specification
type CronTab struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              CronTabSpec `json:"spec"`
}

//CronTabSpec Specification
type CronTabSpec struct {
	Replicas int32              `json:"replicas"`
	Template CronTabPodTemplate `json:"template"`
}

//CronTabPodTemplate Specification
type CronTabPodTemplate struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              apiv1.PodSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//CronTabList List
type CronTabList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CronTabList `json:"items"`
}
