package expansion

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const expansionValuesFileName = "values.yaml"

type values struct {
	Global map[string]string      `json:"global" yaml:"global"`
	Charts map[string]*chartValue `json:"charts" yaml:"charts"`
}

type chartValue struct {
	Values map[string]string `json:"values" yaml:"values"`
}

// TODO: make API to set values by calling this method
// readValues
//nolint
func (d *Driver) readValues() error {
	expansionValuesPath := d.basePath + TKEPlatformExpansionDirName + expansionValuesFileName
	_, err := os.Stat(expansionValuesPath)
	if err != nil {
		if !os.IsNotExist(err) {
			d.log.Info("stat expansionValuesPath failed %v, %v", expansionValuesPath, err)
			return err
		}
		return nil
	}
	b, err := ioutil.ReadFile(expansionValuesPath)
	if err != nil {
		d.log.Errorf("read expansionValuesPath failed %v, %v", expansionValuesPath, err)
		return err
	}

	err = yaml.Unmarshal(b, d.Values)
	if err != nil {
		d.log.Errorf("yaml unmarshal expansionValuesPath file failed %v, %v", expansionValuesPath, err)
		return err
	}
	return nil
}
