/*
 * Copyright 2019 THL A29 Limited, a Tencent company.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package res

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/thoas/go-funk"
	"tkestack.io/tke/pkg/platform/provider/baremetal/constants"
	"tkestack.io/tke/pkg/spec"
	"tkestack.io/tke/pkg/util/ssh"
)

var (
	Docker = Package{
		Name:     "docker",
		Versions: spec.DockerVersions,
	}
	CNIPlugins = Package{
		Name:     "cni-plugins",
		Versions: spec.CNIPluginsVersions,
	}
	ConntrackTools = Package{
		Name:      "conntrack-tools",
		Versions:  spec.ConntrackToolsVersions,
		TargetDir: "/",
	}

	Kubeadm = Package{
		Name:     "kubeadm",
		Versions: spec.K8sValidVersionsWithV,
	}
	KubernetesNode = Package{
		Name:     "kubernetes-node",
		Versions: spec.K8sValidVersionsWithV,
	}
	NvidiaDriver = Package{
		Name:     "NVIDIA",
		Versions: spec.NvidiaDriverVersions,
	}
	NvidiaContainerRuntime = Package{
		Name:     "nvidia-container-runtime",
		Versions: spec.NvidiaContainerRuntimeVersions,
	}
)

type Package struct {
	// Name to package must ends with .tag.gz
	Name     string
	Versions []string
	// TargetDir for untar working dir
	TargetDir string
}

func (p *Package) InstallWithDefault(s ssh.Interface) error {
	return p.Install(s, p.DefaultVersion())

}

func (p *Package) Install(s ssh.Interface, version string) error {
	dstFile, err := p.CopyToNode(s, version)
	if err != nil {
		return err
	}

	if p.TargetDir == "" {
		return errors.New("package TargetDir required")
	}
	cmd := fmt.Sprintf("tar xvaf %s -C %s ", dstFile, p.TargetDir)
	_, err = s.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("untar error: %w", err)
	}

	return nil
}

// CopyToNode copy package which use default version to node and return dst filename
func (p *Package) CopyToNodeWithDefault(s ssh.Interface) (string, error) {
	return p.CopyToNode(s, p.DefaultVersion())
}

// CopyToNode copy package which use default version to node and return dst filename
func (p *Package) CopyToNode(s ssh.Interface, version string) (string, error) {
	srcFile, err := p.ResourceForNode(s, version)
	if err != nil {
		return "", err
	}
	dstFile := path.Join(constants.DstTmpDir, filepath.Base(srcFile))
	err = s.CopyFile(srcFile, dstFile)
	if err != nil {
		return "", err
	}
	return dstFile, nil
}

func (p *Package) ResourceForNode(s ssh.Interface, version string) (string, error) {
	return p.Resource(Arch(s), version)
}

func (p *Package) Resource(arch, version string) (string, error) {
	version, err := p.NormalizeVersion(version)
	if err != nil {
		return "", err
	}
	basename := fmt.Sprintf("linux-%s/%s-linux-%s-%s.tar.gz", arch, p.Name, arch, version)
	srcFile := path.Join(constants.SrcDir, basename)
	if _, err := os.Stat(srcFile); err != nil {
		return "", err
	}
	return srcFile, nil
}

func (p *Package) DefaultVersion() string {
	return p.Versions[0]
}

func (p *Package) NormalizeVersion(version string) (string, error) {
	if p.Versions[0][0] == 'v' && version[0] != 'v' {
		version = "v" + version
	} else if p.Versions[0][0] != 'v' && version[0] == 'v' {
		version = version[1:]
	}

	if funk.ContainsString(p.Versions, version) {
		return version, nil
	}

	return "", errors.New("invalid version")
}

func Arch(s ssh.Interface) string {
	var arch string

	stdout, _, _, _ := s.Exec("arch")
	switch strings.TrimSpace(stdout) {
	case "x86_64":
		arch = "amd64"
	case "aarch64":
		arch = "arm64"
	}

	return arch
}
