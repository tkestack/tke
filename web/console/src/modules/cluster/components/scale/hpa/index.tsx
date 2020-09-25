import React, { useState, useEffect, useContext, useRef, useCallback } from 'react';
import { useSelector } from 'react-redux';
import { HpaScopeProvider, DispatchContext, StateContext } from './context';
import { isEmpty, useRefresh } from '@src/modules/common/utils';
import { cutNsStartClusterId } from '@helper';
import { router } from '@src/modules/cluster/router';
import List from './list';
import Detail from './detail';
import Hpa from  './editor/hpa';
import Yaml from './editor/yaml';
import { getHPAList } from '@src/modules/cluster/WebAPI/scale';
import { useNamespaces } from '../common/hooks';

/**
 * HPAPanel组件，带有scope的局部全局数据的组件
 */
export const HPAPanel = React.memo(() => {
  return (
    <HpaScopeProvider>
      <HPA />
    </HpaScopeProvider>
  );
});

/**
 * HPA组件
 */
export const HPA = React.memo(() => {
  const route = useSelector((state) => state.route);
  const { projectName, clusterId, HPAName, namespaceValue: namespaceValueInURL } = route.queries;
  const { mode } = router.resolve(route);
  const hpaState = useContext(StateContext);
  const hpaDispatch = useContext(DispatchContext);
  const { namespaceValue } = hpaState;

  // 刷新标识
  const { refreshFlag, triggerRefresh } = useRefresh();

  // 获取命名空间数据
  const namespaces = useNamespaces({ projectId: projectName, clusterId });

  /**
   * HPA列表数据获取
   */
  const [HPAData, setHPAData] = useState();
  useEffect(() => {
    async function getHPAData(namespace) {
      const result = await getHPAList({ namespace, clusterId });
      setHPAData(result);
    }
    // namespaceValue处理的逻辑是列表页的，namespaceValueInURL的逻辑是详情页的
    if ((namespaceValue || namespaceValueInURL) && clusterId) {
      const newNamespaceValue = namespaceValue || namespaceValueInURL;
      const namespace = cutNsStartClusterId({ namespace: newNamespaceValue, clusterId });
      getHPAData(namespace);
    }
  }, [namespaceValue, clusterId, refreshFlag, namespaceValueInURL]);

  /**
   * 根据URL中的参数获取具体某个HPA
   */
  const [selectedHpa, setSelectedHpa] = useState({});
  useEffect(() => {
    if (!isEmpty(HPAData)) {
      HPAData.records.forEach(hpa => {
            if (hpa.metadata.name === HPAName) {
              setSelectedHpa(hpa);
            }
          }
      );
    }
  }, [HPAData]);

  /**
   * 页面内容
   */
  let content: React.ReactNode;
  switch (mode) {
    case 'list':
      content = <List namespaces={namespaces} hpaData={HPAData} triggerRefresh={triggerRefresh} />;
      break;
    case 'create':
    case 'modify-hpa':
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
  return (
    <>
      {content}
    </>
  );
});
