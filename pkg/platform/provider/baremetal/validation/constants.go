/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2022 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package validation

const (
	AnywhereValidateItemTunnelConnectivity     = "TunnelConnectivity"
	AnywhereValidateItemSSH                    = "SSH"
	AnywhereValidateItemTimeDiff               = "TimeDiff"
	AnywhereValidateItemOSVersion              = "OS"
	AnywhereValidateItemMachineResourceDiskLib = "MachineResourceDiskLib"
	AnywhereValidateItemMachineResourceDiskLog = "MachineResourceDiskLog"
	AnywhereValidateItemMachineResourceCPU     = "MachineResourceCPU"
	AnywhereValidateItemMachineResourceMemory  = "MachineResourceMemory"
	AnywhereValidateItemDefaultRoute           = "DefaultRoute"
	AnywhereValidateItemReservePorts           = "ReservePorts"
	AnywhereValidateItemHostNetOverlapping     = "HostNetOverlapping"
	AnywhereValidateItemFirewall               = "Firewall"
	AnywhereValidateItemSelinux                = "Selinux"
	AnywhereValidateItemStorage                = "Storage"
	// validate all items
	AnywhereValidateItemAll = "All"
)

const (
	MachineResourceRequstDiskPath     = "/var/lib"
	MachineResourceRequstLogDiskPath  = "/var/log"
	MachineResourceRequstDiskSpace    = 100 // GiB
	MachineResourceRequstLogDiskSpace = 10  // GiB
	MachineResourceRequstCPU          = 4
	MachineResourceRequstMemory       = 8 // GiB
)
