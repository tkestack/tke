package util

import (
	appsv1 "k8s.io/api/apps/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

const reversion = "deployment.kubernetes.io/revision"

type ReplicaSets struct {
	replicaSet []appsv1.ReplicaSet
}

func NewReplicaSets(replicaSet []appsv1.ReplicaSet) ReplicaSets {
	return ReplicaSets{
		replicaSet: replicaSet,
	}
}

func (p ReplicaSets) Len() int {
	return len(p.replicaSet)
}

func (p ReplicaSets) Swap(i, j int) {
	p.replicaSet[i], p.replicaSet[j] = p.replicaSet[j], p.replicaSet[i]
}

func (p ReplicaSets) Less(i, j int) bool {
	ireVersion := MayNilMapGet(p.replicaSet[i].Annotations, reversion)
	jreVersion := MayNilMapGet(p.replicaSet[j].Annotations, reversion)
	return ireVersion > jreVersion
}

func (p ReplicaSets) GetReplicaSets() []appsv1.ReplicaSet {
	return p.replicaSet
}

type EXReplicaSets struct {
	replicaSet []extensionsv1beta1.ReplicaSet
}

func NewEXReplicaSets(replicaSet []extensionsv1beta1.ReplicaSet) EXReplicaSets {
	return EXReplicaSets{
		replicaSet: replicaSet,
	}
}

func (p EXReplicaSets) Len() int {
	return len(p.replicaSet)
}

func (p EXReplicaSets) Swap(i, j int) {
	p.replicaSet[i], p.replicaSet[j] = p.replicaSet[j], p.replicaSet[i]
}

func (p EXReplicaSets) Less(i, j int) bool {
	ireVersion := MayNilMapGet(p.replicaSet[i].Annotations, reversion)
	jreVersion := MayNilMapGet(p.replicaSet[j].Annotations, reversion)
	return ireVersion > jreVersion
}

func (p EXReplicaSets) GetReplicaSets() []extensionsv1beta1.ReplicaSet {
	return p.replicaSet
}

func MayNilMapGet(m map[string]string, k string) string {
	if m == nil {
		return ""
	}
	return m[k]
}
