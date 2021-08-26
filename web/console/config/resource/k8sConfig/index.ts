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

export { deployment } from './deployment';
export { statefulset } from './statefulset';
export { daemonset } from './daemonset';
export { jobs } from './jobs';
export { cronjobs } from './cronjobs';
export { tapps } from './tapp';
export { pods } from './pods';
export { rc } from './replicationcontrollers';
export { rs } from './replicaset';
export { svc } from './service';
export { ingress } from './ingress';
export { np } from './namespace';
export { configmap } from './configmaps';
export { secret } from './secret';
export { pv } from './persistentvolumes';
export { pvc } from './persistentvolumeclaims';
export { sc } from './storageclass';
export * from './otherResource';
export { node } from './node';
export * from './addonResource';
export * from './alarmPolicy';
export * from './alarmRecord';
export * from './notifyChannel';
export * from './audit';
export * from './application';

export { lbcf, lbcf_bg, lbcf_br, lbcf_driver } from './lbcf';

export { serviceForMesh } from './serviceForMesh';
export { gateway } from './gateway';
export { controlPlane } from './controlPlane';
export { virtualService } from './virtualService';
export { destinationRule } from './destinationRule';
