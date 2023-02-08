package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	platform "tkestack.io/tke/api/platform"
	platformv1 "tkestack.io/tke/api/platform/v1"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/platform/types"
	"tkestack.io/tke/pkg/util/log"
)

func Mutate(reponseWriter http.ResponseWriter, request *http.Request) {
	var body []byte
	var err error
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

		if admissionReview.Request.Operation == v1.Update {
			v1OldCluster := platformv1.Cluster{}
			if err := json.Unmarshal(admissionReview.Request.OldObject.Raw, &v1OldCluster); err != nil {
				log.Errorf("Can't unmarshal old cluster, err: %v", err)
				http.Error(reponseWriter, fmt.Sprintf("Can't unmarshal old cluster, err: %v", err), http.StatusInternalServerError)
				return
			}
			oldCluster := platform.Cluster{}
			if err = platformv1.Convert_v1_Cluster_To_platform_Cluster(&v1OldCluster, &oldCluster, nil); err != nil {
				log.Errorf("Can't convert v1oldcluster to oldcluster, err: %v", err)
				http.Error(reponseWriter, fmt.Sprintf("Can't convert v1oldcluster to oldcluster, err: %v", err), http.StatusInternalServerError)
				return
			}
			admissionResponse = MutateClusterUpdate(&cluster, &oldCluster)
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

func MutateClusterUpdate(cluster *platform.Cluster, oldCluster *platform.Cluster) *v1.AdmissionResponse {
	typeCluster := types.Cluster{
		Cluster: cluster,
	}
	oldTypeCluster := types.Cluster{
		Cluster: oldCluster,
	}

	errorList := field.ErrorList{}
	p, err := clusterprovider.GetProvider(cluster.Spec.Type)
	if err != nil {
		errorList = append(errorList, field.NotFound(field.NewPath("spec").Child("type"), cluster.Spec.Type))
		return transferErrorList(&errorList, fmt.Sprintf("cluster %s update mutate failed: %v", oldCluster.Name, errorList.ToAggregate().Errors()))
	}
	jsonPatchByte, errorList := p.MutateUpdate(&typeCluster, &oldTypeCluster)
	if len(errorList) != 0 {
		return transferErrorList(&errorList, fmt.Sprintf("cluster %s update mutate failed: %v", oldCluster.Name, errorList.ToAggregate().Errors()))
	}
	if jsonPatchByte == nil {
		return &v1.AdmissionResponse{
			Allowed:   true,
		}
	}

	patchType := v1.PatchTypeJSONPatch
	return &v1.AdmissionResponse{
		Allowed:   true,
		Patch:     jsonPatchByte,
		PatchType: &patchType,
	}
}
