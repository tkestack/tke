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
	"fmt"
	"k8s.io/klog"
	"time"
	"tkestack.io/tke/test/util/env"

	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"k8s.io/apimachinery/pkg/util/wait"
	"tkestack.io/tke/test/util/cloudprovider"
)

func NewTencentProvider() cloudprovider.Provider {
	p := &provider{}

	credential := common.NewCredential(
		env.SecretID(),
		env.SecretKey(),
	)
	cpf := profile.NewClientProfile()
	p.cvmClient, _ = cvm.NewClient(credential, env.Region(), cpf)
	p.clbClient, _ = clb.NewClient(credential, env.Region(), cpf)
	return p
}

type provider struct {
	cvmClient   *cvm.Client
	clbClient   *clb.Client
	instanceIds []string
	lbIds       []string
}

func (p *provider) CreateInstances(count int64) ([]cloudprovider.Instance, error) {
	klog.Info("Create instances. Count: ", count)
	result := make([]cloudprovider.Instance, count)

	request := cvm.NewRunInstancesRequest()
	err := request.FromJsonString(env.CreateInstancesParam())
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
				Password:   env.Password(),
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
		klog.Info("InstanceId: ", ins.InstanceID, ", InternalIP: ", ins.InternalIP)
		p.instanceIds = append(p.instanceIds, ins.InstanceID)
	}
	time.Sleep(10 * time.Second)

	return result, nil
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

func (p *provider) CreateCLB(eniIps []*string) (*string, error) {
	vips, err := p.CreateCLBs(eniIps, 1)
	if err != nil {
		return nil, err
	}
	return vips[0], err
}

func (p *provider) CreateCLBs(eniIps []*string, count uint64) ([]*string, error) {
	klog.Info("Create LB. EniIPs: ", common.StringValues(eniIps))
	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerName = common.StringPtr("tkestack")
	request.Number = common.Uint64Ptr(count)
	request.AddressIPVersion = common.StringPtr("IPV4")
	request.Forward = common.Int64Ptr(1)
	request.LoadBalancerType = common.StringPtr("OPEN")
	request.ProjectId = common.Int64Ptr(0)
	// Get vpcId from env
	req := cvm.NewRunInstancesRequest()
	err := req.FromJsonString(env.CreateInstancesParam())
	if err != nil {
		return nil, err
	}
	request.VpcId = req.VirtualPrivateCloud.VpcId
	rsp, err := p.clbClient.CreateLoadBalancer(request)
	if err != nil {
		return nil, fmt.Errorf("CreateLoadBalancer failed. %v", err)
	}
	p.lbIds = append(p.lbIds, common.StringValues(rsp.Response.LoadBalancerIds)...)
	klog.Info("LB created. ", common.StringValues(rsp.Response.LoadBalancerIds))

	var vips []*string
	err = p.WaitLBReady(rsp.Response.LoadBalancerIds)
	if err != nil {
		return vips, err
	}
	lbs, err := p.DescribeLBs(rsp.Response.LoadBalancerIds)
	for _, lb := range lbs {
		vips = append(vips, lb.LoadBalancerVips...)
		listenerID, err := p.CreateListener(lb.LoadBalancerId, 6443)
		if err != nil {
			return vips, err
		}
		err = p.RegisterTargets(lb.LoadBalancerId, listenerID, eniIps)
		if err != nil {
			return vips, err
		}
	}
	return vips, err
}

