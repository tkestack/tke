package expansion

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
	platformv1 "tkestack.io/tke/api/platform/v1"
)

const expansionLayoutSpecName = "layout.yaml"
const expansionChartsDir = "charts/"
const expansionFilesDir = "files/"
const expansionFilesGeneratedDir = "files_generated/"
const expansionHooksDir = "hooks/"
const expansionProviderDir = "provider/"
const expansionImagesName = "images.tar.gz"
const expansionConfDir = "conf/"
const expansionApplicationDir = "applications/"

var (
	expansionLayoutSpecPath     string
	expansionChartsPath         string
	expansionFilesPath          string
	expansionFilesGeneratedPath string
	expansionHooksPath          string
	expansionProviderPath       string
	expansionImagesPath         string
	expansionConfPath           string
	expansionApplicationPath    string
)

// layout is the description of expansion package
// 1. it is formed with Charts/Files/Hooks/Provider/Applications/Images path
// 2. it contains hooks of steps such as installerSkipSteps/CreateClusterSkipConditions
type layout struct {
	expansionBasePath               string
	paths                           map[string]string
	Charts                          []string `json:"charts" yaml:"charts"`
	Files                           []string `json:"files" yaml:"files"`
	Hooks                           []string `json:"hooks" yaml:"hooks"`
	Provider                        []string `json:"provider" yaml:"provider"`
	Applications                    []string `json:"applications" yaml:"applications"`
	Images                          []string `json:"images" yaml:"images"`
	InstallerSkipSteps              []string `json:"installer_skip_steps" yaml:"installer_skip_steps"`
	CreateClusterSkipConditions     []string `json:"create_cluster_skip_conditions" yaml:"create_cluster_skip_conditions"`
	CreateClusterDelegateConditions []string `json:"create_cluster_delegate_conditions" yaml:"create_cluster_delegate_conditions"`
}

func (d *Driver) newLayout() {
	d.layout = &layout{
		expansionBasePath: d.expansionBasePath,
	}
	expansionLayoutSpecPath = d.layout.expansionBasePath + expansionLayoutSpecName
	expansionChartsPath = d.layout.expansionBasePath + expansionChartsDir
	expansionFilesPath = d.layout.expansionBasePath + expansionFilesDir
	expansionFilesGeneratedPath = d.layout.expansionBasePath + expansionFilesGeneratedDir
	expansionHooksPath = d.layout.expansionBasePath + expansionHooksDir
	expansionProviderPath = d.layout.expansionBasePath + expansionProviderDir
	expansionImagesPath = d.layout.expansionBasePath + expansionImagesName
	expansionConfPath = d.layout.expansionBasePath + expansionConfDir
	expansionApplicationPath = d.layout.expansionBasePath + expansionApplicationDir
	d.layout.paths = map[string]string{
		"expansionLayoutSpecPath":     expansionLayoutSpecPath,
		"expansionChartsPath":         expansionChartsPath,
		"expansionFilesPath":          expansionFilesPath,
		"expansionFilesGeneratedPath": expansionFilesGeneratedPath,
		"expansionHooksPath":          expansionHooksPath,
		"expansionProviderPath":       expansionProviderPath,
		"expansionImagesPath":         expansionImagesPath,
		"expansionConfPath":           expansionConfPath,
		"expansionApplicationPath":    expansionApplicationPath,
	}
}

// readLayout
func (d *Driver) readLayout() error {
	_, err := os.Stat(expansionLayoutSpecPath)
	if err != nil {
		if !os.IsNotExist(err) {
			d.log.Info("stat expansionLayoutSpecPath failed %v, %v", expansionLayoutSpecPath, err)
			return err
		}
		return nil
	}
	b, err := ioutil.ReadFile(expansionLayoutSpecPath)
	if err != nil {
		d.log.Errorf("read expansionLayoutSpecPath failed %v, %v", expansionLayoutSpecPath, err)
		return err
	}

	err = yaml.Unmarshal(b, d.layout)
	if err != nil {
		d.log.Errorf("yaml unmarshal expansionLayoutSpecPath file failed %v, %v", expansionLayoutSpecPath, err)
		return err
	}
	return nil
}

func (d *Driver) printLayout() {
	b, err := yaml.Marshal(d.layout)
	if err != nil {
		d.log.Errorf("print expansion layout error. %v", err)
		return
	}
	d.log.Infof("expansion layout: %v", string(b))
}

func (l *layout) containsApplications() bool {
	return len(l.Applications) > 0
}
func (l *layout) containsCharts() bool {
	return len(l.Charts) > 0
}
func (l *layout) containsFiles() bool {
	return len(l.Files) > 0
}
func (l *layout) containsHooks() bool {
	return len(l.Hooks) > 0
}
func (l *layout) containsProvider() bool {
	return len(l.Provider) > 0
}
func (l *layout) containsImages() bool {
	return len(l.Images) > 0
}
func (l *layout) containsSkipSteps() bool {
	return len(l.InstallerSkipSteps) > 0
}
func (l *layout) containsSkipConditions() bool {
	return len(l.CreateClusterSkipConditions) > 0
}
func (l *layout) containsDelegateConditions() bool {
	return len(l.CreateClusterDelegateConditions) > 0
}

// makeFlatFiles for generating runtime "copy-files"
// 1. framework will put them into installer copy-file paths
// 2. framework will put them into platform configmap
//nolint
func (l layout) makeFlatFiles() error {
	if !l.containsFiles() {
		return nil
	}
	for _, f := range l.Files {
		src := expansionFilesPath + f
		ff := l.toFlatPath(f)
		flatFile := expansionFilesGeneratedPath + ff
		_, err := os.Stat(flatFile)
		if err == nil {
			continue
		}
		if os.IsNotExist(err) {
			err = copyFile(src, flatFile, 0)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("stat flatFile system error %v %v", flatFile, err)
		}
	}
	return nil
}

func (l *layout) toFlatPath(f string) string {
	return strings.Replace(f, string(os.PathSeparator), expansionFilePathSeparator, -1)
}

func (l *layout) isHookScript(fp string) (platformv1.HookType, bool) {
	var createClusterHookFileTypes = map[platformv1.HookType]bool{
		platformv1.HookPreInstall:         true,
		platformv1.HookPostInstall:        true,
		platformv1.HookPreClusterInstall:  true,
		platformv1.HookPostClusterInstall: true,
	}
	hookType := platformv1.HookType(path.Base(fp))
	if _, ok := createClusterHookFileTypes[hookType]; ok {
		return hookType, true
	}
	return "", false
}
