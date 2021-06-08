package installer

import (
	"context"
	"fmt"
	"helm.sh/helm/v3/pkg/chartutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
	registryv1 "tkestack.io/tke/api/registry/v1"
	"tkestack.io/tke/cmd/tke-installer/app/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	applicationutil "tkestack.io/tke/pkg/application/util"
	registryConfig "tkestack.io/tke/pkg/registry/config"
	chartpath "tkestack.io/tke/pkg/registry/util/chartpath/v1"
	"tkestack.io/tke/pkg/util/log"
)

func newLoadTKE() (*TKE, error) {

	t := &TKE{
		namespace: namespace,
	}
	logOptions := log.NewOptions()
	log.Init(logOptions)
	t.log = log.WithName("tke-installer")

	err := t.loadTKEData()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("see what is cluster: %+v\n", t.Cluster.Status)
	err = t.initDataForDeployTKE()
	if err != nil {
		return nil, fmt.Errorf("config is not ready to do e2e test. %v", err)
	}
	return t, nil
}

func TestTKE_installApplication(t *testing.T) {

	tke, err := newLoadTKE()
	if err != nil {
		t.Fatal(err)
	}

	apps, err := tke.applicationClient.Apps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	name := "keydb"
	chart := &config.Chart{
		Name:            name,
		TenantID:        "default",
		ChartGroupName:  "public",
		Version:         "5.3.3",
		TargetCluster:   "global",
		TargetNamespace: "tcnp",
		Values: chartutil.Values{
			"test_value": []string{
				"value1",
				"value2",
			},
		},
	}

	err = tke.installApplication(context.Background(), config.PlatformApp{
		Name:   name,
		Enable: true,
		Chart:  *chart,
	}, apps.Items)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: test if application installation is successful
}

func TestTKE_uploadChart(t *testing.T) {

	tke, err := newLoadTKE()
	if err != nil {
		t.Fatal(err)
	}

	client := applicationutil.NewHelmClientWithoutRESTClient()
	conf := registryConfig.RepoConfiguration{
		Scheme:        "http",
		DomainSuffix:  tke.Para.Config.Registry.Domain(),
		Admin:         tke.Para.Config.Registry.Username(),
		AdminPassword: string(tke.Para.Config.Registry.Password()),
	}

	ct := registryv1.ChartGroup{
		Spec: registryv1.ChartGroupSpec{
			Name:        "public",
			TenantID:    "default",
			DisplayName: "public",
			Visibility:  "Public",
			Type:        "System",
		},
	}
	chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(conf, ct)
	if err != nil {
		t.Fatal(err)
	}

	f := os.Getenv("TEST_HELM_CHART_FILE")
	_, err = client.Push(&helmaction.PushOptions{
		ChartPathOptions: chartPathBasicOptions,
		ChartFile:        f,
		ForceUpload:      true,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTKE_getApplication(t *testing.T) {

	tke, err := newLoadTKE()
	if err != nil {
		t.Fatal(err)
	}

	name := "keydb"
	ns := "tcnp"

	apps, err := tke.applicationClient.Apps("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, app := range apps.Items {
		if app.Spec.Name == name && app.Namespace == ns {
			t.Logf("find %v/%v,%v", ns, name, app.Name)
		}
	}

}
