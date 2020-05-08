package resourcelock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logagentv1client "tkestack.io/tke/api/client/clientset/versioned/typed/logagent/v1"
	v1 "tkestack.io/tke/api/logagent/v1"
)

// NotifyConfigMapLock defines the structure of using configmap resources to implement
// distributed locks.
type LogagentConfigMapLock struct {
	// ConfigMapMeta should contain a Name and a Namespace of a
	// ConfigMapMeta object that the LeaderElector will attempt to lead.
	ConfigMapMeta metav1.ObjectMeta
	Client        logagentv1client.ConfigMapsGetter
	LockConfig    Config
	cm            *v1.ConfigMap
}

// Get returns the election record from a ConfigMap Annotation
func (cml *LogagentConfigMapLock) Get() (*LeaderElectionRecord, error) {
	var record LeaderElectionRecord
	var err error
	cml.cm, err = cml.Client.ConfigMaps().Get(context.Background(), cml.ConfigMapMeta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if cml.cm.Annotations == nil {
		cml.cm.Annotations = make(map[string]string)
	}
	if recordBytes, found := cml.cm.Annotations[LeaderElectionRecordAnnotationKey]; found {
		if err := json.Unmarshal([]byte(recordBytes), &record); err != nil {
			return nil, err
		}
	}
	return &record, nil
}

// Create attempts to create a LeaderElectionRecord annotation
func (cml *LogagentConfigMapLock) Create(ler LeaderElectionRecord) error {
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}

	cml.cm, err = cml.Client.ConfigMaps().Create(context.Background(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cml.ConfigMapMeta.Name,
			Namespace: cml.ConfigMapMeta.Namespace,
			Annotations: map[string]string{
				LeaderElectionRecordAnnotationKey: string(recordBytes),
			},
		},
	}, metav1.CreateOptions{})
	return err
}

// Update will update an existing annotation on a given resource.
func (cml *LogagentConfigMapLock) Update(ler LeaderElectionRecord) error {
	if cml.cm == nil {
		return errors.New("endpoint not initialized, call get or create first")
	}
	recordBytes, err := json.Marshal(ler)
	if err != nil {
		return err
	}
	cml.cm.Annotations[LeaderElectionRecordAnnotationKey] = string(recordBytes)
	cml.cm, err = cml.Client.ConfigMaps().Update(context.Background(), cml.cm, metav1.UpdateOptions{})
	return err
}

// Describe is used to convert details on current resource lock
// into a string
func (cml *LogagentConfigMapLock) Describe() string {
	return fmt.Sprintf("%v/%v", cml.ConfigMapMeta.Namespace, cml.ConfigMapMeta.Name)
}

// Identity returns the Identity of the lock
func (cml *LogagentConfigMapLock) Identity() string {
	return cml.LockConfig.Identity
}
