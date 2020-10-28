/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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

package tencent

import (
	"k8s.io/klog"
	"os"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"k8s.io/apimachinery/pkg/util/wait"
	"tkestack.io/tke/test/util/cloudprovider"
)

func NewTencentProvider() cloudprovider.Provider {
	p := &provider{}

	credential := common.NewCredential(
		os.Getenv("SECRET_ID"),
		os.Getenv("SECRET_KEY"),
	)
	cpf := profile.NewClientProfile()
	p.cvmClient, _ = cvm.NewClient(credential, os.Getenv("REGION"), cpf)

	return p
}

type provider struct {
	cvmClient   *cvm.Client
	instanceIds []string
}

func (p *provider) CreateInstances(count int64) ([]cloudprovider.Instance, error) {
	klog.Info("Create instances. Count: ", count)
	result := make([]cloudprovider.Instance, count)

	request := cvm.NewRunInstancesRequest()
	err := request.FromJsonString(os.Getenv("CREATE_INSTANCES_PARAM"))
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	request.InstanceCount = &count
	response, err := p.cvmClient.RunInstances(request)
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	err = wait.PollImmediate(5*time.Second, 10*time.Minute, func() (bool, error) {
		describeInstancesRequest := cvm.NewDescribeInstancesRequest()
		describeInstancesRequest.InstanceIds = response.Response.InstanceIdSet
		describeInstancesResponse, err := p.cvmClient.DescribeInstances(describeInstancesRequest)
		if err != nil {
			klog.Error(err)
			return false, nil
		}
		for _, one := range describeInstancesResponse.Response.InstanceSet {
			if *one.InstanceState != "RUNNING" {
				return false, nil
			}
		}
		for i, one := range describeInstancesResponse.Response.InstanceSet {
			result[i] = cloudprovider.Instance{
				InstanceID: *one.InstanceId,
				InternalIP: *one.PrivateIpAddresses[0],
				PublicIP:   *one.PublicIpAddresses[0],
				Port:       22,
				Username:   "root",
				Password:   os.Getenv("PASSWORD"),
			}
		}
		return true, nil
	})
	if err != nil {
		klog.Error(err)
		_ = p.DeleteInstances(common.StringValues(response.Response.InstanceIdSet))
		return nil, err
	}

	for _, ins := range result {
		klog.Info("InstanceId: ", ins.InstanceID, ", PublicIP: ", ins.PublicIP, ", InternalIP: ", ins.InternalIP)
		p.instanceIds = append(p.instanceIds, ins.InstanceID)
	}

	return result, nil
}

func (p *provider) DeleteAllInstances() error {
	return p.DeleteInstances(p.instanceIds)
}

func (p *provider) DeleteInstances(instanceIDs []string) error {
	klog.Info("Delete instances: ", instanceIDs)
	if len(instanceIDs) == 0 {
		return nil
	}
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = common.StringPtrs(instanceIDs)
	_, err := p.cvmClient.TerminateInstances(request)
	if err != nil {
		return err
	}

	return nil
}
