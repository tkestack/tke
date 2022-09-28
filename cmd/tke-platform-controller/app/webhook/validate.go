package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platform "tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	"tkestack.io/tke/api/platform/validation"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/util/log"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func init() {
	_ = v1.AddToScheme(runtimeScheme)
}

func Validate(reponseWriter http.ResponseWriter, request *http.Request) {
	var body []byte
	var err error
	log.Infof("receive validate request, request: %v", request)
	if request.Body != nil {
		if body, err = ioutil.ReadAll(request.Body); err != nil {
			log.Errorf("request body read failed, err: %v", err)
			http.Error(reponseWriter, fmt.Sprintf("request body read failed, err: %v", err), http.StatusBadRequest)
			return
		}
		if len(body) == 0 {
			log.Errorf("request body length 0")
			http.Error(reponseWriter, "request body length 0", http.StatusBadRequest)
			return
		}
	} else {
		log.Errorf("request body nil")
		http.Error(reponseWriter, "request body nil", http.StatusBadRequest)
		return
	}

	contentType := request.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Errorf("Content-Type=%s, expect `application/json`", contentType)
		http.Error(reponseWriter, fmt.Sprintf("Content-Type=%s, expect `application/json`", contentType), http.StatusUnsupportedMediaType)
		return
	}

	admissionReview := v1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &admissionReview); err != nil {
		log.Errorf("decode request body to admission review failed, err: %v", err)
		http.Error(reponseWriter, fmt.Sprintf("decode request body to admission review failed, err: %v", err), http.StatusBadRequest)
		return
	}

	var admissionResponse *v1.AdmissionResponse
	switch admissionReview.Request.Kind.Kind {
	case "Cluster":
		v1Cluster := platformv1.Cluster{}
		if err := json.Unmarshal(admissionReview.Request.Object.Raw, &v1Cluster); err != nil {
			log.Errorf("Can't unmarshal cluster, err: %v", err)
			http.Error(reponseWriter, fmt.Sprintf("Can't unmarshal cluster, err: %v", err), http.StatusInternalServerError)
			return
		}

		cluster := platform.Cluster{}
		if err = platformv1.Convert_v1_Cluster_To_platform_Cluster(&v1Cluster, &cluster, nil); err != nil {
			log.Errorf("Can't convert v1cluster to cluster, err: %v", err)
			http.Error(reponseWriter, fmt.Sprintf("Can't convert v1cluster to cluster, err: %v", err), http.StatusInternalServerError)
			return
		}

		if admissionReview.Request.Operation == v1.Create {
			admissionResponse = ValidateCluster(&cluster)
		}
		if admissionReview.Request.Operation == v1.Update {
			oldCluster := platform.Cluster{}
			if err := json.Unmarshal(admissionReview.Request.Object.Raw, &oldCluster); err != nil {
				log.Errorf("Can't unmarshal cluster, err: %v", err)
				http.Error(reponseWriter, fmt.Sprintf("Can't unmarshal cluster, err: %v", err), http.StatusInternalServerError)
				return
			}
			admissionResponse = ValidateClusterUpdate(&cluster, &oldCluster)
		}
	default:
		log.Errorf("Can't recognized request kind %v", admissionReview.Request.Kind)
		http.Error(reponseWriter, fmt.Sprintf("Can't recognized request kind %v", admissionReview.Request.Kind), http.StatusBadRequest)
		return
	}

	admissionReview.Response = admissionResponse
	admissionReview.Response.UID = admissionReview.Request.UID

	admissionReviewBytes, err := json.Marshal(admissionReview)
	if err != nil {
		log.Errorf("Can't encode response: %v", err)
		http.Error(reponseWriter, fmt.Sprintf("Can't encode response: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err := reponseWriter.Write(admissionReviewBytes); err != nil {
		log.Errorf("Can't write response: %v", err)
		http.Error(reponseWriter, fmt.Sprintf("Can't write response: %v", err), http.StatusInternalServerError)
		return
	}
}

func ValidateCluster(cluster *platform.Cluster) *v1.AdmissionResponse {
	typeCluster := types.Cluster{
		Cluster: cluster,
	}
	errorList := validation.ValidateCluster(&typeCluster)
	if len(errorList) == 0 {
		return &v1.AdmissionResponse{
			Allowed: true,
		}
	}
	return transferErrorList(&errorList, fmt.Sprintf("cluster %s create validate failed: %v", cluster.Name, errorList.ToAggregate().Errors()))
}

func ValidateClusterUpdate(cluster *platform.Cluster, oldCluster *platform.Cluster) *v1.AdmissionResponse {
	typeCluster := types.Cluster{
		Cluster: cluster,
	}
	oldTypeCluster := types.Cluster{
		Cluster: oldCluster,
	}
	errorList := validation.ValidateClusterUpdate(&typeCluster, &oldTypeCluster)
	if len(errorList) == 0 {
		return &v1.AdmissionResponse{
			Allowed: true,
		}
	}
	return transferErrorList(&errorList, fmt.Sprintf("cluster %s update validate failed: %v", oldCluster.Name, errorList.ToAggregate().Errors()))
}

func transferErrorList(errorList *field.ErrorList, failedMessage string) *v1.AdmissionResponse {
	causes := make([]metav1.StatusCause, 0)
	for _, validateError := range *errorList {
		cause := metav1.StatusCause{
			Type:    metav1.CauseType(validateError.Type),
			Message: validateError.Detail,
			Field:   validateError.Field,
		}
		causes = append(causes, cause)
	}
	return &v1.AdmissionResponse{
		Allowed: false,
		Result: &metav1.Status{
			Code:    400,
			Message: failedMessage,
			Details: &metav1.StatusDetails{
				Causes: causes,
			},
		},
	}
}
