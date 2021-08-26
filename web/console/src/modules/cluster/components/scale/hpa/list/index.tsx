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

import React, { useState, useEffect, useContext, useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
  Table,
  Justify,
  Button,
  Card,
  Layout,
  Pagination,
  Text,
  TagSearchBox,
  Modal,
  AttributeValue
} from 'tea-component';
import { RecordSet } from '@tencent/ff-redux/src';
// import { router } from '@src/modules/cluster/router.project';
import { router } from '@src/modules/cluster/router';
import { downloadCsv } from '@helper/downloadCsv';
import { DisplayFiledProps, Resource, ResourceInfo } from '@src/modules/common/models';
import { NestedMetricsResourceMap, MetricsResourceMap } from '../constant';
import { CHANGE_NAMESPACE, StateContext, DispatchContext } from '../context';
import { isEmpty, useModal } from '@src/modules/common/utils';
import { removeHPA } from '@src/modules/cluster/WebAPI/scale';
import { LinkButton, Clip } from '@src/modules/common/components';
import NamespaceSelect from './namespaceSelect';

const { Body, Content } = Layout;
const { autotip } = Table.addons;

interface ListProps {
  namespaces: RecordSet<Resource>;
  hpaData: RecordSet<Resource>;
  triggerRefresh: () => void;
}

