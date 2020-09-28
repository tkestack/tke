/**
 * namespace下拉选择组件
 */
import React, { useState, useEffect, useContext, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { CHANGE_NAMESPACE, StateContext, DispatchContext } from '../context';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { isEmpty } from '@src/modules/common/utils';
// import { router } from '@src/modules/cluster/router.project';
import { router } from '@src/modules/cluster/router';
import {
  Select,
  Text
} from '@tencent/tea-component';
import { RecordSet } from '@tencent/ff-redux/src';
import { Resource } from '@src/modules/common/models';

// 下边props没有用这个interface因为使用后Select会你报一些类型的问题，感觉还不太好整合，有时间整合下
interface NamespaceSelectProps {
  namespaces: RecordSet<Resource>;
}
let i = 0;
const NamespaceSelect = React.memo((props: {
  namespaces: any;
}) => {
  const route = useSelector((state) => state.route);
  const urlParams = router.resolve(route);
  const { clusterId, projectName } = route.queries;
  const { namespaces = {}} = props;
  const hpaDispatch = useContext(DispatchContext);
  const hpaState = useContext(StateContext);
  const { namespaceValue } = hpaState;

  /**
   * 初始化namespace选项为第一个，对应初始化浏览器URL数据
   */
  useEffect(() => {
    if (!isEmpty(namespaces) && namespaces.recordCount) {
      hpaDispatch({
        type: CHANGE_NAMESPACE,
        payload: { namespaceValue: namespaces.records[0].value }
      });

      // 如果是业务侧，
      if (projectName && !namespaces.records[0].spec.clusterName) {
        router.navigate(urlParams, { ...route.queries, clusterId: namespaces.records[0].spec.clusterName });
      }
    }
  }, [namespaces]);

  const selectStyle = useMemo(() => ({ display: 'inline-block', fontSize: '12px', verticalAlign: 'middle' }), []);
  return (
    <div style={selectStyle}>
      <Text theme="label" verticalAlign="middle">
        {t('命名空间')}
      </Text>
      <Select
        type="native"
        appearence="button"
        size="s"
        options={namespaces.records}
        style={{ width: '130px', marginRight: '5px' }}
        value={namespaceValue}
        onChange={selectedNamespace => {
          hpaDispatch({
            type: CHANGE_NAMESPACE,
            payload: { namespaceValue: selectedNamespace }
          });

          // namespace选项改变时对应改变路由中的参数
          if (projectName) {
            namespaces.records.forEach(item => {
              if (item.value === selectedNamespace) {
                  router.navigate(urlParams, { ...route.queries, clusterId: item.spec.clusterName });
              }
            });
          } else {
            router.navigate(urlParams, { ...route.queries, np: selectedNamespace });
          }
        }}
        placeholder={namespaces.recordCount ? '' : t('无可用命名空间')}
      />
    </div>
  );
});

export default NamespaceSelect;
