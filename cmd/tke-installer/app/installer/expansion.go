package installer

import (
	"fmt"
	"github.com/thoas/go-funk"
	"tkestack.io/tke/pkg/expansion"
	"tkestack.io/tke/pkg/mesh/util/json"
	"tkestack.io/tke/pkg/util/log"
)

func (t *TKE) newExpansionDriver() error {
	var err error
	t.expansionDriver, err = expansion.NewExpansionDriver(log.WithName("tke-expansions"))
	if err != nil {
		return err
	}
	// test tke.json from expansion
	if err := t.parseInstallerConfig(&TKE{}); err != nil {
		return err
	}
	return nil
}

// parseInstallerConfig to get tke.json from expansion, and then merge it into TKE-installer's tke.json
func (t *TKE) parseInstallerConfig(mock *TKE) error {
	raw := t.expansionDriver.InstallerConfig
	err := json.Unmarshal([]byte(raw), mock)
	if err != nil {
		return fmt.Errorf("bad expansion specified tke config %v", err)
	}
	t.log.Infof("tke.json from expansion: %+v", mock)
	return nil
}

func (t *TKE) mergeExpansionTKEConfig() error {

	// TODO: or we can use github.com/imdario/mergo to do this, which needs a fully test
	mock := &TKE{}
	if err := t.parseInstallerConfig(mock); err != nil {
		return err
	}

	// 1. append skipSteps.
	if len(mock.Para.Config.SkipSteps) > 0 {
		if t.Para.Config.SkipSteps == nil {
			t.Para.Config.SkipSteps = mock.Para.Config.SkipSteps
		} else {
			for _, step := range mock.Para.Config.SkipSteps {
				if !funk.ContainsString(t.Para.Config.SkipSteps, step) {
					t.Para.Config.SkipSteps = append(t.Para.Config.SkipSteps, step)
				}
			}
		}
	}

	// 2. append skipConditions.
	if len(mock.Cluster.Spec.Features.SkipConditions) > 0 {
		if t.Cluster.Spec.Features.SkipConditions == nil {
			t.Cluster.Spec.Features.SkipConditions = mock.Cluster.Spec.Features.SkipConditions
		} else {
			for _, cond := range mock.Cluster.Spec.Features.SkipConditions {
				if !funk.ContainsString(t.Cluster.Spec.Features.SkipConditions, cond) {
					t.Cluster.Spec.Features.SkipConditions = append(t.Cluster.Spec.Features.SkipConditions, cond)
				}
			}
		}
	}

	// 3. append copy-files, expansion/tke.json do not override installer/tke.json, just append non-exist dst-file
	if len(mock.Cluster.Spec.Features.Files) > 0 {
		if t.Cluster.Spec.Features.Files == nil {
			t.Cluster.Spec.Features.Files = mock.Cluster.Spec.Features.Files
		} else {
			for _, file := range mock.Cluster.Spec.Features.Files {
				var duplicated bool
				for _, f := range t.Cluster.Spec.Features.Files {
					if f.Dst == file.Dst {
						duplicated = true
						break
					}
				}
				if !duplicated {
					t.Cluster.Spec.Features.Files = append(t.Cluster.Spec.Features.Files, file)
				}
			}
		}
	}

	// 4. append hooks: expansion/tke.json do not override installer/tke.json, just append non-exist hooks
	if len(mock.Cluster.Spec.Features.Hooks) > 0 {
		if t.Cluster.Spec.Features.Hooks == nil {
			t.Cluster.Spec.Features.Hooks = mock.Cluster.Spec.Features.Hooks
		} else {
			for hook, hookFile := range mock.Cluster.Spec.Features.Hooks {
				if _, ok := t.Cluster.Spec.Features.Hooks[hook]; !ok {
					t.Cluster.Spec.Features.Hooks[hook] = hookFile
				}
			}
		}
	}

	// 5. merge extra args
	mergeMap(&t.Cluster.Spec.APIServerExtraArgs, &mock.Cluster.Spec.APIServerExtraArgs)
	t.log.Infof("APIServerExtraArgs %+v", t.Cluster.Spec.APIServerExtraArgs)
	mergeMap(&t.Cluster.Spec.ControllerManagerExtraArgs, &mock.Cluster.Spec.ControllerManagerExtraArgs)
	mergeMap(&t.Cluster.Spec.SchedulerExtraArgs, &mock.Cluster.Spec.SchedulerExtraArgs)
	mergeMap(&t.Cluster.Spec.KubeletExtraArgs, &mock.Cluster.Spec.KubeletExtraArgs)
	mergeMap(&t.Cluster.Spec.DockerExtraArgs, &mock.Cluster.Spec.DockerExtraArgs)
	if mock.Cluster.Spec.Etcd != nil {
		t.Cluster.Spec.Etcd = mock.Cluster.Spec.Etcd
	}

	// save to disk
	t.backup()
	return nil
}

func mergeMap(high, low *map[string]string) {
	if *low == nil {
		return
	}
	if *high == nil {
		*high = *low
		return
	}
	for k, v := range *low {
		if _, ok := (*high)[k]; !ok {
			(*high)[k] = v
		}
	}
}
