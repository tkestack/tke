package types

import (
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/apis/audit"
	"strings"
	"tkestack.io/tke/pkg/util/log"
)

type Event struct {
	AuditID    string `json:"auditID"`
	Stage      string `json:"stage"`
	RequestURI string `json:"requestURI"`
	Verb       string `json:"verb"`
	UserName   string `json:"userName"`
	UserAgent  string `json:"userAgent"`

	Resource   string    `json:"resource"`
	Namespace  string    `json:"namespace"`
	Name       string    `json:"name"`
	UID        types.UID `json:"uid"`
	APIGroup   string    `json:"apiGroup"`
	APIVersion string    `json:"apiVersion"`
	SourceIPs  string    `json:"sourceIPs"`
	//ObjectRef *audit.ObjectReference

	Status  string `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	Details string `json:"details"`
	Code    int32  `json:"code"`
	//ResponseStatus *metav1.Status

	RequestObject  string `json:"requestObject"`
	ResponseObject string `json:"responseObject"`

	RequestReceivedTimestamp int64 `json:"requestReceivedTimestamp"`
	StageTimestamp           int64 `json:"stageTimestamp"`

	ClusterName string `json:"clusterName"`
}

func convertK8sEvent(event audit.Event) ([]*Event, error) {
	ev := Event{
		AuditID:                  string(event.AuditID),
		Stage:                    string(event.Stage),
		RequestURI:               event.RequestURI,
		Verb:                     event.Verb,
		UserName:                 event.User.Username,
		SourceIPs:                strings.Join(event.SourceIPs, ","),
		UserAgent:                event.UserAgent,
		RequestObject:            convertUnknown(event.RequestObject),
		ResponseObject:           convertUnknown(event.ResponseObject),
		RequestReceivedTimestamp: event.RequestReceivedTimestamp.Unix() * 1000,
		StageTimestamp:           event.StageTimestamp.Unix() * 1000,
		//Annotations : event.Annotations,
	}
	fillInObjectRef(&ev, event.ObjectRef)
	fillInStatus(&ev, event.ResponseStatus)
	if ev.Code >= 300 {
		// a Failure event
		if event.ResponseObject != nil {
			status := metav1.Status{}
			if err := json.Unmarshal(event.ResponseObject.Raw, &status); err == nil {
				if status.Status == "Failure" {
					fillInStatus(&ev, &status)
				}
			}
		}
	}
	if ev.Status == "" {
		if ev.Code >= 200 && ev.Code < 300 {
			ev.Status = "Success"
		} else {
			ev.Status = "Failure"
		}
	}
	if ev.Name == "" && ev.Verb == "create" &&
		event.ResponseObject != nil && event.ResponseStatus != nil &&
		event.ResponseStatus.Code >= 200 && event.ResponseStatus.Code < 300 {
		obj, _, err := unstructured.UnstructuredJSONScheme.Decode(event.ResponseObject.Raw, nil, nil)
		if err == nil {
			_, _, _, name, uid := extractMetadata(obj)
			if name != "" {
				ev.Name = name
				ev.UID = uid
			}
		}
	}
	return []*Event{&ev}, nil
}

func ConvertEvents(events []audit.Event) []*Event {
	var result []*Event
	for _, item := range events {
		evs, err := convertK8sEvent(item)
		if err != nil {
			log.Errorf("failed convert: %v", err)
		} else {
			result = append(result, evs...)
		}
	}
	return result
}

func extractMetadata(obj runtime.Object) (apiVersion, kind, namespace, name string, uid types.UID) {
	metaAccessor := meta.NewAccessor()
	name, _ = metaAccessor.Name(obj)
	namespace, _ = metaAccessor.Namespace(obj)
	kind, _ = metaAccessor.Kind(obj)
	uid, _ = metaAccessor.UID(obj)
	apiVersion, _ = metaAccessor.APIVersion(obj)
	return
}

func convertUnknown(obj *runtime.Unknown) string {
	if obj == nil {
		return ""
	}
	return string(obj.Raw)
}

func fillInObjectRef(event *Event, ref *audit.ObjectReference) {
	if ref == nil {
		return
	}
	event.Resource = ref.Resource
	event.Namespace = ref.Namespace
	event.Name = ref.Name
	event.UID = ref.UID
	event.APIGroup = ref.APIGroup
	event.APIVersion = ref.APIVersion
}

func fillInStatus(event *Event, status *metav1.Status) {
	if status == nil {
		return
	}
	event.Status = status.Status
	event.Message = status.Message
	event.Reason = string(status.Reason)
	event.Details = marshalAnything(status.Details)
	event.Code = status.Code
}

func marshalAnything(obj interface{}) string {
	res, _ := json.Marshal(obj)
	return string(res)
}
