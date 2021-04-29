package addon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"tkestack.io/tke/pkg/util/containerregistry"
)

const versionMapPathPrefix = "/app/addon/images/"

func GetVersionMap(addon string) (map[string]containerregistry.Image, error) {
	if len(addon) == 0 {
		return nil, fmt.Errorf("the name of an addon must be specified")
	}

	path := versionMapPathPrefix + addon
	// check whether the file exists
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open addon's version map from %s: %v", path, err)
	}
	defer file.Close()

	versionMap := make(map[string]containerregistry.Image)

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read addon's version map from %s: %v", path, err)
	}

	err = json.Unmarshal(b, &versionMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal version map %s: %v", addon, err)
	}

	return versionMap, nil
}

// GetLatestVersion returns latest version
func GetLatestVersion(addon string) (string, error) {
	versionMap, err := GetVersionMap(addon)
	if err != nil {
		return "", fmt.Errorf("get version map error: %v", err)
	}
	cv, ok := versionMap["latest"]
	if !ok {
		return "", fmt.Errorf("the component version definition corresponding to version %s could not be found，return default instead", cv.Tag)
	}

	return cv.Tag, nil
}
