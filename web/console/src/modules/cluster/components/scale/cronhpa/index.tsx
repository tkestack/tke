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

import React, { useState, useEffect, useContext, useRef, useCallback, useMemo } from 'react';
import { useSelector } from 'react-redux';
import { HpaScopeProvider, DispatchContext, StateContext } from './context';
import { isEmpty, useRefresh, usePrevious } from '@src/modules/common/utils';
import { cutNsStartClusterId } from '@helper';
import { router } from '@src/modules/cluster/router';
import List from './list';
import Detail from './detail';
import Hpa from './editor/hpa';
import Yaml from './editor/yaml';
import {
  fetchProjectNamespaceList,
  fetchCronHpaRecords,
  fetchAddons,
  fetchNamespaceList
} from '@src/modules/cluster/WebAPI/scale';
import { useNamespaces } from '@src/modules/cluster/components/scale/common/hooks';
import { RecordSet } from '@tencent/ff-redux';
import { Resource } from '@src/modules/common';

/**
 * HPAPanel组件，带有scope的局部全局数据的组件
 */
export const CronHpaPanel = React.memo(() => {
  return (
    <HpaScopeProvider>
      <CronHpa />
    </HpaScopeProvider>
  );
});

/**
 * HPA组件
 */
export const CronHpa = React.memo(() => {
  const { route } = useSelector(state => ({ route: state.route }));
  const { projectName, clusterId, HPAName, namespaceValue: namespaceValueInURL, np } = route.queries;
  const { mode } = router.resolve(route);
  const hpaState = useContext(StateContext);
  const hpaDispatch = useContext(DispatchContext);
  const { namespaceValue } = hpaState;
  // 刷新标识
  const { refreshFlag, triggerRefresh } = useRefresh();

  // 获取命名空间数据
  const namespaces = useNamespaces({ projectId: projectName, clusterId });

  /**
   * 是否安装CronHpa组件
   * HPA列表数据获取
   */
  const [isCronHpaInstalled, setIsCronHpaInstalled] = useState<boolean | undefined>();
  const [cronHpaRecords, setCronHpaRecords] = useState<RecordSet<Resource, any> | undefined>();
  const previousClusterId = usePrevious(clusterId);
  useEffect(() => {
    // 切换clusterId后，fetchAddons 和 fetchCronHpaRecords的请求是有先后顺序限制的，这样能分开在两个useEffect中吗？？？
    async function getCronHpaRecords() {
      let hasCronHap = isCronHpaInstalled;
      // clusterId !== previousClusterId节省调用次数，!projectName--平台侧
      if (clusterId && (clusterId !== previousClusterId || !projectName)) {
        const addons = await fetchAddons({ clusterId });
        hasCronHap = addons && addons['CronHPA'] ? true : false;
        setIsCronHpaInstalled(hasCronHap);
      }

      // namespaceValue处理的逻辑是列表页的，namespaceValueInURL的逻辑是详情页的，np处理的是平台侧的
      const newNamespaceValue = namespaceValue || namespaceValueInURL || np;

      if (hasCronHap && newNamespaceValue && clusterId) {
        const namespace = cutNsStartClusterId({ namespace: newNamespaceValue, clusterId });
        const result = await fetchCronHpaRecords({ namespace, clusterId });
        setCronHpaRecords(result);
      }
    }

    // 设置isCronHpaInstalled为初始状态
    if (clusterId !== previousClusterId) {
      setIsCronHpaInstalled(undefined);
    }

    // if (clusterId) {
    //   getAddons(clusterId);
    // }
    // namespaceValue处理的逻辑是列表页的，namespaceValueInURL的逻辑是详情页的，np处理的是平台侧的
    setCronHpaRecords({ recordCount: 0, records: [] });
    getCronHpaRecords();
    // if ((namespaceValue || namespaceValueInURL || np) && clusterId) {
    //   const newNamespaceValue = namespaceValue || namespaceValueInURL || np;
    //   const namespace = cutNsStartClusterId({ namespace: newNamespaceValue, clusterId });
    //   getCronHpaRecords(namespace, clusterId);
    // }
  }, [namespaceValue, clusterId, refreshFlag, namespaceValueInURL, np]);

  /**
   * 根据URL中的参数获取具体某个HPA
   */
  const [selectedHpa, setSelectedHpa] = useState({});
  useEffect(() => {
    if (!isEmpty(cronHpaRecords)) {
      cronHpaRecords.records.forEach(hpa => {
        if (hpa.metadata.name === HPAName) {
          setSelectedHpa(hpa);
        }
      });
    }
  }, [cronHpaRecords]);

  /**
   * 页面内容
   */
  let content: React.ReactNode;
  switch (mode) {
    case 'list':
      content = (
        <List
          namespaces={namespaces}
          cronHpaData={cronHpaRecords}
          triggerRefresh={triggerRefresh}
          isCronHpaInstalled={isCronHpaInstalled}
        />
      );
      break;
    case 'create':
    case 'modify-cronhpa':
      content = <Hpa selectedHpa={selectedHpa} />;
      break;
    case 'modify-yaml':
      content = <Yaml />;
      break;
    case 'detail':
      content = <Detail selectedHpa={selectedHpa} />;
      break;
    default:
      content = '';
  }
  return <>{content}</>;
});
