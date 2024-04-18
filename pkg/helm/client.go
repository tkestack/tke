package helmchart

import (
	"fmt"
	"os"
	"path/filepath"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/registry"
	helmaction "tkestack.io/tke/pkg/application/helm/action"
	"tkestack.io/tke/pkg/util/log"
)

func Download(chartPathOptions helmaction.ChartPathOptions) (string, error) {
	actionConfig := new(action.Configuration)
	var err error
	actionConfig.RegistryClient, err = registry.NewClient()
	if err != nil {
		return "", err
	}
	client := action.NewPullWithOpts(action.WithConfig(actionConfig))

	settings, err := helmaction.NewSettings(chartPathOptions.ChartRepo)
	if err != nil {
		log.Errorf("NewSettings failed,err:%s", err.Error())
		return "", err
	}
	client.Settings = settings
	client.DestDir = settings.RepositoryCache
	client.Untar = false

	chartPathOptions.ApplyTo(&client.ChartPathOptions)

	if err := os.MkdirAll(settings.RepositoryCache, 0755); err != nil {
		log.Errorf("MkdirAll failed,err:%s", err.Error())
		return "", err
	}
	_, err = client.Run(chartPathOptions.Chart)
	if err != nil {
		log.Errorf("client run failed,err:%s", err.Error())
		return "", err
	}

	destTmpfile := filepath.Join(client.DestDir, fmt.Sprintf("%s-%s.tgz", chartPathOptions.Chart, chartPathOptions.Version))
	temp, err := helmaction.ExpandFile(destTmpfile, settings.RepositoryCache)
	if err != nil {
		log.Errorf("client expandFile failed,err:%s", err.Error())
		return "", err
	}
	destfile := filepath.Join(temp, chartPathOptions.Chart)
	return destfile, nil
}
