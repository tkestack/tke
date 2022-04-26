package storage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/api/logagent"
	"tkestack.io/tke/pkg/logagent/util"

	"tkestack.io/tke/pkg/util/log"
)

// LogfileProxyREST implements the REST endpoint.
type LogfileProxyREST struct {
	//rest.Storage
	store          *registry.Store
	platformClient platformversionedclient.PlatformV1Interface
}

// ConnectMethods returns the list of HTTP methods that can be proxied
func (r *LogfileProxyREST) ConnectMethods() []string {
	return []string{"GET"}
}

// NewConnectOptions returns versioned resource that represents proxy parameters
func (r *LogfileProxyREST) NewConnectOptions() (runtime.Object, bool, string) {
	return &logagent.LogFileProxyOptions{}, false, ""
}

// Connect returns a handler for the kube-apiserver proxy
func (r *LogfileProxyREST) Connect(ctx context.Context, clusterName string, opts runtime.Object, responder rest.Responder) (http.Handler, error) {
	clusterObject, err := r.store.Get(ctx, clusterName, &metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	agentConfig := clusterObject.(*logagent.LogAgent)
	proxyOpts := opts.(*logagent.LogFileProxyOptions)
	hostIP, err := util.GetClusterPodIP(ctx, agentConfig.Spec.ClusterName, proxyOpts.Namespace, proxyOpts.Pod, r.platformClient)
	if err != nil {
		return nil, errors.NewInternalError(fmt.Errorf("unable to get host ip with config %+v", *proxyOpts))
	}
	return &logFileProxyHandler{
		location:  &url.URL{Scheme: "http", Host: hostIP + ":" + util.LogagentPort},
		namespace: proxyOpts.Namespace,
		pod:       proxyOpts.Pod,
		container: proxyOpts.Container,
	}, nil
}

//
// New creates a new LogCollector proxy options object
func (r *LogfileProxyREST) New() runtime.Object {
	return &logagent.LogFileProxyOptions{}
}

type logFileProxyHandler struct {
	location  *url.URL
	namespace string
	pod       string
	container string
}

func (h *logFileProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	loc := *h.location
	loc.RawQuery = req.URL.RawQuery
	prefix := "/v1/logfile/download"
	// WithContext creates a shallow clone of the request with the new context.
	newReq := req.WithContext(context.Background())
	newReq.Header = netutil.CloneHeader(req.Header)
	loc.Path = prefix
	newReq.URL = &loc
	reverseProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: h.location.Scheme, Host: h.location.Host})
	reverseProxy.FlushInterval = 100 * time.Millisecond
	reverseProxy.ErrorLog = log.StdErrLogger()
	reverseProxy.ServeHTTP(w, newReq)
}
