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

import React, { useState, useEffect } from 'react';
import { AntdLayout } from '@src/modules/common/layouts';
import { Button, H5, Form, Select, Checkbox, InputNumber } from 'tea-component';
import { getK8sValidVersions } from '@src/webApi/cluster';
import { compareVersion } from '@helper/version';
import { RootProps } from '../ClusterApp';
import { updateCluster } from '@src/webApi/cluster';
import { useForm, Controller } from 'react-hook-form';
import { getReactHookFormStatusWithMessage } from '@helper';

export function ClusterUpdate({ route, actions }: RootProps) {
  const defaultUpgradeConfig = {
    version: null,
    drainNodeBeforeUpgrade: true,
    maxUnready: 20,
    autoMode: true
  };

  const { handleSubmit, control, watch } = useForm({
    mode: 'onBlur',
    defaultValues: defaultUpgradeConfig
  });

  const { clusterId, clusterVersion } = route.queries;

  function goBack() {
    history.back();
  }

  async function perform(values) {
    await updateCluster({ ...values, clusterName: clusterId });

    goBack();
  }

  console.log(React.version);

  const [k8sValidVersions, setK8sValidVersions] = useState([]);

  useEffect(() => {
    async function fetchK8sValidVersions() {
      const versions = await getK8sValidVersions();
      const k8sValidVersions = versions.map(v => ({
        value: v,
        disabled: compareVersion(clusterVersion, v) >= 0
      }));
      setK8sValidVersions(k8sValidVersions);
    }

    fetchK8sValidVersions();
  }, []);

  return (
    <AntdLayout
      title="升级Master"
      footer={
        <>
          <Button type="primary" style={{ marginRight: 10 }} onClick={handleSubmit(perform)}>
            提交
          </Button>

          <Button onClick={goBack}>取消</Button>
        </>
      }
    >
      <H5 style={{ marginBottom: 35 }}>更新集群的K8S版本{clusterVersion}至：</H5>

      <Form>
        <Controller
          control={control}
          name="version"
          rules={{ required: '请选择将要升级的k8s版本！' }}
          render={({ field, ...others }) => (
            <Form.Item
              label="升级目标版本"
              extra="注意：master升级支持一个次版本升级到下一个次版本，或者同样次版本的补丁版。"
              {...getReactHookFormStatusWithMessage(others)}
            >
              <Select {...field} options={k8sValidVersions} />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="autoMode"
          render={({ field }) => (
            <Form.Item
              label="自动升级Worker"
              extra="注意：启用自定升级Worker，在升级完Master后，将自动升级集群下所有Worker节点。"
            >
              <Checkbox {...field}>启用自动升级</Checkbox>
            </Form.Item>
          )}
        />

        {watch('autoMode') && (
          <>
            <Controller
              control={control}
              name="drainNodeBeforeUpgrade"
              render={({ field }) => (
                <Form.Item
                  label="驱逐节点"
                  extra="若选择升级前驱逐节点，该节点所有pod将在升级前被驱逐，此时节点如有pod使用emptyDir类卷会导致驱逐失败而影响升级流程"
                >
                  <Checkbox {...field}>驱逐节点</Checkbox>
                </Form.Item>
              )}
            />

            <Controller
              control={control}
              name="maxUnready"
              render={({ field }) => (
                <Form.Item
                  label="最大不可用Pod占比"
                  extra="注意如果节点过少，而设置比例过低，没有足够多的节点承载pod的迁移会导致升级卡死。如果业务对pod可用比例较高，请考虑选择升级前不驱逐节点。"
                >
                  <InputNumber {...field} min={0} max={100} unit="%" />
                </Form.Item>
              )}
            />
          </>
        )}
      </Form>
    </AntdLayout>
  );
}
