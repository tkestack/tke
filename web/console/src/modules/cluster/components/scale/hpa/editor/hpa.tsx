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
/**
 * HPA创建及修改组件
 */
import React, { useState, useEffect, useContext, useMemo, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import {
  Layout,
  Card,
  Select,
  Text,
  Button,
  Form,
  Input,
  InputNumber
} from '@tencent/tea-component';
import { LinkButton } from '@src/modules/common/components';
import { Resource } from '@src/modules/common/models';
import { useForm, useFieldArray, Controller, NestedValue } from 'react-hook-form';
import { isEmpty } from '@src/modules/common/utils';
import { insertCSS, uuid } from '@tencent/ff-redux/libs/qcloud-lib';
// import { router } from '@src/modules/cluster/router.project';
import { router } from '@src/modules/cluster/router';
import { cutNsStartClusterId } from '@helper';
import {
  fetchNamespaceList,
  createHPA,
  modifyHPA,
  fetchResourceList,
  fetchProjectNamespaceList
} from '@src/modules/cluster/WebAPI/scale';
import { MetricsResourceMap, NestedMetricsResourceMap } from '../constant';
import { useNamespaces } from '@src/modules/cluster/components/scale/common/hooks';
import { RecordSet } from '@tencent/ff-redux';

/**
 * 组件样式
 */
insertCSS(
    'hpa-editor-panel',
    `
      .hpa-edit-strategy-ul { margin-bottom : 10px; }
      .hpa-edit-strategy-li + .hpa-edit-strategy-li { margin-top: 5px; }
    `
);

/**
 * 常量定义
 */
const { Body, Content } = Layout;

// url中路由字段
const ModifyHpa = 'modify-hpa';

// 触发策略select选项
const StrategyOptions = [
  {
    id: uuid(),
    text: 'CPU利用率',
    value: 'cpuUtilization',
    disabled: false
  },
  {
    id: uuid(),
    text: '内存利用率',
    value: 'memoryUtilization',
    disabled: false
  },
  {
    id: uuid(),
    text: 'CPU使用量',
    value: 'cpuAverage',
    disabled: false
  },
  {
    id: uuid(),
    text: '内存使用量',
    value: 'memoryAverage',
    disabled: false
  }
];
const DisabledResourceMap = {
  cpu: ['cpu', 'cpuUtilization', 'cpuAverage'],
  cpuUtilization: ['cpu', 'cpuUtilization', 'cpuAverage'],
  cpuAverage: ['cpu', 'cpuUtilization', 'cpuAverage'],
  memory: ['memory', 'memoryUtilization', 'memoryAverage'],
  memoryUtilization: ['memory', 'memoryUtilization', 'memoryAverage'],
  memoryAverage: ['memory', 'memoryUtilization', 'memoryAverage']
};
const KindsMap = {
  deployments: 'Deployment',
  statefulsets: 'StatefulSet',
  tapps: 'TApp'
};
const ResourceTypeMap = {
  Deployment: 'deployments',
  StatefulSet: 'statefulsets',
  TApp: 'tapps'
};

/**
 * 组件实现
 */
const Hpa = React.memo((props: {
  selectedHpa?: any;
}) => {
  const { route, addons } = useSelector((state) => ({ route: state.route, addons: state.subRoot.addons }));
  const urlParams = router.resolve(route);
  const { mode } = urlParams;
  const { clusterId, projectName, HPAName } = route.queries;
  const { selectedHpa = {}} = props;

  const isModify = useMemo(() => mode === ModifyHpa, [mode]);

  /**
   * 表单初始化
   */
  const { register, watch, handleSubmit, reset, control, errors, setValue } = useForm<{
    name?: string;
    namespace?: string;
    resourceType: string;
    resource: string;
    strategy: any;
    minReplicas: number;
    maxReplicas: number;
  }>({
    mode: 'onBlur',
    defaultValues: {
      name: '',
      namespace: '',
      resourceType: '',
      resource: '',
      strategy: [{ key: '', value: 0 }],
      minReplicas: 0,
      maxReplicas: 0
    }
  });
  const { fields, append, remove } = useFieldArray({
    control,
    name: 'strategy'
  });
  const { namespace, resourceType, strategy, minReplicas } = watch();
  console.log('errors is:', strategy, errors);
  // modify的初始化
  const [selectedHpaNamespace, setSelectedHpaNamespace] = useState('');
  useEffect(() => {
    if (!isEmpty(selectedHpa)) {
      const { name, namespace } = selectedHpa.metadata;
      const { minReplicas, maxReplicas, scaleTargetRef, metrics } = selectedHpa.spec;
      setSelectedHpaNamespace(namespace);
      const selectedHPAStrategy =  metrics.map((item, index) => {
        const { name, targetAverageValue, targetAverageUtilization } = item.resource;
        let theKey = '';
        if (name === 'cpu' || name === 'memory') {
          const target = targetAverageValue ? 'targetAverageValue' : 'targetAverageUtilization';
          const { key } = NestedMetricsResourceMap[name][target];
          theKey = key;
        } else {
          const { key } = MetricsResourceMap[name];
          theKey = key;
        }
        return { key: theKey, value: parseInt(targetAverageValue || targetAverageUtilization) };
      });
      reset({
        resourceType: ResourceTypeMap[scaleTargetRef.kind],
        resource: scaleTargetRef.name,
        strategy: selectedHPAStrategy,
        minReplicas,
        maxReplicas
      });
    }
  }, [selectedHpa]);

  // 获取命名空间数据
  const namespaces = useNamespaces({ projectId: projectName, clusterId });

  /**
   * 设置【工作负载类型】数据
   */
  const resourceTypes = useMemo(() => {
    const initialValue = [{
      id: uuid(),
      text: 'Deployment',
      value: 'deployments'
    }, {
      id: uuid(),
      text: 'StatefulSet',
      value: 'statefulsets'
    }];
    let hasTapp = false;
    if (!isEmpty(addons) && !isEmpty(addons.TappController)) {
      hasTapp = true;
    }
    if (hasTapp) {
      initialValue.push({
        id: uuid(),
        text: 'Tapp',
        value: 'tapps'
      });
    }
    return initialValue;
  }, []);

  /**
   * 关联工作负载列表数据处理
   */
  const [resources, setResources] = useState({
    recordCount: 0,
    records: []
  });
  useEffect(() => {
    // 请求resourceType对应的列表数据
    async function getResourceList(resourceType, clusterId, namespace) {
      const resourceData = await fetchResourceList({ resourceType, clusterId, namespace });
      setResources(resourceData);
    }
    if (resourceType && clusterId && (namespace || selectedHpaNamespace)) {
      const selectNamespace = namespace || selectedHpaNamespace;
      getResourceList(resourceType, clusterId, selectNamespace);
    }
  }, [resourceType, clusterId, namespace, selectedHpaNamespace]);

  /**
   * 触发策略数据处理
   */
  // const [showStrategyOptions, setShowStrategyOptions] = useState(StrategyOptions);
  // useEffect(() => {
  //   if (strategy) {
  //     const selectedStrategyKeys = strategy.map(item => {
  //       return item.key;
  //     });
  //     let disabledStrategyKeys = [];
  //     selectedStrategyKeys.forEach(item => {
  //       if (item) {
  //         disabledStrategyKeys = [...disabledStrategyKeys, ...DisabledResourceMap[item]];
  //       }
  //     });
  //     const newStrategyOptions = StrategyOptions.map(item => {
  //       if (disabledStrategyKeys.indexOf(item.value) !== -1) {
  //         item.disabled = true;
  //       } else {
  //         item.disabled = false;
  //       }
  //       return item;
  //     });
  //     setShowStrategyOptions(newStrategyOptions);
  //   }
  // }, [strategy && strategy.length]);

  /**
   * 表单提交数据处理
   * @param data
   */
  const onSubmit = data => {
    const { name, namespace, minReplicas, maxReplicas, resource, resourceType, strategy } = data;
    const metrics = strategy.map(item => {
      // cpu || memory
      const name = item.key.replace('Utilization', '').replace('Average', '');
      // targetAverageValue || targetAverageUtilization
      const strategyKey = item.key.endsWith('Utilization') ? 'targetAverageUtilization' : 'targetAverageValue';
      return {
        type: 'Resource',
        resource: {
          name,
          [strategyKey]: +item.value
        }
      };
    });
    let hpaData;
    if (isModify) {
      selectedHpa.spec = {
        minReplicas: +minReplicas,
        maxReplicas: +maxReplicas,
        metrics,
        scaleTargetRef: {
          apiVersion: selectedHpa.spec.scaleTargetRef.apiVersion,
          kind: KindsMap[resourceType],
          name: resource
        }
      };
      hpaData = selectedHpa;
    } else {
      const newNamespace = cutNsStartClusterId({ namespace: data.namespace, clusterId });
      hpaData = {
        kind: 'HorizontalPodAutoscaler',
        apiVersion: 'autoscaling/v2beta1',
        metadata: {
          name,
          namespace: newNamespace,
          label: {
            'qcloud-app': name
          }
        },
        spec: {
          minReplicas: +minReplicas,
          maxReplicas: +maxReplicas,
          metrics,
          scaleTargetRef: {
            apiVersion: resourceType === 'tapps' ? 'apps.tkestack.io/v1' : 'apps/v1',
            kind: KindsMap[resourceType],
            name: resource
          }
        }
      };
    }

    async function addHpa() {
      const addHpaResult = await createHPA({
        namespace: data.namespace,
        clusterId,
        hpaData
      });
      if (addHpaResult) {
        // 跳转到列表页面
        router.navigate({ ...urlParams, mode: 'list' }, route.queries);
      }
    }
    async function updateHpa() {
      const addHpaResult = await modifyHPA({
        name: HPAName,
        namespace: selectedHpa.metadata.namespace,
        clusterId,
        hpaData
      });
      if (addHpaResult) {
        // 跳转到列表页面
        router.navigate({ ...urlParams, mode: 'list' }, route.queries);
      }
    }

    if (isModify) {
      updateHpa();
    } else {
      addHpa();
    }
  };
  return (
    <Layout>
      <Body>
        <Content>
          <Content.Header
            showBackButton
            onBackButtonClick={() => history.back()}
            title={isModify ? t('修改HPA') : t('新建HPA')}
          />
          <Content.Body>
            <Card>
              <Card.Body>
                <form onSubmit={handleSubmit(onSubmit)}>
                  <Form>
                    <Form.Item
                      required
                      label={t('名称')}
                      showStatusIcon={false}
                      status={errors.name ? 'error' : 'success'}
                      message={
                        errors.name ? errors.name.message : (isModify ? '' : t(
                          '最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾'
                          ))
                      }
                    >
                      {
                        isModify ?
                          <Text parent="div" align="left" reset>{selectedHpa.metadata ? selectedHpa.metadata.name : ''}</Text>
                          :
                          (
                            <Controller
                              as={Input}
                              name="name"
                              control={control}
                              rules={{
                                required: t('名称不能为空'),
                                maxLength: {
                                  value: 63,
                                  message: t('名称不能超过63个字符')
                                },
                                pattern: {
                                  value: /^[a-z]([-a-z0-9]*[a-z0-9])?$/,
                                  message: t('名称格式不正确')
                                }
                              }}
                            />
                          )
                      }
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('命名空间')}
                      showStatusIcon={false}
                      status={errors.namespace ? 'error' : 'success'}
                      message={errors.namespace && errors.namespace.message}
                    >
                      {
                        isModify ?
                          <Text parent="div" align="left" reset>{selectedHpa.metadata ? selectedHpa.metadata.namespace : ''}</Text>
                          :
                          (
                            <Controller
                              as={
                                <Select
                                  searchable
                                  boxSizeSync
                                  type="simulate"
                                  appearence="button"
                                  size="m"
                                  options={isEmpty(namespaces) ? [] : namespaces.records}
                                />
                              }
                              name="namespace"
                              control={control}
                              rules={{ required: t('命名空间不能为空') }}
                            />
                         )
                      }
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('工作负载类型')}
                      showStatusIcon={false}
                      status={errors.resourceType ? 'error' : 'success'}
                      message={errors.resourceType && errors.resourceType.message}
                    >
                      <Controller
                        as={
                          <Select
                            searchable
                            boxSizeSync
                            type="simulate"
                            appearence="button"
                            size="m"
                            options={resourceTypes}
                          />
                        }
                        name="resourceType"
                        control={control}
                        rules={{ required: t('工作负载类型不能为空') }}
                      />
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('关联工作负载')}
                      showStatusIcon={false}
                      status={errors.resource ? 'error' : 'success'}
                      message={errors.resource && errors.resource.message}
                    >
                      <Controller
                        as={
                          <Select
                            searchable
                            boxSizeSync
                            type="simulate"
                            appearence="button"
                            size="m"
                            options={isEmpty(resources) ? [] : resources.records}
                          />
                        }
                        name="resource"
                        control={control}
                        rules={{ required: t('关联工作负载不能为空') }}
                      />
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('触发策略')}
                      showStatusIcon={false}
                    >
                      <ul className="hpa-edit-strategy-ul">
                        {
                          fields.map((item, index) => {
                            return (
                              <li key={item.id} className="hpa-edit-strategy-li">
                                <Controller
                                  as={
                                    <Select
                                      boxSizeSync
                                      type="simulate"
                                      appearence="button"
                                      size="m"
                                      options={StrategyOptions}
                                      className={errors.strategy && errors.strategy[index] && errors.strategy[index].key ? 'is-error' : ''}
                                    />
                                  }
                                  name={`strategy[${index}].key`}
                                  control={control}
                                  defaultValue={item.key}
                                  rules={{
                                    required: t('请选择触发策略'),
                                    validate: value => {
                                      // 之前的选择在这个值对应的的[]中出现了，这里就要报错重复或者二者不能同时选择
                                      const disabledSelectValue = DisabledResourceMap[value];
                                      let tip = '',
                                          repeatFlag = false,
                                          similarFlag = false;
                                      strategy && strategy.forEach((item, i) => {
                                        // 非当前值中判断
                                        if (index !== i) {
                                          // 非当前值中有相同策略选择
                                          if (item.key === value) {
                                            repeatFlag = true;
                                          }
                                          // 非当前值中有相似策略选择
                                          if (disabledSelectValue.indexOf(item.key) !== -1) {
                                            similarFlag = true;
                                          }
                                        }
                                      });
                                      if (repeatFlag) {
                                        tip = '相同指标不能重复设置';
                                      } else if (similarFlag) {
                                        const element = value.replace('Utilization', '').replace('Average', '');
                                        tip = (element === 'cpu' ? 'CPU' : '内存') + '的利用率和使用量不能同时设置';
                                      }
                                      console.log('触发策略select value： ', value, tip, strategy, repeatFlag, similarFlag);
                                      if (tip) {
                                        return tip;
                                      }
                                    }
                                  }}
                                />
                                <Text style={{ fontSize: '14px' }}> </Text>
                                {
                                  strategy && strategy[index] && strategy[index].key && MetricsResourceMap[strategy[index].key].unit === '%'
                                    ?
                                      <Controller
                                        as={
                                          <InputNumber
                                            step={1}
                                            min={0}
                                            max={100}
                                            className={errors.strategy && errors.strategy[index] && errors.strategy[index].value ? 'is-error' : ''}
                                            unit="%"
                                          />}
                                        name={`strategy[${index}].value`}
                                        size="s"
                                        control={control}
                                        defaultValue={item.value}
                                        rules={{
                                          required: t('请输入0-100之间的整数')
                                        }}
                                      />
                                    :
                                      <Controller
                                        as={
                                          <InputNumber
                                            step={1}
                                            min={0}
                                            className={errors.strategy && errors.strategy[index] && errors.strategy[index].value ? 'is-error' : ''}
                                            unit={
                                              strategy && strategy[index] && strategy[index].key
                                                  ?
                                                MetricsResourceMap[strategy[index].key].unit
                                                  :
                                                ''
                                            }
                                          />}
                                        name={`strategy[${index}].value`}
                                        size="s"
                                        control={control}
                                        defaultValue={item.value}
                                        rules={{
                                          required: t('请输入大于等于0的数')
                                        }}
                                      />
                                }
                                {
                                  strategy && strategy.length > 1 &&
                                  <LinkButton onClick={(e) => {
                                    e.preventDefault();
                                    remove(index);
                                  }}>
                                    <i className="icon-cancel-icon" />
                                  </LinkButton>
                                }
                                {
                                  errors.strategy && errors.strategy[index] &&
                                    <Text parent="div" theme="danger" reset>
                                      {errors.strategy[index].key ? errors.strategy[index].key.message : ''}
                                    </Text>
                                }
                              </li>
                            );
                          })
                        }
                      </ul>
                      {/*{errors.strategy && <Text parent="div" theme="danger" reset><Trans>不能有空内容</Trans></Text>}*/}
                      <LinkButton onClick={(e) => {
                        e.preventDefault();
                        append({ key: '', value: 0 });
                      }}>
                        <Trans>新增指标</Trans>
                      </LinkButton>
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('实例范围')}
                      showStatusIcon={false}
                      status={errors.minReplicas || errors.maxReplicas ? 'error' : 'success'}
                      message={(errors.minReplicas && errors.minReplicas.message) || (errors.maxReplicas && errors.maxReplicas.message)}
                    >
                      <Controller
                        as={<InputNumber step={1} min={0} />}
                        name="minReplicas"
                        size="s"
                        control={control}
                        rules={{
                          required: t('最小值不能为空'),
                          min: {
                            value: 0,
                            message: t('最小值需是大于等于0的整数')
                          }
                        }}
                      />
                      <Text style={{ fontSize: '14px', verticalAlign: 'middle' }}> ~ </Text>
                      <Controller
                        as={<InputNumber step={1} min={0} />}
                        name="maxReplicas"
                        size="s"
                        control={control}
                        rules={{
                          required: t('最大值不能为空'),
                          validate: value => {
                            return value > minReplicas || t('最大值需是大于最小值的整数');
                          }
                        }}
                      />
                    </Form.Item>
                  </Form>
                  <Form.Action>
                    <Button htmlType="submit" type="primary"><Trans>保存</Trans></Button>
                    <Button htmlType="button" onClick={() => { history.back() }}><Trans>取消</Trans></Button>
                  </Form.Action>
                </form>
              </Card.Body>
            </Card>
          </Content.Body>
        </Content>
      </Body>
    </Layout>
  );
});
export default Hpa;
