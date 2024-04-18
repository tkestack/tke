package helmchart

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	applicationv1 "tkestack.io/tke/api/application/v1"
	platformversionedclient "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	"tkestack.io/tke/pkg/application/config"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	helmconfig "tkestack.io/tke/pkg/application/helm/config"
	helmutil "tkestack.io/tke/pkg/application/helm/util"
	"tkestack.io/tke/pkg/application/util"
	registryutil "tkestack.io/tke/pkg/registry/util"
	"tkestack.io/tke/pkg/util/file"
	"tkestack.io/tke/pkg/util/log"
)

type Helper struct {
	installOptions *helmaction.InstallOptions
	actionInstall  *action.Install
	actionClient   *helmaction.Client
	actionConfig   *action.Configuration
	app            *applicationv1.App
	cfg            *action.Configuration
}

func NewHelper(ctx context.Context, platformClient platformversionedclient.PlatformV1Interface, app *applicationv1.App, repo config.RepoConfiguration) (*Helper, error) {
	values, err := helmutil.MergeValues(app.Spec.Values.Values, app.Spec.Values.RawValues, string(app.Spec.Values.RawValuesType))
	if err != nil {
		return nil, errors.Wrap(err, "helper MergeValues")
	}
	cfg, err := util.NewActionConfigWithProvider(ctx, platformClient, app)
	if err != nil {
		log.Errorf("failed to new action config, err:%s", err.Error())
		return nil, errors.Wrap(err, "failed to get Helm action configuration")
	}
	restConfig, err := cfg.RESTClientGetter.ToRESTConfig()
	if err != nil {
		log.Errorf("failed to new action config, err:%s", err.Error())
		return nil, errors.Wrap(err, "failed to get Helm action configuration")
	}

	const defaultTimeout = 600 * time.Second
	var clientTimeout = defaultTimeout
	if app.Spec.Chart.InstallPara.Timeout > 0 {
		clientTimeout = app.Spec.Chart.InstallPara.Timeout
	}

	restClientGetter := &helmconfig.RESTClientGetter{RestConfig: restConfig}
	// we should set namespace here. If not, release will be installed in target namespace, but resources will not be installed in target namespace
	restClientGetter.Namespace = &app.Spec.TargetNamespace
	client := helmaction.NewClient("", restClientGetter)
	// chartPathBasicOptions, err := chartpath.BuildChartPathBasicOptions(repo, app.Spec.Chart)
	chartPathBasicOptions, err := buildChartPathBasicOptions(repo, app.Spec.Chart)
	if err != nil {
		return nil, errors.Wrap(err, "helper BuildChartPathBasicOptions")
	}
	destfile, err := client.Pull(&helmaction.PullOptions{
		ChartPathOptions: chartPathBasicOptions,
	})
	chartPathBasicOptions.ExistedFile = destfile
	options := &helmaction.InstallOptions{
		DryRun:           false,
		Namespace:        app.Spec.TargetNamespace,
		ReleaseName:      app.Spec.Name,
		DependencyUpdate: true,
		Values:           values,
		Timeout:          clientTimeout,
		ChartPathOptions: chartPathBasicOptions,
		CreateNamespace:  app.Spec.Chart.InstallPara.CreateNamespace,
		Atomic:           app.Spec.Chart.InstallPara.Atomic,
		Wait:             app.Spec.Chart.InstallPara.Wait,
		WaitForJobs:      app.Spec.Chart.InstallPara.WaitForJobs,
	}

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(restClientGetter, options.Namespace, "", log.Debugf)
	if err != nil {
		return nil, errors.Wrap(err, "helper actionConfig.Init")
	}
	actionConfig.RegistryClient, err = registry.NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "helper registry.NewClient")
	}
	installCli := action.NewInstall(actionConfig)
	installCli.DryRun = options.DryRun
	installCli.DependencyUpdate = options.DependencyUpdate
	installCli.Timeout = options.Timeout
	installCli.Namespace = options.Namespace
	installCli.ReleaseName = options.ReleaseName
	installCli.Description = options.Description
	installCli.IsUpgrade = options.IsUpgrade
	installCli.Atomic = options.Atomic
	installCli.CreateNamespace = options.CreateNamespace
	installCli.Wait = options.Wait
	installCli.WaitForJobs = options.WaitForJobs

	options.ChartPathOptions.ApplyTo(&installCli.ChartPathOptions)
	return &Helper{
		installOptions: options,
		actionInstall:  installCli,
		actionClient:   client,
		actionConfig:   actionConfig,
		app:            app,
		cfg:            cfg,
	}, nil
}

