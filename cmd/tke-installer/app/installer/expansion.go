package installer

import (
	"context"
	"tkestack.io/tke/pkg/util/docker"
)

func (t *TKE) mergeExpansionsProviderConf() {
	err := t.expansionDriver.MergeProvider()
	// TODO: what to do with error
	if err != nil {
		t.log.Errorf("MergeProvider failed %v", err)
	}
}

func (t *TKE) mergeExpansionsConfig() {
	t.Para.Config.SkipSteps = t.expansionDriver.MergeInstallerSkipSteps(t.Para.Config.SkipSteps)
	// TODO: this is not working yet
	//t.expansionDriver.MergeCustomizedImages(&t.Config.PrepareCustomK8sImages)
}
func (t *TKE) mergeExpansionsCluster() {
	//t.expansionDriver.MergeInstallerSkipSteps(&t.Para.Config.SkipSteps)
	t.expansionDriver.MergeCluster(t.Cluster)
	//t.expansionDriver.MergeCustomizedImages(&t.Config.PrepareCustomK8sImages)
	t.backup()
}

func (t *TKE) loadOperatorImage(ctx context.Context) error {
	return t.expansionDriver.LoadOperatorImage(ctx)
}

func (t *TKE) patchHookFiles(ctx context.Context) error {
	return t.expansionDriver.PatchHookFiles(ctx)
}

func (t *TKE) startOperator(ctx context.Context) error {
	return t.expansionDriver.StartOperator(ctx)
}

func (t *TKE) loadExpansionImages(ctx context.Context, docker *docker.Docker) error {
	return t.expansionDriver.LoadExpansionImages(ctx, docker)
}

func (t *TKE) expansionWriteKubeconfigFile(ctx context.Context, data []byte) error {
	return t.expansionDriver.WriteKubeconfigFile(ctx, data)
}

func (t *TKE) installExpansionApplications(ctx context.Context) error {
	// TODO: transfer tke installer variables to expansion
	tkeValues := t.values()
	return t.expansionDriver.InstallApplications(ctx, t.applicationClient, tkeValues)
}

func (t *TKE) values() map[string]string {
	ret := make(map[string]string)
	// TODO: set useful things
	return ret
}
