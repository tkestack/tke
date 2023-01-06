package util

import corev1 "k8s.io/api/core/v1"

type Pods struct {
	pods []corev1.Pod
}

func NewPods(pods []corev1.Pod) Pods {
	return Pods{
		pods: pods,
	}
}

func (p Pods) Len() int {
	return len(p.pods)
}

func (p Pods) Swap(i, j int) {
	p.pods[i], p.pods[j] = p.pods[j], p.pods[i]
}

func (p Pods) Less(i, j int) bool {
	// created at the same time, sort by the first letter of the pod name
	if p.pods[i].CreationTimestamp.Time == p.pods[j].CreationTimestamp.Time {
		return p.pods[i].Name > p.pods[j].Name
	}
	// the earliest created time comes first
	return p.pods[j].CreationTimestamp.After(p.pods[i].CreationTimestamp.Time)
}

func (p Pods) GetPods() []corev1.Pod {
	return p.pods
}