func (h *Helper) GetChartRequested() (*chart.Chart, error) {
	settings, err := helmaction.NewSettings(h.installOptions.ChartRepo)
	if err != nil {
		return nil, errors.Wrap(err, "helmaction.NewSettings")
	}

	// unpack first if need
	root := settings.RepositoryCache
	if h.installOptions.ExistedFile != "" && file.IsFile(h.installOptions.ExistedFile) {
		temp, err := helmaction.ExpandFile(h.installOptions.ExistedFile, settings.RepositoryCache)
		if err != nil {
			return nil, errors.Wrap(err, "helmaction.ExpandFile")
		}
		root = temp
		defer func() {
			os.RemoveAll(temp)
		}()
	}

	var cp string
	chartDir, err := securejoin.SecureJoin(root, h.installOptions.Chart)
	if err != nil {
		return nil, errors.Wrap(err, "securejoin.SecureJoin")
	}

	cp, err = h.actionInstall.ChartPathOptions.LocateChart(chartDir, settings)
	if err != nil {
		return nil, errors.Wrap(err, "ChartPathOptions.LocateChart")
	}

	p := getter.All(settings)

	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, errors.Wrap(err, "loader.Load")
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, errors.Wrap(err, "isChartInstallable")
	}

	if chartRequested.Metadata.Deprecated {
		log.Warnf("This chart %s/%s is deprecated", h.installOptions.ChartRepo, h.installOptions.Chart)
	}
	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if h.actionInstall.DependencyUpdate {
				if err := h.actionClient.DependencyUpdate(cp, p, settings, h.installOptions.Verify, h.installOptions.Keyring); err != nil {
					return nil, errors.Wrap(err, "client.DependencyUpdate")
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "loader.Load")
				}
			} else {
				return nil, errors.Wrap(err, "action.CheckDependencies")
			}
		}
	}
	return chartRequested, nil
}

func (h *Helper) GetClusterRelease() (*release.Release, error) {
	var releaseName = h.app.Spec.Name
	latestRelease, err := h.cfg.Releases.Last(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			log.Infof("release not found, don't need to update release")
			return latestRelease, err
		}
		log.Errorf("failed to get release latest version, err:%s", err.Error())
		return nil, errors.Wrapf(err, "failed to get release '%s' latest version", releaseName)
	}

	return latestRelease, err
}

func (h *Helper) GetClusterManifest() (string, error) {
	var releaseName = h.app.Spec.Name
	latestRelease, err := h.cfg.Releases.Last(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			log.Infof("release not found, don't need to update release")
			return "", err
		}
		log.Errorf("failed to get release latest version, err:%s", err.Error())
		return "", errors.Wrapf(err, "failed to get release '%s' latest version", releaseName)
	}

	return latestRelease.Manifest, err
}

func (h *Helper) GetClusterResource() (kube.ResourceList, error) {
	var releaseName = h.app.Spec.Name
	latestRelease, err := h.cfg.Releases.Last(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			log.Infof("release not found, don't need to update release")
			return nil, err
		}
		log.Errorf("failed to get release latest version, err:%s", err.Error())
		return nil, errors.Wrapf(err, "failed to get release '%s' latest version", releaseName)
	}

	return h.cfg.KubeClient.Build(bytes.NewBufferString(latestRelease.Manifest), false)
}

