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
import { dataFormatConfig } from './dataFormat';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** container所在的字段 */
const containerMapPath = {
  pod: 'spec.containers',
  basic: 'spec.template.spec.containers',
  cronjob: 'spec.jobTemplate.spec.template.spec.containers'
};

/** volumnes所在的k8s的字段 */
const volumesMapPath = {
  pod: 'spec.volumes',
  basic: 'spec.template.spec.volumes',
  cronjob: 'spec.template.spec.volumes'
};

/** 这里是 资源详情页的 基础信息的配置 */
export const commonDetailInfo = (path: string = 'basic') => ({
  volume: {
    volumes: {
      dataField: [volumesMapPath[path]]
    }
  },
  container: {
    containers: {
      dataField: [containerMapPath[path]],
      displayField: {
        name: {
          dataField: ['name'],
          dataFormat: dataFormatConfig['text'],
          label: t('容器名称'),
          tips: ''
        },
        image: {
          dataField: ['image'],
          dataFormat: dataFormatConfig['text'],
          label: t('镜像'),
          tips: ''
        },
        imagePullPolicy: {
          dataField: ['imagePullPolicy'],
          dataFormat: dataFormatConfig['text'],
          label: t('镜像更新策略'),
          tips: ''
        },
        requestCpu: {
          dataField: ['resources.requests.cpu'],
          dataFormat: dataFormatConfig['text'],
          label: 'CPU Requested',
          tips: ''
        },
        limitsCpu: {
          dataField: ['resources.limits.cpu'],
          dataFormat: dataFormatConfig['text'],
          label: 'CPU Limited',
          tips: ''
        },
        requestMem: {
          dataField: ['resources.requests.memory'],
          dataFormat: dataFormatConfig['text'],
          label: t('内存 Requested'),
          tips: ''
        },
        limitsMem: {
          dataField: ['resources.limits.memory'],
          dataFormat: dataFormatConfig['text'],
          label: t('内存 Limited'),
          tips: ''
        },
        workingDir: {
          dataField: ['workingDir'],
          dataFormat: dataFormatConfig['text'],
          label: t('工作目录'),
          tips: ''
        },
        command: {
          dataField: ['command'],
          dataFormat: dataFormatConfig['array'],
          label: t('运行命令'),
          tips: ''
        },
        args: {
          dataField: ['args'],
          dataFormat: dataFormatConfig['array'],
          label: t('运行参数'),
          tips: ''
        },
        env: {
          dataField: ['env'],
          dataFormat: dataFormatConfig['env'],
          label: t('环境变量'),
          tips: ''
        },
        volumeMounts: {
          dataField: ['volumeMounts'],
          dataFormat: dataFormatConfig['volume'],
          label: t('挂载点'),
          tips: ''
        },
        livenessProbe: {
          dataField: ['livenessProbe'],
          dataFormat: dataFormatConfig['probe'],
          label: t('存活检查'),
          tips: ''
        },
        readinessProbe: {
          dataField: ['readinessProbe'],
          dataFormat: dataFormatConfig['probe'],
          label: t('就绪检查'),
          tips: ''
        }
      }
    }
  }
});

/** k8s对象tabList */
export const workloadCommonTabList = [
  {
    id: 'pod',
    label: t('Pod管理')
  },
  {
    id: 'event',
    label: t('事件')
  },
  {
    id: 'log',
    label: t('日志')
  },
  {
    id: 'info',
    label: t('详情')
  },
  {
    id: 'yaml',
    label: 'YAML'
  }
];

export const commonTabList = [
  {
    id: 'info',
    label: t('详情')
  },
  {
    id: 'event',
    label: t('事件')
  },
  {
    id: 'yaml',
    label: 'YAML'
  }
];

/** 默认的 notExistValue */
export const defaulNotExistedValue = '-';
