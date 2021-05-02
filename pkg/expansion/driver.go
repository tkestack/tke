package expansion

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"tkestack.io/tke/pkg/util/log"
)

// something cannot import from tke lib
const TKENamespace = "tke"
const TKEPlatformExpansionConfigmapName = "platform-expansion"
const TKEPlatformExpansionFilesConfigmapName = "platform-expansion-files"

const TKEPlatformRootDir = "/app/"
const TKEPlatformDataDir = "data/"
const TKEPlatformBase = TKEPlatformRootDir + TKEPlatformDataDir
const TKEPlatformExpansionDirName = "expansions/"
const InstallerProviderBasePath = "provider/"

//const InstallerProviderBareMetalPath = InstallerProviderBasePath + "baremetal/"
//const InstallerHooksDir = "hooks/"
//
const defaultExpansionConfigName = "expansion.yaml"
const TKEConfigName = "tke.json"

//const defaultExpansionConfigPath = defaultExpansionBase + defaultExpansionConfigName

// expansion path and separators
const absolutePath = "/"
const expansionFilePathSeparator = "__"

//
const envExpansionBase = "EXPANSION_BASE"
const fileSuffixYaml = ".yaml" //nolint

// Driver is the model to handle expansion layout.
type Driver struct {
	log               log.Logger
	basePath          string
	expansionBasePath string
	enabled           bool
	ExpansionName     string `json:"expansion_name" yaml:"expansion_name"`
	ExpansionVersion  string `json:"expansion_version" yaml:"expansion_version"`
	RegistryNamespace string `json:"registry_namespace,omitempty" yaml:"registry_namespace" comment:"default is expansionName"`
	InstallerConfig   string `json:"installer_config" yaml:"installer_config"`
	layout            *layout
	Values            *values

	// TODO: not designed yet
	K8sVersion string `json:"k8s_version" yaml:"k8s_version"`
	// TODO: save image lists in case of installer restart to avoid load images again
	ImagesLoaded bool `json:"images_loaded" yaml:"images_loaded"`
	// TODO: if true, will prevent to pass the same expansion to cluster C
	DisablePassThroughToPlatform bool `json:"disable_pass_through_to_platform" yaml:"disable_pass_through_to_platform"`
}

func (d *Driver) enableApplications() bool {
	return d.isEnabled() && d.layout.containsApplications()
}
func (d *Driver) enableCharts() bool {
	return d.isEnabled() && d.layout.containsCharts()
}
func (d *Driver) enableFiles() bool {
	return d.isEnabled() && d.layout.containsFiles()
}
func (d *Driver) enableHooks() bool {
	return d.isEnabled() && d.layout.containsHooks()
}
func (d *Driver) enableProvider() bool {
	return d.isEnabled() && d.layout.containsProvider()
}
func (d *Driver) enableImages() bool {
	return d.isEnabled() && d.layout.containsImages()
}
func (d *Driver) enableSkipSteps() bool {
	return d.isEnabled() && d.layout.containsSkipSteps()
}
func (d *Driver) enableSkipConditions() bool {
	return d.isEnabled() && d.layout.containsSkipConditions()
}
func (d *Driver) enableDelegateConditions() bool {
	return d.isEnabled() && d.layout.containsDelegateConditions()
}
func (d *Driver) enable() {
	d.enabled = true
}
func (d *Driver) isEnabled() bool {
	return d.enabled
}

// NewExpansionDriver returns an expansionDriver instance which has all expansion layout items loaded.
func NewExpansionDriver(logger log.Logger) (*Driver, error) {
	driver := &Driver{
		log:             logger,
		Values:          &values{},
		InstallerConfig: "{}",
	}
	driver.basePath = os.Getenv(envExpansionBase)
	if driver.basePath == "" {
		driver.basePath = TKEPlatformBase
	}
	driver.expansionBasePath = driver.basePath + TKEPlatformExpansionDirName
	driver.log.Infof("tke-expansion base path %v", driver.basePath)

	driver.newLayout()

	err := driver.readConfig()
	if err != nil {
		return driver, err
	}
	driver.log.Infof("new expansion driver %+v", driver)
	err = driver.scanLayout()
	if err != nil {
		return driver, err
	}
	err = driver.readTKEConfig()
	if err != nil {
		return driver, err
	}

	// TODO: drop this when all finish
	driver.nolint()
	return driver, nil
}

