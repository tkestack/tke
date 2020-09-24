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
  Input
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
import { MetricsResourceMap } from '../constant';
import { useNamespaces } from '@src/modules/cluster/components/scale/common/hooks';

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
  tapps: 'Tapp'
};
const ResourceTypeMap = {
  Deployment: 'deployments',
  StatefulSet: 'statefulsets',
  Tapp: 'tapps'
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
    strategy: NestedValue<any[]>;
    minReplicas: number;
    maxReplicas: number;
  }>({
    mode: 'onBlur',
    defaultValues: {
      name: '',
      namespace: '',
      resourceType: '',
      resource: '',
      strategy: [{ key: '', value: '' }],
      minReplicas: undefined,
      maxReplicas: undefined
    }
  });
  const { fields, append, remove } = useFieldArray({
    control,
    name: 'strategy'
  });
  const { namespace, resourceType, strategy, minReplicas } = watch();

  // modify的初始化
  const [selectedHpaNamespace, setSelectedHpaNamespace] = useState('');
  useEffect(() => {
    if (!isEmpty(selectedHpa)) {
      const { name, namespace } = selectedHpa.metadata;
      const { minReplicas, maxReplicas, scaleTargetRef, metrics } = selectedHpa.spec;
      setSelectedHpaNamespace(namespace);
      const selectedHPAStrategy =  metrics.map((item, index) => {
        const { name, targetAverageValue, targetAverageUtilization } = item.resource;
        return { key: MetricsResourceMap[name].key, value: parseFloat(targetAverageValue || targetAverageUtilization) };
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
  const [resources, setResources] = useState();
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
  const [showStrategyOptions, setShowStrategyOptions] = useState(StrategyOptions);
  useEffect(() => {
    if (strategy) {
      const selectedStrategyKeys = strategy.map(item => {
        return item.key;
      });
      let disabledStrategyKeys = [];
      selectedStrategyKeys.forEach(item => {
        if (item) {
          disabledStrategyKeys = [...disabledStrategyKeys, ...DisabledResourceMap[item]];
        }
      });
      const newStrategyOptions = StrategyOptions.map(item => {
        if (disabledStrategyKeys.indexOf(item.value) !== -1) {
          item.disabled = true;
        } else {
          item.disabled = false;
        }
        return item;
      });
      setShowStrategyOptions(newStrategyOptions);
    }
  }, [strategy && strategy.length]);

  /**
   * 表单提交数据处理
   * @param data
   */
  const onSubmit = data => {
    const { name, namespace, minReplicas, maxReplicas, resource, resourceType, strategy } = data;
    const metrics = strategy.map(item => {
      return {
        type: 'Resource',
        resource: {
          name: item.key.replace('Utilization', ''),
          targetAverageUtilization: +item.value
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
            apiVersion: 'apps/v1',
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
                              <li key={item.id} className={errors.strategy && errors.strategy[index] ? 'hpa-edit-strategy-li is-error' : 'hpa-edit-strategy-li'}>
                                <Controller
                                  as={
                                    <Select
                                      boxSizeSync
                                      type="simulate"
                                      appearence="button"
                                      size="m"
                                      options={showStrategyOptions}
                                    />
                                  }
                                  name={`strategy[${index}].key`}
                                  control={control}
                                  defaultValue={item.key}
                                  rules={{
                                    required: t('不能有空内容'),
                                  }}
                                />
                                <Text style={{ fontSize: '14px' }}> </Text>
                                <Controller
                                  as={Input}
                                  name={`strategy[${index}].value`}
                                  size="s"
                                  control={control}
                                  defaultValue={item.value}
                                  rules={{
                                    required: t('不能有空内容')
                                  }}
                                />
                                <Text style={{ fontSize: '14px', verticalAlign: 'middle' }}> {strategy && strategy[index] && strategy[index].key ? MetricsResourceMap[strategy[index].key].unit : ''}</Text>
                                <LinkButton onClick={(e) => {
                                  e.preventDefault();
                                  remove(index);
                                }}>
                                  <i className="icon-cancel-icon" />
                                </LinkButton>
                              </li>
                            );
                          })
                        }
                      </ul>
                      {errors.strategy && <Text parent="div" theme="danger" reset><Trans>不能有空内容</Trans></Text>}
                      <LinkButton onClick={(e) => {
                        e.preventDefault();
                        append({ key: '', value: '' });
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
                        as={<Input type="number" />}
                        name="minReplicas"
                        size="s"
                        control={control}
                        rules={{
                          required: t('最小值不能为空'),
                          min: {
                            value: 1,
                            message: t('最小值大于等于1')
                          }
                        }}
                      />
                      <Text style={{ fontSize: '14px', verticalAlign: 'middle' }}> ~ </Text>
                      <Controller
                        as={<Input type="number" />}
                        name="maxReplicas"
                        size="s"
                        control={control}
                        rules={{
                          required: t('最大值不能为空'),
                          validate: value => {
                            return value > minReplicas || t('最大值应大于最小值的整数');
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
