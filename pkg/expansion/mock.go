package expansion

import (
	"github.com/thoas/go-funk"
)

// TKE we cannot import it cause of cycle
type TKE struct {
	Para    *CreateClusterPara `json:"para"`
	Cluster *Cluster           `json:"cluster"`
}
type CreateClusterPara struct {
	Config Config `json:"Config"`
}
type Config struct {
	SkipSteps []string `json:"skipSteps,omitempty"`
}
type Cluster struct {
	*V1Cluster
}
type V1Cluster struct {
	// +optional
	Spec ClusterSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}
type ClusterSpec struct {
	// +optional
	Features ClusterFeature `json:"features,omitempty" protobuf:"bytes,11,opt,name=features,casttype=ClusterFeature"`
}
type ClusterFeature struct {
	// +optional
	SkipConditions []string `json:"skipConditions,omitempty" protobuf:"bytes,7,opt,name=skipConditions"`
	// +optional
	Files []File `json:"files,omitempty" protobuf:"bytes,8,opt,name=files"`
	// +optional
	Hooks map[string]string `json:"hooks,omitempty" protobuf:"bytes,9,opt,name=hooks"`
}
type File struct {
	Src string `json:"src" protobuf:"bytes,1,name=src"` // Only support regular file
	Dst string `json:"dst" protobuf:"bytes,2,name=dst"`
}

// genTKEConfig is a helper when building expansion package. It generates tke.json from expansion layout
func (l *layout) genTKEConfig() *TKE {
	var tke = &TKE{
		Para: &CreateClusterPara{
			Config: Config{
				SkipSteps: nil,
			},
		},
		Cluster: &Cluster{
			V1Cluster: &V1Cluster{
				Spec: ClusterSpec{
					Features: ClusterFeature{
						SkipConditions: nil,
						Files:          nil,
						Hooks:          nil,
					},
				},
			},
		},
	}
	tke.Para.Config.SkipSteps = l.mergeInstallerSkipSteps(tke.Para.Config.SkipSteps)
	l.mergeCluster(tke.Cluster)
	return tke
}

// mergeInstallerSkipSteps sets up install skip steps of TKEStack-installer by expansion specifying
func (l *layout) mergeInstallerSkipSteps(skipSteps []string) []string {
	// merge skip steps
	if !l.containsSkipSteps() {
		return skipSteps
	}
	if skipSteps == nil {
		skipSteps = make([]string, 0)
	}
	for _, step := range l.InstallerSkipSteps {
		if !funk.ContainsString(skipSteps, step) {
			skipSteps = append(skipSteps, step)
		}
	}
	return skipSteps
}

// mergeCluster
// 1. register copy-files with TKEStack hook config.
// 2. sets up create-cluster skip steps by expansion specifying
func (l *layout) mergeCluster(cluster *Cluster) {
	// merge file hook config
	if l.containsFiles() {
		if len(cluster.Spec.Features.Hooks) == 0 {
			cluster.Spec.Features.Hooks = make(map[string]string)
		}
		for _, f := range l.Files {
			ff := l.toFlatPath(f)
			src, dst := expansionFilesGeneratedPath+ff, absolutePath+f
			cluster.Spec.Features.Files = append(cluster.Spec.Features.Files, File{
				Src: src,
				Dst: dst,
			})
			hookType, ok := l.isHookScript(f)
			if ok {
				cluster.Spec.Features.Hooks[string(hookType)] = dst
			}
		}
	}
	// merge skip conditions
	if l.containsSkipConditions() {
		if len(cluster.Spec.Features.SkipConditions) == 0 {
			cluster.Spec.Features.SkipConditions = make([]string, 0)
		}
		for _, skipC := range l.CreateClusterSkipConditions {
			if !funk.ContainsString(cluster.Spec.Features.SkipConditions, skipC) {
				cluster.Spec.Features.SkipConditions = append(cluster.Spec.Features.SkipConditions, skipC)
			}
		}
	}

}