const List = React.memo((props: ListProps) => {
  const route = useSelector(state => state.route);
  const urlParams = router.resolve(route);
  const { clusterId } = route.queries;
  const { namespaces, hpaData, triggerRefresh } = props;
  const { namespaceValue } = useContext(StateContext);
  const { isShowing, toggle } = useModal();
  const [removeHpaName, setRemoveHpaName] = useState();

  /**
   * 列表内容处理
   */
  const [loading, setLoading] = useState(true);
  const [emptyText, setEmptyText] = useState('');
  useEffect(() => {
    if (isEmpty(hpaData)) {
      setLoading(true);
      setEmptyText('');
    } else if (!hpaData.recordCount) {
      setLoading(false);
      setEmptyText('您选择的该资源的列表为空，您可以切换到其他命名空间');
    } else if (hpaData.recordCount) {
      setLoading(false);
      setEmptyText('');
    }
  }, [hpaData]);

  /**
   * TagSearchBox和list数据的管理
   */
  const attributes: AttributeValue[] = [
    {
      type: 'input',
      key: 'resourceName',
      name: t('名称')
    }
  ];
  const [hpaList, setHpaList] = useState([]);
  const [tags, setTags] = useState([]);
  useEffect(() => {
    if (!isEmpty(hpaData)) {
      if (tags.length) {
        const newHpaList = [];
        hpaData.records.forEach(item => {
          if (item.metadata.name.indexOf(tags[0].values[0].name) !== -1) {
            newHpaList.push(item);
          }
        });
        setHpaList(newHpaList);
      } else {
        setHpaList(hpaData.records);
      }
    }
  }, [hpaData, tags]);

  /**
   * 下载
   */
  const downloadHandle = useCallback(resourceList => {
    const head = ['名称', '关联工作负载', '触发策略', '最小实例数', '最大实例数'];
    function getTriggerStrategys(metrics) {
      return metrics.map((item, index) => {
        const { name, targetAverageValue, targetAverageUtilization } = item.resource;
        const { meaning, unit } = MetricsResourceMap[name];
        const content = targetAverageValue ? meaning + targetAverageValue : meaning + targetAverageUtilization + unit;
        return content;
      });
    }
    const rows = resourceList.map(hpa => {
      const triggerStrategyArr = getTriggerStrategys(hpa.spec.metrics);
      const triggerStrategyStr = isEmpty(triggerStrategyArr) ? '-' : triggerStrategyArr.join(';');
      return [
        hpa.metadata.name,
        `${hpa.spec.scaleTargetRef.kind}:${hpa.spec.scaleTargetRef.name}`,
        triggerStrategyStr,
        hpa.spec.minReplicas,
        hpa.spec.maxReplicas
      ];
    });
    downloadCsv(rows, head, 'tke_hpa_' + new Date().getTime() + '.csv');
  }, []);

  return (
    <Layout>
      <Body>
        <Content>
          <Content.Body>
            <Table.ActionPanel>
              <Justify
                left={
                  <Button
                    type="primary"
                    onClick={() => {
                      delete route.queries.HPAName;
                      router.navigate({ ...urlParams, mode: 'create' }, route.queries);
                    }}
                  >
                    <Trans>新建</Trans>
                  </Button>
                }
                right={
                  <>
                    <NamespaceSelect namespaces={namespaces} />
                    <div style={{ width: 350, display: 'inline-block' }}>
                      <TagSearchBox
                        className="myTagSearchBox"
                        attributes={attributes}
                        value={tags}
                        onChange={tags => {
                          setTags(tags);
                        }}
                      />
                    </div>
                    <Button
                      onClick={() => {
                        triggerRefresh();
                      }}
                      icon="refresh"
                    />
                    <Button
                      onClick={() => {
                        downloadHandle(hpaList);
                      }}
                      icon="download"
                    />
                  </>
                }
              />
            </Table.ActionPanel>
            <Card>
              <Table
                records={hpaList}
                recordKey="id"
                columns={[
                  {
                    key: 'name',
                    header: t('名称'),
                    render: hpa => (
                      <>
                        <LinkButton
                          onClick={() => {
                            router.navigate(
                              { ...urlParams, mode: 'detail' },
                              { ...route.queries, namespaceValue, HPAName: hpa.metadata.name }
                            );
                          }}
                        >
                          <Text id={'hpaName' + hpa.id}>{hpa.metadata.name}</Text>
                        </LinkButton>
                        <Clip target={'#hpaName' + hpa.id} />
                      </>
                    )
                  },
                  {
                    key: 'workload',
                    header: t('关联工作负载'),
                    render: hpa => (
                      <>
                        <Text id={'hpaWorkload' + hpa.id}>
                          {hpa.spec.scaleTargetRef.kind}:{hpa.spec.scaleTargetRef.name}
                        </Text>
                        <Clip target={'#hpaWorkload' + hpa.id} />
                      </>
                    )
                  },
                  {
                    key: 'triggerStrategy',
                    header: t('触发策略'),
                    render: hpa => {
                      return hpa.spec.metrics.map((item, index) => {
                        const { name, targetAverageValue, targetAverageUtilization } = item.resource;
                        const target = targetAverageValue ? 'targetAverageValue' : 'targetAverageUtilization';
                        let content = '';
                        if (name === 'cpu' || name === 'memory') {
                          const target = targetAverageValue ? 'targetAverageValue' : 'targetAverageUtilization';
                          const { meaning, unit } = NestedMetricsResourceMap[name][target];
                          content =
                            (targetAverageValue ? meaning + targetAverageValue : meaning + targetAverageUtilization) +
                            unit;
                        } else {
                          const { meaning, unit } = MetricsResourceMap[name];
                          content =
                            (targetAverageValue ? meaning + targetAverageValue : meaning + targetAverageUtilization) +
                            unit;
                        }
                        return (
                          <Text key={index} parent="div">
                            {content}
                          </Text>
                        );
                      });
                    }
                  },
                  {
                    key: 'min',
                    header: t('最小实例数'),
                    width: '10%',
                    render: hpa => <Text>{hpa.spec.minReplicas}</Text>
                  },
                  {
                    key: 'max',
                    header: t('最大实例数'),
                    width: '10%',
                    render: hpa => <Text>{hpa.spec.maxReplicas}</Text>
                  },
                  {
                    key: 'action',
                    header: t('操作'),
                    render: hpa => (
                      <>
                        <LinkButton
                          tipDirection="left"
                          onClick={() => {
                            router.navigate(
                              { ...urlParams, mode: 'modify-hpa' },
                              { ...route.queries, namespaceValue, HPAName: hpa.metadata.name }
                            );
                          }}
                        >
                          <Trans>修改配置</Trans>
                        </LinkButton>
                        <LinkButton
                          tipDirection="left"
                          onClick={() => {
                            router.navigate(
                              { ...urlParams, mode: 'modify-yaml' },
                              { ...route.queries, namespaceValue, HPAName: hpa.metadata.name }
                            );
                          }}
                        >
                          <Trans>编辑YAML</Trans>
                        </LinkButton>
                        <LinkButton
                          tipDirection="left"
                          onClick={() => {
                            setRemoveHpaName(hpa.metadata.name);
                            toggle();
                          }}
                        >
                          <Trans>删除</Trans>
                        </LinkButton>
                      </>
                    )
                  }
                ]}
                addons={[
                  autotip({
                    isLoading: loading,
                    emptyText: emptyText
                  })
                ]}
              />
              <Pagination recordCount={hpaList.length} pageSizeVisible={false} />
            </Card>
            <Modal visible={isShowing} caption="删除资源" onClose={toggle}>
              <Modal.Body>
                <Trans>您确定要删除HorizontalPodAutoscaler: {{ removeHpaName }}吗？</Trans>
              </Modal.Body>
              <Modal.Footer>
                <Button
                  type="primary"
                  onClick={async () => {
                    const isRemove = await removeHPA({ namespace: namespaceValue, clusterId, name: removeHpaName });
                    if (isRemove) {
                      toggle();
                      triggerRefresh();
                    }
                  }}
                >
                  确定
                </Button>
                <Button type="weak" onClick={toggle}>
                  取消
                </Button>
              </Modal.Footer>
            </Modal>
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
});

export default List;
