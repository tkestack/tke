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

import { generateResourceInfo } from '../common';

/** addon的相关配置 */
export const addon = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon',
    requestType: {
      list: 'clusteraddontypes'
    }
  });
};

/** helm的相关配置 */
export const addon_helm = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_helm',
    requestType: {
      list: 'helms'
    }
  });
};

/** gpumanage的相关配置 */
export const addon_gpumanager = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_gpumanager',
    requestType: {
      list: 'gpumanagers'
    }
  });
};

/** logCollector的相关配置 */
export const addon_logcollector = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_logcollector',
    requestType: {
      list: 'logcollectors'
    }
  });
};

/** tappcontroller的相关配置 */
export const addon_tappcontroller = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_tappcontroller',
    requestType: {
      list: 'tappcontrollers'
    }
  });
};

/** csioperator的相关配置 */
export const addon_csioperator = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_csioperator',
    requestType: {
      list: 'csioperators'
    }
  });
};

/** lbcf的相关配置 */
export const addon_lbcf = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_lbcf',
    requestType: {
      list: 'lbcfs'
    }
  });
};

/** cronhpa的相关配置 */
export const addon_cronhpa = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_cronhpa',
    requestType: {
      list: 'cronhpas'
    }
  });
};

/** coredns的相关配置 */
export const addon_coredns = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_coredns',
    requestType: {
      list: 'corednss'
    }
  });
};

/** galaxy的相关配置 */
export const addon_galaxy = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_galaxy',
    requestType: {
      list: 'galaxies'
    }
  });
};

/** Prometheus的相关配置 */
export const addon_prometheus = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_prometheus',
    requestType: {
      list: 'prometheuses'
    }
  });
};

/** VolumeDecorator的相关配置 */
export const addon_volumedecorator = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_volumedecorator',
    requestType: {
      list: 'volumedecorators'
    }
  });
};

/** IPAM的相关配置 */
export const addon_ipam = (k8sVersion: string) => {
  return generateResourceInfo({
    k8sVersion,
    resourceName: 'addon_ipam',
    requestType: {
      list: 'ipams'
    }
  });
};
