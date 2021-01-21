/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
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
 *
 */

package configfiles

import (
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	meshconfig "tkestack.io/tke/pkg/mesh/apis/config"
	meshscheme "tkestack.io/tke/pkg/mesh/apis/config/scheme"
	"tkestack.io/tke/pkg/mesh/config/codec"
	utilfs "tkestack.io/tke/pkg/util/filesystem"
)

// Loader loads configuration from a storage layer
type Loader interface {
	// Load loads and returns the MeshConfiguration from the storage layer,
	// or an error if a configuration could not be loaded
	Load() (*meshconfig.MeshConfiguration, error)
}

// fsLoader loads configuration from `configDir`
type fsLoader struct {
	// fs is the filesystem where the config files exist; can be mocked for testing
	fs utilfs.Filesystem
	// meshCodecs is the scheme used to decode config files
	meshCodecs *serializer.CodecFactory
	// meshFile is an absolute path to the file containing a serialized MeshConfiguration
	meshFile string
}

// NewFsLoader returns a Loader that loads a MeshConfiguration from the `meshFile`
func NewFsLoader(fs utilfs.Filesystem, meshFile string) (Loader, error) {
	_, meshCodecs, err := meshscheme.NewSchemeAndCodecs()
	if err != nil {
		return nil, err
	}

	return &fsLoader{
		fs:         fs,
		meshCodecs: meshCodecs,
		meshFile:   meshFile,
	}, nil
}

func (loader *fsLoader) Load() (*meshconfig.MeshConfiguration, error) {
	data, err := loader.fs.ReadFile(loader.meshFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubelet config file %q, error: %v", loader.meshFile, err)
	}

	// no configuration is an error, some parameters are required
	if len(data) == 0 {
		return nil, fmt.Errorf("kubelet config file %q was empty", loader.meshFile)
	}

	cfg, err := codec.DecodeConfiguration(loader.meshCodecs, data)
	if err != nil {
		return nil, err
	}

	// make all paths absolute
	resolveRelativePaths(nil, filepath.Dir(loader.meshFile))
	return cfg, nil
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