func (p *provider) WaitLBReady(lbIds []*string) error {
	klog.Info("Wait LB ready. lbIds: ", common.StringValues(lbIds))
	err := wait.Poll(time.Second, time.Minute, func() (done bool, err error) {
		lbs, err := p.DescribeLBs(lbIds)
		if err != nil {
			return false, nil
		}
		for _, lb := range lbs {
			if *lb.Status != 1 {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("wait LB ready failed. %v", err)
	}
	klog.Info("LB is ready")
	return err
}

func (p *provider) DeleteLBs(lbIds []*string) error {
	klog.Info("Delete LBs: ", common.StringValues(lbIds))
	if len(lbIds) == 0 {
		return nil
	}
	req := clb.NewDeleteLoadBalancerRequest()
	req.LoadBalancerIds = lbIds
	rsp, err := p.clbClient.DeleteLoadBalancer(req)
	if err != nil {
		return fmt.Errorf("DeleteLoadBalancer failed. %v", err)
	}
	err = wait.Poll(time.Second, time.Minute, func() (bool, error) {
		req1 := clb.NewDescribeTaskStatusRequest()
		req1.TaskId = rsp.Response.RequestId
		rsp1, err := p.clbClient.DescribeTaskStatus(req1)
		return *rsp1.Response.Status == 0, err
	})
	if err != nil {
		return fmt.Errorf("wait LB to be deleted failed. %v", err)
	}
	return err
}

func (p *provider) DescribeLBs(lbIds []*string) ([]*clb.LoadBalancer, error) {
	req := clb.NewDescribeLoadBalancersRequest()
	req.LoadBalancerIds = lbIds
	req.Offset = common.Int64Ptr(0)
	req.Limit = common.Int64Ptr(20)
	rsp, err := p.clbClient.DescribeLoadBalancers(req)
	if err != nil {
		return nil, fmt.Errorf("DescribeLoadBalancers failed. %v", err)
	}
	return rsp.Response.LoadBalancerSet, err
}

func (p *provider) CreateListener(lbID *string, port int64) (*string, error) {
	klog.Infof("Create listener. lbID: %v; port: %v", *lbID, port)
	req := clb.NewCreateListenerRequest()
	req.LoadBalancerId = lbID
	req.Ports = []*int64{common.Int64Ptr(port)}
	req.Protocol = common.StringPtr("TCP")
	req.Scheduler = common.StringPtr("WRR")
	rsp, err := p.clbClient.CreateListener(req)
	if err != nil {
		return nil, fmt.Errorf("create listener failed. %v", err)
	}
	err = wait.Poll(time.Second, time.Minute, func() (bool, error) {
		req1 := clb.NewDescribeTaskStatusRequest()
		req1.TaskId = rsp.Response.RequestId
		rsp1, err := p.clbClient.DescribeTaskStatus(req1)
		return *rsp1.Response.Status == 0, err
	})
	if err != nil {
		return nil, fmt.Errorf("wait listener to be ready failed. %v", err)
	}
	klog.Info("Listener created. ListenerId: ", *rsp.Response.ListenerIds[0])
	return rsp.Response.ListenerIds[0], err
}

func (p *provider) RegisterTargets(lbId, listenerID *string, eniIPs []*string) error {
	klog.Info("Register targets")
	listener, err := p.DescribeListeners(lbId)
	if err != nil {
		return err
	}

	req := clb.NewRegisterTargetsRequest()
	req.LoadBalancerId = lbId
	req.ListenerId = listenerID
	req.Targets = []*clb.Target{}
	for _, eniIP := range eniIPs {
		req.Targets = append(req.Targets, &clb.Target{
			Port:   listener.Port,
			Weight: common.Int64Ptr(10),
			EniIp:  eniIP,
		})
	}
	rsp, err := p.clbClient.RegisterTargets(req)
	if err != nil {
		return err
	}
	err = wait.Poll(time.Second, time.Minute, func() (bool, error) {
		req1 := clb.NewDescribeTaskStatusRequest()
		req1.TaskId = rsp.Response.RequestId
		rsp1, err := p.clbClient.DescribeTaskStatus(req1)
		return *rsp1.Response.Status == 0, err
	})
	if err != nil {
		return fmt.Errorf("wait register targets success failed. %v", err)
	}
	klog.Info("Targets registered")
	return err
}

func (p *provider) DescribeListeners(lbId *string) (*clb.Listener, error) {
	req := clb.NewDescribeListenersRequest()
	req.LoadBalancerId = lbId
	rsp, err := p.clbClient.DescribeListeners(req)
	if err != nil {
		return nil, fmt.Errorf("DescribeListeners failed. %v", err)
	}
	return rsp.Response.Listeners[0], err
}

func (p *provider) TearDown() error {
	err := p.DeleteLBs(common.StringPtrs(p.lbIds))
	if err != nil {
		klog.Error(err)
	}
	err = p.DeleteInstances(p.instanceIds)
	if err != nil {
		klog.Error(err)
	}
	return err
}
