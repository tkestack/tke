package helmchart

import (
	"strings"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	helmconfig "tkestack.io/tke/pkg/application/helm/config"
	"tkestack.io/tke/pkg/util/log"
)

func getCaps(restConfig *rest.Config) (*chartutil.Capabilities, error) {
	restClientGetter := &helmconfig.RESTClientGetter{RestConfig: restConfig}
	dc, err := restClientGetter.ToDiscoveryClient()
	if err != nil {
		return nil, errors.Wrap(err, "could not get Kubernetes discovery client")
	}
	// force a discovery cache invalidation to always fetch the latest server version/capabilities.
	dc.Invalidate()
	kubeVersion, err := dc.ServerVersion()
	if err != nil {
		return nil, errors.Wrap(err, "could not get server version from Kubernetes")
	}
	// Issue #6361:
	// Client-Go emits an error when an API service is registered but unimplemented.
	// We trap that error here and print a warning. But since the discovery client continues
	// building the API object, it is correctly populated with all valid APIs.
	// See https://github.com/kubernetes/kubernetes/issues/72051#issuecomment-521157642
	apiVersions, err := action.GetVersionSet(dc)
	if err != nil {
		if discovery.IsGroupDiscoveryFailedError(err) {
			log.Errorf("WARNING: The Kubernetes server has an orphaned API service. Server reports: %s", err)
			log.Errorf("WARNING: To fix this, kubectl delete apiservice <service-name>")
		} else {
			return nil, errors.Wrap(err, "could not get apiVersions from Kubernetes")
		}
	}

	return &chartutil.Capabilities{
		APIVersions: apiVersions,
		KubeVersion: chartutil.KubeVersion{
			Version: kubeVersion.GitVersion,
			Major:   kubeVersion.Major,
			Minor:   kubeVersion.Minor,
		},
		HelmVersion: chartutil.DefaultCapabilities.HelmVersion,
	}, nil
}

func Render(tkeYaml string, values map[string]interface{}, restConfig *rest.Config) (string, error) {
	c := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "tke",
			Version: "1.0.0",
		},
		Templates: []*chart.File{
			{Name: "templates/tke", Data: []byte(tkeYaml)},
		},
	}

	caps, err := getCaps(restConfig)
	if err != nil {
		return "", err
	}

	vals := map[string]interface{}{
		"Capabilities": caps,
		"Release": map[string]interface{}{
			// "Name":      options.Name,
			// "Namespace": options.Namespace,
			// "IsUpgrade": options.IsUpgrade,
			// "IsInstall": options.IsInstall,
			// "Revision":  options.Revision,
			"Service": "Helm",
		},
		"Values": values,
	}

	v, err := chartutil.CoalesceValues(c, vals)
	if err != nil {
		return "", err
	}
	out, err := engine.RenderWithClient(c, v, restConfig)
	// out, err := engine.Render(c, v)
	if err != nil {
		return "", err
	}

	return out["tke/templates/tke"], nil
}

func CoalesceValues(key, value string, values interface{}) {
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		values.(map[string]interface{})[keys[0]] = value
	} else {
		if values.(map[string]interface{})[keys[0]] == nil {
			values.(map[string]interface{})[keys[0]] = map[string]interface{}{}
		}
		CoalesceValues(strings.Join(keys[1:], "."), value, values.(map[string]interface{})[keys[0]])
	}
}