// readConfig reads expansion config from config file
func (d *Driver) readConfig() error {
	expansionConfigPath := fmt.Sprintf("%s%s", d.expansionBasePath, defaultExpansionConfigName)
	_, err := os.Stat(expansionConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			d.log.Errorf("stat expansionConfigPath failed %v, %v", expansionConfigPath, err)
			return err
		}
		d.log.Infof("do not find expansion config file %v, expansion disabled", expansionConfigPath)
		return nil
	}
	b, err := ioutil.ReadFile(expansionConfigPath)
	if err != nil {
		d.log.Errorf("read expansionConfigPath failed %v, %v", expansionConfigPath, err)
		return err
	}

	err = yaml.Unmarshal(b, d)
	if err != nil {
		d.log.Errorf("yaml unmarshal expansionConfigPath file failed %v, %v", expansionConfigPath, err)
		return err
	}
	d.enable()
	d.log.Infof("expansion driver inited %+v", d)
	return nil
}

// scanLayout looks up for all files/charts/images in layout.yaml, verifies them.
func (d *Driver) scanLayout() error {
	if !d.isEnabled() {
		return nil
	}
	err := d.readLayout()
	if err != nil {
		d.log.Errorf("read expansion config failed %v", err)
		return err
	}

	// prepare generated dir
	for _, dir := range []string{
		expansionConfPath,
		expansionFilesGeneratedPath,
	} {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			d.log.Errorf("mkdir expansion generate path failed %v,%v", dir, err)
			return err
		}
	}

	// verify
	err = d.verify()
	if err != nil {
		d.log.Errorf("expansion verify failed %v", err)
		return err
	}

	return nil
}

func (d *Driver) verify() error {
	// TODO: verify files,charts,hooks,provider,images
	return nil
}

//nolint
func (d *Driver) backup() error {
	// TODO: backup expansion config file
	return nil
}

// generate
// 1. rend and copy "copyFiles" from layout to files_generate directory
// 2. rend and override "provider files"
// TODO: make API to call this method when values ready
//nolint
func (d *Driver) generate() error {
	// prepare "copy files"
	err := d.layout.makeFlatFiles()
	if err != nil {
		d.log.Errorf("expansion makeFlatFiles failed %v", err)
		return err
	}
	// prepare "provider files"
	err = d.makeProviderFiles()
	if err != nil {
		d.log.Errorf("expansion rendProviderFiles failed %v", err)
		return err
	}
	return nil
}

// makeProviderFiles overrides provider files in TKEStack installer by expansion specifying
//nolint
func (d *Driver) makeProviderFiles() error {
	if !d.enableProvider() {
		return nil
	}
	for _, f := range d.layout.Provider {
		src := expansionProviderPath + f
		dst := d.basePath + InstallerProviderBasePath + f
		err := copyFile(src, dst, 0)
		if err != nil {
			return fmt.Errorf("merge provider failed, copy %v to %v, %v", src, dst, err)
		}
	}
	return nil
}

func (d *Driver) readTKEConfig() error {
	if !d.isEnabled() {
		return nil
	}
	tkeConfigPath := fmt.Sprintf("%s%s", d.expansionBasePath, TKEConfigName)
	_, err := os.Stat(tkeConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			d.log.Errorf("stat tkeConfigPath failed %v, %v", tkeConfigPath, err)
			return err
		}
		d.log.Infof("do not find installer config file %v", tkeConfigPath)
		return nil
	}
	b, err := ioutil.ReadFile(tkeConfigPath)
	if err != nil {
		d.log.Errorf("read tkeConfigPath failed %v, %v", tkeConfigPath, err)
		return err
	}
	// TODO: check if it can be unMarshaled
	d.InstallerConfig = string(b)
	return nil
}

func (d *Driver) nolint() {
	_ = d.enableImages()
	_ = d.enableApplications()
	_ = d.enableProvider()
	_ = d.enableFiles()
	_ = d.enableCharts()
	_ = d.enableHooks()
	_ = d.enableSkipConditions()
	_ = d.enableSkipSteps()
	_ = d.enableDelegateConditions()
}
