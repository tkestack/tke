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

package ssh

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
)

// GetNetworkInterface return network interface name by ip
func GetNetworkInterface(s Interface, ip string) string {
	stdout, _, _, _ := s.Execf("ip a | grep '%s' |awk '{print $NF}'", ip)

	return strings.Replace(stdout, "\n", "", -1)
}

// GetDefaultRouteInterface returns default router network interface
func GetDefaultRouteInterface(s Interface) string {
	stdout, _, _, _ := s.Exec(`ip route get 1.1.1.1 | grep -oP 'dev \K\S+'`)

	return strings.Replace(stdout, "\n", "", -1)
}

// GetDefaultRouteIP returns default router network interface
func GetDefaultRouteIP(s Interface) string {
	stdout, _, _, _ := s.Exec(`ip route get 1.1.1.1 | grep -oP 'src \K\S+'`)

	return strings.Replace(stdout, "\n", "", -1)
}

// Timestamp returns target node timestamp.
func Timestamp(s Interface) (int, error) {
	stdout, err := s.CombinedOutput("date +%s")
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(stdout)))
}

func BackupFile(s Interface, file string) (string, error) {
	backup := path.Join(fmt.Sprintf("%s-%s", file, time.Now().Format("20060102150405")))
	cmd := fmt.Sprintf("mv %s %s", file, backup)
	_, err := s.CombinedOutput(cmd)
	if err != nil {
		return backup, fmt.Errorf("backup %q error: %w", file, err)
	}

	return backup, nil
}

func RestoreFile(s Interface, file string) error {
	i := strings.LastIndex(file, "-")
	if i <= 0 {
		return fmt.Errorf("invalid file name %q", file)
	}
	cmd := fmt.Sprintf("mv %s %s", file, file[0:i])
	_, err := s.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("restore %q error: %w", file, err)
	}

	return nil
}

// MemoryCapacity returns the machine's total memory from /proc/meminfo.
// Returns the total memory capacity as an uint64 (number of bytes).
func MemoryCapacity(s Interface) (uint64, error) {
	stdout, err := s.CombinedOutput(`grep 'MemTotal:' /proc/meminfo | grep -oP '\d+'`)
	if err != nil {
		return 0, err
	}

	memInKB, err := strconv.ParseUint(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, err
	}

	return memInKB * 1024, err
}

// NumCPU returns the number of logical CPUs.
func NumCPU(s Interface) (int, error) {
	stdout, err := s.CombinedOutput(`nproc --all`)
	if err != nil {
		return 0, err
	}

	cpu, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		return 0, err
	}

	return cpu, nil
}

// DiskAvail returns available disk space in GiB.
func DiskAvail(s Interface, path string) (int, error) {
	cmd := fmt.Sprintf(`df -BG %s | tail -1 | awk '{print $4}' | grep -oP '\d+'`, path)
	stdout, err := s.CombinedOutput(cmd)
	if err != nil {
		return 0, err
	}

	disk, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
	if err != nil {
		return 0, err
	}

	return disk, nil
}

// OSVersion returns os version.
func OSVersion(s Interface) (os string, err error) {
	var id, version string
	releasePath := "/etc/os-release"
	stdout, err := s.CombinedOutput("cat " + releasePath)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(stdout), "\n") {
		if strings.Contains(line, "=") {
			item := strings.Split(line, "=")
			item[0] = strings.TrimPrefix(item[0], "\"")
			item[0] = strings.TrimSuffix(item[0], "\"")
			item[1] = strings.TrimPrefix(item[1], "\"")
			item[1] = strings.TrimSuffix(item[1], "\"")
			switch item[0] {
			case "ID":
				id = item[1]
			case "VERSION_ID":
				version = item[1]
			}
		}
	}
	if len(id) == 0 {
		return "", fmt.Errorf("can not get os ID from %s", releasePath)
	}

	if len(version) == 0 {
		return "", fmt.Errorf("can not get os version ID from %s", releasePath)
	}

	return id + version, nil
}

func ReservePorts(s Interface, ip string, ports []int) (isInused bool, message string, err error) {
	var cmd string
	for _, port := range ports {
		cmd += fmt.Sprintf(`timeout 3 bash -c "</dev/tcp/%s/%d" &>/dev/null; echo $?; `, ip, port)
	}
	out, _, _, _ := s.Exec(cmd)
	out = strings.TrimSuffix(out, "\n")
	results := strings.Split(out, "\n")
	if len(results) != len(ports) {
		return false, "", fmt.Errorf("check results length does not match need check ports length, get results output is: %s", out)
	}
	for i, result := range results {
		// if return code is 124, it means that the connection is timeout
		if result == "124" {
			return false, "", fmt.Errorf("connect %s:%d timeout", ip, ports[i])
		}
		if result != "1" {
			message += fmt.Sprintf("%d ", ports[i])
		}
	}
	if len(message) != 0 {
		return true, fmt.Sprintf("ports %sis in used", message), nil
	}
	return false, "", nil
}

func FirewallEnabled(s Interface) (enabled bool, err error) {
	ostype, err := OSVersion(s)
	if err != nil {
		return false, err
	}
	switch {
	case strings.HasPrefix(ostype, "tencentos"):
		stdout, err := s.CombinedOutput("ps -ef | grep firewalld | grep -v grep | wc -l")
		if err != nil {
			return false, err
		}
		res := strings.TrimSpace(string(stdout))
		return res == "1", nil
	case strings.Contains(ostype, "ubuntu"):
		stdout, _, exit, err := s.Exec("ufw status | awk '{print $2}'")
		if err != nil || exit != 0 {
			return false, err
		}
		res := strings.TrimSpace(stdout)
		return res == "active", nil
	default:
		stdout, err := s.CombinedOutput("ps -ef | grep firewalld | grep -v grep | wc -l")
		if err != nil {
			return false, err
		}
		res := strings.TrimSpace(string(stdout))
		return res == "1", nil
	}

}

func SelinuxEnabled(s Interface) (enabled bool, err error) {
	// https://www.thegeekdiary.com/how-to-check-whether-selinux-is-enabled-or-disabled/
	_, _, exit, err := s.Exec("selinuxenabled")
	if err != nil {
		return false, err
	}
	return exit == 0, nil
}

func CheckNFS(s Interface, server string, path string) (err error) {
	_, stderr, exit, err := s.Execf("mkdir -p /tmp/nfs/ && mount -t nfs -o soft,timeo=15,retry=0 %s:%s /tmp/nfs/ && umount /tmp/nfs/ && rm -rf /tmp/nfs/", server, path)
	if exit != 0 || err != nil {
		return fmt.Errorf("check nfs failed:exit %d:stderr %s:error %s", exit, stderr, err)
	}
	return nil
}
