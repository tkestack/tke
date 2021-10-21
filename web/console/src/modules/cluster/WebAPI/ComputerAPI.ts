/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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
import { QueryState, RecordSet } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../config';
import {
    Method, operationResult, reduceNetworkRequest, reduceNetworkWorkflow
} from '../../../../helpers';
import { RequestParams } from '../../common/models';
import { Computer, ComputerFilter } from '../models';
import { ComputerLabelEdition, ComputerOperator, ComputerTaintEdition } from '../models/Computer';

// 提示框
const tips = seajs.require('tips');

/**
 * 更新节点label
 * @param computerLabelEdition: ComputerLabelEdition[] 节点编辑的集合
 * @param opt: any  操作的集合
 */
export async function updateComputerLabel(computerLabelEdition: ComputerLabelEdition[], opt: any) {
  try {
    let resourceInfo = resourceConfig().node;
    let url = `/${resourceInfo.basicEntry}/${resourceInfo.version}/nodes/${computerLabelEdition[0].computerName}`;
    let labels = computerLabelEdition[0].labels.reduce((prev, next) => {
      return Object.assign({}, prev, {
        [next.key]: next.value
      });
    }, {});
    Object.keys(computerLabelEdition[0].originLabel).forEach(key => {
      if (labels[key] === undefined) {
        labels[key] = null;
      }
    });
    let jsonData = {
      metadata: { labels }
    };
    // 去除当中不需要的数据
    jsonData = JSON.parse(JSON.stringify(jsonData));
    let params: RequestParams = {
      method: Method.patch,
      url,
      userDefinedHeader: {
        'Content-Type': 'application/strategic-merge-patch+json'
      },
      data: JSON.stringify(jsonData)
    };

    let response = await reduceNetworkRequest(params, opt.clusterId);

    if (response.code === 0) {
      tips.success(t('更新成功'), 2000);
      return operationResult(computerLabelEdition);
    } else {
      return operationResult(computerLabelEdition, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(computerLabelEdition, reduceNetworkWorkflow(error));
  }
}

/**
 * 批量封锁 drain 指定的 computer
 * @param {Computer[]} computers 要 UnSchedule的节点
 * @param operator:clusterId | regionId | isDrainComputer  附带的一些操作
 */
export async function drainComputer(computers: Computer[], operator: ComputerOperator) {
  try {
    let { clusterId } = operator,
      resourceInfo = resourceConfig().cluster,
      url = `/${resourceInfo.basicEntry}/${resourceInfo.group}/${resourceInfo.version}/clusters/${clusterId}/drain/${computers[0].metadata.name}`;
    let params: RequestParams = {
      method: Method.post,
      url
    };
    let response = await reduceNetworkRequest(params, clusterId);
    if (response.code === 0) {
      return operationResult(computers);
    } else {
      return operationResult(computers, response);
    }
  } catch (error) {
    return operationResult(computers, error);
  }
}

export async function updateComputerTaints(computerTaintEdition: ComputerTaintEdition[], opt: any) {
  try {
    let resourceInfo = resourceConfig().node;
    let url = `/${resourceInfo.basicEntry}/${resourceInfo.version}/nodes/${computerTaintEdition[0].computerName}`;
    // 构建更新ingress 转发配置的json格式
    let jsonData = {
      spec: {
        taints: computerTaintEdition[0].taints.map(item => ({
          key: item.key,
          value: item.value,
          effect: item.effect
        }))
      }
    };
    // 去除当中不需要的数据
    jsonData = JSON.parse(JSON.stringify(jsonData));
    let params: RequestParams = {
      method: Method.patch,
      url,
      userDefinedHeader: {
        'Content-Type': 'application/merge-patch+json'
      },
      data: JSON.stringify(jsonData)
    };

    let response = await reduceNetworkRequest(params, opt.clusterId);

    if (response.code === 0) {
      tips.success(t('更新成功'), 2000);
      return operationResult(computerTaintEdition);
    } else {
      return operationResult(computerTaintEdition, reduceNetworkWorkflow(response));
    }
  } catch (error) {
    return operationResult(computerTaintEdition, reduceNetworkWorkflow(error));
  }
}
