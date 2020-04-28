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

package docker

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"tkestack.io/tke/pkg/spec"

	pkgerrors "github.com/pkg/errors"
)

// Docker wraps several docker commands.
//
// A Docker instance can be reused after calling its methods.
type Docker struct {
	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Cmd Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	Stdout io.Writer
	Stderr io.Writer
}

// New returns the Docker struct for executing docker commands.
func New() *Docker {
	docker := &Docker{}
	return docker
}

// runCmd starts to execute the command specified by cmdString.
func (d *Docker) runCmd(cmdString string) error {
	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Stdout = d.Stdout
	cmd.Stderr = d.Stderr
	return cmd.Run()
}

// getCmdOutput runs the command specified by cmdString and returns its standard output.
func (d *Docker) getCmdOutput(cmdString string) ([]byte, error) {
	// print cmdString before run
	if d.Stdout != nil {
		_, _ = d.Stdout.Write([]byte(cmdString + "\n"))
	}
	cmd := exec.Command("sh", "-c", cmdString)
	return cmd.Output()
}

// GetImages returns docker images which match given image prefix.
func (d *Docker) GetImages(imagePrefix string) ([]string, error) {
	cmdString := fmt.Sprintf("docker images --format='{{.Repository}}:{{.Tag}}' --filter='reference=%s'", imagePrefix)
	out, err := d.getCmdOutput(cmdString)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "docker images error")
	}
	images := strings.Split(strings.TrimSpace(string(out)), "\n")
	return images, nil
}

// PushImageWithArch pushes an image which has a suffix about arch.
//
// It will create/amend local manifest of the image,
// and push the updated local manifest to registry if need.
// (For speed up processing, it is better to push manifests after all changes have made.)
func (d *Docker) PushImageWithArch(image string, manifestName string,
	arch string, variant string, needPushManifest bool) error {
	err := d.PushImage(image)
	if err != nil {
		return err
	}

	err = d.CreateManifest(image, manifestName)
	if err != nil {
		return err
	}

	err = d.AnnotateManifest(image, manifestName, arch, variant)
	if err != nil {
		return err
	}

	if needPushManifest {
		err = d.PushManifest(manifestName, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// PushArm64Variants accepts an arm64 image, and creates another two variants that refer to this image.
// The manifest of this arm64 image is updated accordingly.
// Current variants: unknown, v8.
func (d *Docker) PushArm64Variants(image string, name string, tag string) error {
	manifestName := fmt.Sprintf("%s:%s", name, tag)
	for _, variant := range spec.Arm64Variants {
		// variantImage: ${BIN}-arm64-${VARIANT}:${VERSION}
		variantImage := fmt.Sprintf("%s-%s-%s:%s", name, spec.Arm64, variant, tag)

		err := d.TagImage(image, variantImage)
		if err != nil {
			return err
		}

		err = d.PushImageWithArch(variantImage, manifestName, spec.Arm64, variant, false)
		if err != nil {
			return err
		}

		err = d.RemoveImage(variantImage)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetNameArchTag returns the name, arch & tag of the given image.
// If the tag is <none>, return err.
// If the image doesn't contain an arch suffix, arch = "".
func (d *Docker) GetNameArchTag(image string) (name string, arch string, tag string, err error) {
	name, tag, err = d.SplitImageNameAndTag(image)
	if err != nil { // invalid image
		return "", "", "", err
	}
	name, arch = d.SplitNameAndArch(name)
	return
}

// SplitImageNameAndTag returns the name & tag of the given image.
// If the tag is <none>, return err.
func (d *Docker) SplitImageNameAndTag(image string) (name string, tag string, err error) {
	nameAndTag := strings.Split(image, ":")
	if len(nameAndTag) != 2 {
		return "", "", fmt.Errorf("fail to get name and tag for image: %v", image)
	}
	if nameAndTag[1] == "<none>" {
		return "", "", fmt.Errorf("image %s is invalid", image)
	}
	return nameAndTag[0], nameAndTag[1], nil
}

// SplitNameAndArch returns the real name & arch of the given name.
// If the name doesn't contain an arch suffix, arch = "".
func (d *Docker) SplitNameAndArch(name string) (string, string) {
	archRegex := regexp.MustCompile(fmt.Sprintf(`(.+)-(%s)$`, strings.Join(spec.Archs, "|")))
	archMatches := archRegex.FindStringSubmatch(name)
	if archMatches == nil {
		return name, ""
	}
	return archMatches[1], archMatches[2]
}

// LoadImages loads images from a tar archive file.
func (d *Docker) LoadImages(imagesFile string) error {
	cmdString := fmt.Sprintf("docker load -i %s", imagesFile)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker load error")
	}
	return nil
}

// TagImage creates a tag destImage that refers to srcImage.
func (d *Docker) TagImage(srcImage string, destImage string) error {
	cmdString := fmt.Sprintf("docker tag %s %s", srcImage, destImage)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker tag error")
	}
	return nil
}

// PushImage pushes an image.
func (d *Docker) PushImage(image string) error {
	cmdString := fmt.Sprintf("docker push %s", image)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker push error")
	}
	return nil
}

// RemoveImage removes a local image.
func (d *Docker) RemoveImage(image string) error {
	cmdString := fmt.Sprintf("docker rmi %s ", image)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker rmi error")
	}
	return nil
}

// RemoveContainers forces to remove one or more running containers.
func (d *Docker) RemoveContainers(containers... string) error {
	// Force the removal of containers. Do not return error.
	cmdString := fmt.Sprintf("docker rm -f %s 2> /dev/null || true", strings.Join(containers, " "))
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker rm error")
	}
	return nil
}

// RunImage derives a running container from an image.
func (d *Docker) RunImage(image string, options string, runArgs string) error {
	cmdString := fmt.Sprintf("docker run %s %s %s", options, image, runArgs)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker run error")
	}
	return nil
}

// ClearLocalManifests clears all local manifest lists.
// It is better to call this method before you want to create a manifest list.
func (d *Docker) ClearLocalManifests() error {
	cmdString := "rm -rf ${HOME}/.docker/manifests/"
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "remove local manifest files error")
	}
	return nil
}

// CreateManifest creates a local manifest list. (!IMPORTANT: local,local,local!)
func (d *Docker) CreateManifest(image string, manifestName string) error {
	cmdString := fmt.Sprintf("docker manifest create --amend --insecure %s %s", manifestName, image)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker manifest create error")
	}
	return nil
}

// AnnotateManifest adds additional information to a local image manifest. (!IMPORTANT: local,local,local!)
func (d *Docker) AnnotateManifest(image string, manifestName string, arch string, variant string) error {
	if arch == "" {
		return fmt.Errorf("docker manifest annotate error: Image %s doesn't contain arch info", image)
	}

	variantArg := ""
	if variant != "" {
		variantArg = fmt.Sprintf("--variant %s", variant)
	}
	cmdString := fmt.Sprintf("docker manifest annotate %s %s --arch %s %s",
		manifestName, image, arch, variantArg)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker manifest annotate error")
	}
	return nil
}

// PushManifest pushes a manifest list.
func (d *Docker) PushManifest(manifestName string, needPurge bool) error {
	purgeArg := ""
	if needPurge {
		// Remove the local manifest list after push. !IMPORTANT: Remove local!
		purgeArg = "--purge"
	}
	cmdString := fmt.Sprintf("docker manifest push --insecure %s %s ", purgeArg, manifestName)
	err := d.runCmd(cmdString)
	if err != nil {
		return pkgerrors.Wrap(err, "docker manifest push error")
	}
	return nil
}
