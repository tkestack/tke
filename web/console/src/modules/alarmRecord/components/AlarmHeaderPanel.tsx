import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { FormPanelSelect, FormPanelSelectProps } from '@tencent/ff-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux/libs/qcloud-lib';
import { allActions } from '@src/modules/alarmRecord/actions';
import { router } from '@src/modules/alarmRecord/router';
const { useState, useEffect } = React;

export const AlarmRecordHeadPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { cluster, route } = state;
  const clusterData = cluster.list.data;
  const urlParams = router.resolve(route);
  const queryClusterId = route.queries.clusterId;

  // 获取集群列表
  useEffect(() => {
    actions.cluster.applyFilter({});
  }, []);

  // 初始化选择
  useEffect(() => {
    if (clusterData.recordCount > 0) {
      // 不带clusterID 刷新页面
      let selectedCluster = clusterData.records[0];

      // 带clusterID 刷新页面
      if (queryClusterId) {
        selectedCluster = clusterData.records.find(c => c.metadata.name === queryClusterId);
      }

      actions.cluster.select(selectedCluster);
      router.navigate(urlParams, { ...route.queries, clusterId: selectedCluster.metadata.name });
    }
  }, [clusterData.recordCount > 0]);

  let selectProps: FormPanelSelectProps = {
    type: 'simulate',
    appearence: 'button',
    label: '集群',
    model: cluster,
    action: actions.cluster,
    valueField: record => record.metadata.name,
    displayField: record => `${record.metadata.name} (${record.spec.displayName})`,
    onChange: (clusterId: string) => {
      const selectedCluster = clusterData.records.find(c => c.metadata.name === clusterId);
      actions.cluster.select(selectedCluster);
      router.navigate(urlParams, { ...route.queries, clusterId: selectedCluster.metadata.name });
    }
  };

  return (
    <>
      <div className="tc-15-dropdown" style={{ marginLeft: '20px', display: 'inline-block', minWidth: '30px' }}>
        {t('集群')}
      </div>
      <FormPanelSelect {...selectProps} />
    </>
  );
};