func (h *Helper) GetManifest() (string, error) {
	chartRequested, err := h.GetChartRequested()
	if err != nil {
		return "", errors.Wrap(err, "unable to get chart request")
	}
	caps, err := h.cfg.GetCapabilities()
	if err != nil {
		return "", errors.Wrap(err, "unable to build kubernetes objects from release manifest")
	}

	// special case for helm template --is-upgrade
	// isUpgrade := i.IsUpgrade && i.DryRun
	releaseOptions := chartutil.ReleaseOptions{
		Name:      h.app.Spec.Name,
		Namespace: h.app.Spec.TargetNamespace,
		Revision:  1,
		IsInstall: true,
		IsUpgrade: false,
	}
	valuesToRender, err := chartutil.ToRenderValues(chartRequested, h.installOptions.Values, releaseOptions, caps)
	if err != nil {
		return "", errors.Wrap(err, "chartutil.ToRenderValues")
	}
	var manifestDoc *bytes.Buffer
	var manifest string
	// 需要根据调用集群来渲染
	// 所以dryRun必须为false
	_, manifestDoc, _, err = h.cfg.RenderResources(chartRequested, valuesToRender, h.installOptions.ReleaseName, "", false, false, false, nil, h.installOptions.DryRun)
	// Even for errors, attach this if available
	if manifestDoc != nil {
		manifest = manifestDoc.String()
	}
	return manifest, err
}

func (h *Helper) GetResourceList() (kube.ResourceList, error) {
	chartRequested, err := h.GetChartRequested()
	if err != nil {
		return nil, errors.Wrap(err, "unable to get chart request")
	}
	caps, err := h.cfg.GetCapabilities()
	if err != nil {
		return nil, errors.Wrap(err, "unable to build kubernetes objects from release manifest")
	}

	// special case for helm template --is-upgrade
	// isUpgrade := i.IsUpgrade && i.DryRun
	releaseOptions := chartutil.ReleaseOptions{
		Name:      h.app.Spec.Name,
		Namespace: h.app.Spec.TargetNamespace,
		Revision:  1,
		IsInstall: true,
		IsUpgrade: false,
	}
	valuesToRender, err := chartutil.ToRenderValues(chartRequested, h.installOptions.Values, releaseOptions, caps)
	if err != nil {
		return nil, errors.Wrap(err, "chartutil.ToRenderValues")
	}
	var manifestDoc *bytes.Buffer
	var manifest string
	// 需要根据调用集群来渲染
	// 所以dryRun必须为false
	_, manifestDoc, _, err = h.cfg.RenderResources(chartRequested, valuesToRender, h.installOptions.ReleaseName, "", false, false, false, nil, h.installOptions.DryRun)
	// Even for errors, attach this if available
	if manifestDoc != nil {
		manifest = manifestDoc.String()
	}
	resources, err := h.actionConfig.KubeClient.Build(bytes.NewBufferString(manifest), true)
	if err != nil {
		return nil, errors.Wrap(err, "unable to build kubernetes objects from release manifest")
	}
	return resources, nil
}

const defaultDirectoryPermission = 0755

// write the <data> to <output-dir>/<name>. <append> controls if the file is created or content will be appended
func writeToFile(outputDir string, name string, data string, append bool) error {
	outfileName := strings.Join([]string{outputDir, name}, string(filepath.Separator))

	err := ensureDirectoryForFile(outfileName)
	if err != nil {
		return err
	}

	f, err := createOrOpenFile(outfileName, append)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("---\n# Source: %s\n%s\n", name, data))

	if err != nil {
		return err
	}

	fmt.Printf("wrote %s\n", outfileName)
	return nil
}

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(filename)
}

// check if the directory exists to create file. creates if don't exists
func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, defaultDirectoryPermission)
}

// isChartInstallable validates if a chart can be installed
//
// Application chart type is only installable
func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

// buildChartPathBasicOptions will judge chartgroup type and return well-structured ChartPathOptions
func buildChartPathBasicOptions(repo config.RepoConfiguration, appChart applicationv1.Chart) (opt helmaction.ChartPathOptions, err error) {
	if strings.Contains(appChart.RepoURL, "tke-addon.tencentcloudcr.com") {
		password, err := registryutil.VerifyDecodedPassword(appChart.RepoPassword)
		if err != nil {
			return opt, err
		}

		opt.RepoURL = appChart.RepoURL
		opt.Username = appChart.RepoUsername
		opt.Password = password
	} else {
		opt.RepoURL = fmt.Sprintf("%s://%s", repo.Scheme, repo.DomainSuffix)
		opt.Username = repo.Admin
		opt.Password = repo.AdminPassword
	}

	opt.ChartRepo = appChart.TenantID + "/" + appChart.ChartGroupName
	opt.Chart = appChart.ChartName
	opt.Version = appChart.ChartVersion
	return opt, nil
}
