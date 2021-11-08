package storage

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	platforminternalclient "tkestack.io/tke/api/client/clientset/internalversion/typed/platform/internalversion"
	"tkestack.io/tke/api/platform"
	"tkestack.io/tke/pkg/apiserver/authentication"
	"tkestack.io/tke/pkg/platform/apiserver/filter"
	clusterprovider "tkestack.io/tke/pkg/platform/provider/cluster"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/pkiutil"
)

type CustomResourceHandler struct {
	LoopbackClientConfig *rest.Config
}

var pool clientX509Pool

type clientX509Pool struct {
	sm sync.Map
}

type clientX509Cache struct {
	clientCertData []byte
	clientKeyData  []byte
}

var (
	// Scheme is the default instance of runtime.Scheme to which types in the TKE API are already registered.
	Scheme = runtime.NewScheme()
	// Codecs provides access to encoding and decoding for the scheme
	Codecs = serializer.NewCodecFactory(Scheme)

	unversionedVersion = schema.GroupVersion{Group: "", Version: "v1"}
	unversionedTypes   = []runtime.Object{
		&metav1.Status{},
	}
)

func init() {
	Scheme.AddUnversionedTypes(unversionedVersion, unversionedTypes...)
}

// ServeHTTP is a proxy for unregister custom resource
func (n *CustomResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestInfo, ok := apirequest.RequestInfoFrom(r.Context())
	if !ok {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("no RequestInfo found in the context")),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	// crd resource start with /apis
	if requestInfo.APIPrefix != "apis" {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(fmt.Errorf("request crd is validate")),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	platformClient := platforminternalclient.NewForConfigOrDie(n.LoopbackClientConfig)
	config, err := getConfig(r.Context(), platformClient)

	if err != nil {
		responsewriters.ErrorNegotiated(
			apierrors.NewInternalError(err),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	TLSClientConfig := &tls.Config{}
	TLSClientConfig.InsecureSkipVerify = true
	if config.TLSClientConfig.CertData != nil && config.TLSClientConfig.KeyData != nil {
		cert, err := tls.X509KeyPair(nil, config.TLSClientConfig.KeyData)
		if err != nil {
			responsewriters.ErrorNegotiated(
				apierrors.NewGenericServerResponse(http.StatusUnauthorized, requestInfo.Verb, schema.GroupResource{
					Group:    requestInfo.APIGroup,
					Resource: requestInfo.Resource,
				}, requestInfo.Name,
					err.Error(), 0, true),
				Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
			)
			return
		}
		TLSClientConfig.Certificates = []tls.Certificate{cert}
	} else if config.BearerToken == "" {
		responsewriters.ErrorNegotiated(
			apierrors.NewGenericServerResponse(http.StatusUnauthorized, requestInfo.Verb, schema.GroupResource{
				Group:    requestInfo.APIGroup,
				Resource: requestInfo.Resource,
			}, requestInfo.Name,
				fmt.Sprintf("%s has NO BearerToken", filter.ClusterFrom(r.Context())),
				0, true),
			Codecs, schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}, w, r,
		)
		return
	}

	reserveProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   config.Host,
	})
	reserveProxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       TLSClientConfig,
	}

	reserveProxy.Director = buildDirector(r, config)
	reserveProxy.FlushInterval = 100 * time.Millisecond
	reserveProxy.ServeHTTP(w, r)
}

func buildDirector(r *http.Request, config *rest.Config) func(req *http.Request) {
	clusterName := filter.ClusterFrom(r.Context())
	userName, tenantID := authentication.UsernameAndTenantID(r.Context())
	return func(req *http.Request) {
		req.Header.Set(filter.ClusterNameHeaderKey, clusterName)
		req.Header.Set("X-Remote-User", userName)
		req.Header.Set("X-Remote-Extra-TenantID", tenantID)
		if config.BearerToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.BearerToken))
		}
		req.URL = &url.URL{
			Scheme: "https",
			Host:   config.Host,
			Path:   r.RequestURI,
		}
	}
}

func getConfig(ctx context.Context, platformClient platforminternalclient.PlatformInterface) (*rest.Config, error) {

	clusterName := filter.ClusterFrom(ctx)
	if clusterName == "" {
		return nil, errors.NewBadRequest("clusterName is required")
	}

	cluster, err := platformClient.Clusters().Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if cluster.Status.Locked != nil && *cluster.Status.Locked {
		return nil, fmt.Errorf("cluster %s has been locked", cluster.ObjectMeta.Name)
	}

	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if len(tenantID) > 0 && cluster.Spec.TenantID != tenantID {
		return nil, errors.NewNotFound(platform.Resource("clusters"), cluster.ObjectMeta.Name)
	}

	clusterWrapper, err := clusterprovider.GetCluster(ctx, platformClient, cluster, username)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{}
	if cluster.AuthzWebhookEnabled() {
		clientCertData, clientKeyData, err := getOrCreateClientCert(ctx, clusterWrapper.ClusterCredential)
		if err != nil {
			return nil, err
		}
		config, err = clusterWrapper.RESTConfigForClientX509(config, clientCertData, clientKeyData)
		if err != nil {
			return nil, err
		}
	} else {
		config, err = clusterWrapper.RESTConfig(config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func getOrCreateClientCert(ctx context.Context, credential *platform.ClusterCredential) ([]byte, []byte, error) {
	groups := authentication.Groups(ctx)
	username, tenantID := authentication.UsernameAndTenantID(ctx)
	if tenantID != "" {
		groups = append(groups, fmt.Sprintf("tenant:%s", tenantID))
	}

	ns, ok := request.NamespaceFrom(ctx)
	if ok {
		groups = append(groups, fmt.Sprintf("namespace:%s", ns))
	}

	cache, ok := pool.sm.Load(makeClientKey(username, groups))
	if ok {
		return cache.(*clientX509Cache).clientCertData, cache.(*clientX509Cache).clientKeyData, nil
	}

	clientCertData, clientKeyData, err := pkiutil.GenerateClientCertAndKey(username, groups, credential.CACert,
		credential.CAKey)
	if err != nil {
		return nil, nil, err
	}

	pool.sm.Store(makeClientKey(username, groups), &clientX509Cache{clientCertData: clientCertData,
		clientKeyData: clientKeyData})

	log.Debugf("generateClientCert success. username:%s groups:%v\n clientCertData:\n %s clientKeyData:\n %s",
		username, groups, clientCertData, clientKeyData)

	return clientCertData, clientKeyData, nil
}

func makeClientKey(username string, groups []string) string {
	return fmt.Sprintf("%s###%v", username, groups)
}
