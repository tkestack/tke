/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package configfiles

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"path/filepath"
	registryconfig "tkestack.io/tke/pkg/registry/apis/config"
	registryscheme "tkestack.io/tke/pkg/registry/apis/config/scheme"
	"tkestack.io/tke/pkg/registry/config/codec"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
)

// Loader loads configuration from a storage layer
type Loader interface {
	// Load loads and returns the RegistryConfiguration from the storage layer, or an error if a configuration could not be loaded
	Load() (*registryconfig.RegistryConfiguration, error)
}

// fsLoader loads configuration from `configDir`
type fsLoader struct {
	// fs is the filesystem where the config files exist; can be mocked for testing
	fs utilfs.Filesystem
	// registryCodecs is the scheme used to decode config files
	registryCodecs *serializer.CodecFactory
	// registryFile is an absolute path to the file containing a serialized RegistryConfiguration
	registryFile string
}

// NewFsLoader returns a Loader that loads a RegistryConfiguration from the `registryFile`
func NewFsLoader(fs utilfs.Filesystem, registryFile string) (Loader, error) {
	_, registryCodecs, err := registryscheme.NewSchemeAndCodecs()
	if err != nil {
		return nil, err
	}

	return &fsLoader{
		fs:             fs,
		registryCodecs: registryCodecs,
		registryFile:   registryFile,
	}, nil
}

func (loader *fsLoader) Load() (*registryconfig.RegistryConfiguration, error) {
	data, err := loader.fs.ReadFile(loader.registryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubelet config file %q, error: %v", loader.registryFile, err)
	}

	// no configuration is an error, some parameters are required
	if len(data) == 0 {
		return nil, fmt.Errorf("kubelet config file %q was empty", loader.registryFile)
	}

	kc, err := codec.DecodeRegistryConfiguration(loader.registryCodecs, data)
	if err != nil {
		return nil, err
	}

	// make all paths absolute
	resolveRelativePaths(registryconfig.RegistryConfigurationPathRefs(kc), filepath.Dir(loader.registryFile))
	return kc, nil
}

// resolveRelativePaths makes relative paths absolute by resolving them against `root`
func resolveRelativePaths(paths []*string, root string) {
	for _, path := range paths {
		// leave empty paths alone, "no path" is a valid input
		// do not attempt to resolve paths that are already absolute
		if len(*path) > 0 && !filepath.IsAbs(*path) {
			*path = filepath.Join(root, *path)
		}
	}
}
