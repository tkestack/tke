/**
 * HPA创建及修改组件
 */
import React, { useState, useEffect, useContext, useMemo, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Layout, Card, Select, Text, Button, Form, Input, InputNumber } from '@tencent/tea-component';
import { LinkButton } from '@src/modules/common/components';
import { Resource } from '@src/modules/common/models';
import { useForm, useFieldArray, Controller, NestedValue } from 'react-hook-form';
import { isEmpty } from '@src/modules/common/utils';
import { insertCSS, uuid } from '@tencent/ff-redux/libs/qcloud-lib';
import { router } from '@src/modules/cluster/router';
import { cutNsStartClusterId } from '@helper';
import { createCronHpa, modifyCronHpa, fetchResourceList } from '@src/modules/cluster/WebAPI/scale';
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
const ModifyHpa = 'modify-cronhpa';

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
const Hpa = React.memo((props: { selectedHpa?: any }) => {
  const { route, addons } = useSelector(state => ({ route: state.route, addons: state.subRoot.addons }));
  const urlParams = router.resolve(route);
  const { mode } = urlParams;
  const { clusterId, projectName, HPAName } = route.queries;
  const { selectedHpa = {} } = props;

  const isModify = useMemo(() => mode === ModifyHpa, [mode]);

  /**
   * 表单初始化
   */
  const {
    register,
    watch,
    handleSubmit,
    reset,
    control,
    setValue,
    formState: { errors }
  } = useForm<{
    name?: string;
    namespace?: string;
    resourceType: string;
    resource: string;
    strategy: { key: string; value: number }[];
    minReplicas: number;
    maxReplicas: number;
  }>({
    mode: 'onBlur',
    defaultValues: {
      name: '',
      namespace: '',
      resourceType: '',
      resource: '',
      strategy: [{ key: '', value: 0 }]
    }
  });
  const { fields, append, remove } = useFieldArray({
    control,
    name: 'strategy'
  });
  const { namespace, resourceType, strategy } = watch();

  // modify的初始化
  const [selectedHpaNamespace, setSelectedHpaNamespace] = useState('');
  useEffect(() => {
    if (!isEmpty(selectedHpa)) {
      const { name, namespace } = selectedHpa.metadata;
      const { scaleTargetRef, crons } = selectedHpa.spec;
      setSelectedHpaNamespace(namespace);
      const selectedHPAStrategy = crons.map((item, index) => {
        const { schedule, targetReplicas } = item;
        return { key: schedule, value: parseInt(targetReplicas) };
      });
      reset({
        resourceType: ResourceTypeMap[scaleTargetRef.kind],
        resource: scaleTargetRef.name,
        strategy: selectedHPAStrategy
      });
    }
  }, [selectedHpa]);

  // 获取命名空间数据
  const namespaces = useNamespaces({ projectId: projectName, clusterId });

  /**
   * 设置【工作负载类型】数据
   */
  const resourceTypes = useMemo(() => {
    const initialValue = [
      {
        id: uuid(),
        text: 'Deployment',
        value: 'deployments'
      },
      {
        id: uuid(),
        text: 'StatefulSet',
        value: 'statefulsets'
      }
    ];
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
   * 表单提交数据处理
   * @param data
   */
  const onSubmit = data => {
    const { name, namespace, resource, resourceType, strategy } = data;
    const crons = strategy.map(item => {
      const { key, value } = item;
      return {
        schedule: key,
        targetReplicas: +value
      };
    });
    let cronHpaData;
    if (isModify) {
      selectedHpa.spec = {
        crons,
        scaleTargetRef: {
          apiVersion: selectedHpa.spec.scaleTargetRef.apiVersion,
          kind: KindsMap[resourceType],
          name: resource
        }
      };
      cronHpaData = selectedHpa;
    } else {
      const newNamespace = cutNsStartClusterId({ namespace: data.namespace, clusterId });
      cronHpaData = {
        kind: 'CronHPA',
        apiVersion: 'extensions.tkestack.io/v1',
        metadata: {
          name,
          namespace: newNamespace
        },
        spec: {
          crons,
          scaleTargetRef: {
            apiVersion: resourceType === 'tapps' ? 'apps.tkestack.io/v1' : 'apps/v1',
            kind: KindsMap[resourceType],
            name: resource
          }
        }
      };
    }

    async function addCronHpa() {
      const addHpaResult = await createCronHpa({
        namespace: data.namespace,
        clusterId,
        cronHpaData
      });
      if (addHpaResult) {
        // 跳转到列表页面
        router.navigate({ ...urlParams, mode: 'list' }, route.queries);
      }
    }
    async function updateCronHpa() {
      const addHpaResult = await modifyCronHpa({
        name: HPAName,
        namespace: selectedHpa.metadata.namespace,
        clusterId,
        cronHpaData
      });
      if (addHpaResult) {
        // 跳转到列表页面
        router.navigate({ ...urlParams, mode: 'list' }, route.queries);
      }
    }

    if (isModify) {
      updateCronHpa();
    } else {
      addCronHpa();
    }
  };
  return (
    <Layout>
      <Body>
        <Content>
          <Content.Header
            showBackButton
            onBackButtonClick={() => history.back()}
            title={isModify ? t('修改CronHPA') : t('新建CronHPA')}
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
                        errors.name
                          ? errors.name.message
                          : isModify
                          ? ''
                          : t(
                              '最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾'
                            )
                      }
                    >
                      {isModify ? (
                        <Text parent="div" align="left" reset>
                          {selectedHpa.metadata ? selectedHpa.metadata.name : ''}
                        </Text>
                      ) : (
                        <Controller
                          render={({ field }) => <Input {...field} />}
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
                      )}
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('命名空间')}
                      showStatusIcon={false}
                      status={errors.namespace ? 'error' : 'success'}
                      message={errors.namespace && errors.namespace.message}
                    >
                      {isModify ? (
                        <Text parent="div" align="left" reset>
                          {selectedHpa.metadata ? selectedHpa.metadata.namespace : ''}
                        </Text>
                      ) : (
                        <Controller
                          render={({ field }) => (
                            <Select
                              {...field}
                              searchable
                              boxSizeSync
                              type="simulate"
                              appearence="button"
                              size="m"
                              options={isEmpty(namespaces) ? [] : namespaces.records}
                            />
                          )}
                          name="namespace"
                          control={control}
                          rules={{ required: t('命名空间不能为空') }}
                        />
                      )}
                    </Form.Item>
                    <Form.Item
                      required
                      label={t('工作负载类型')}
                      showStatusIcon={false}
                      status={errors.resourceType ? 'error' : 'success'}
                      message={errors.resourceType && errors.resourceType.message}
                    >
                      <Controller
                        render={({ field }) => (
                          <Select
                            {...field}
                            searchable
                            boxSizeSync
                            type="simulate"
                            appearence="button"
                            size="m"
                            options={resourceTypes}
                          />
                        )}
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
                        render={({ field }) => (
                          <Select
                            {...field}
                            searchable
                            boxSizeSync
                            type="simulate"
                            appearence="button"
                            size="m"
                            options={isEmpty(resources) ? [] : resources.records}
                          />
                        )}
                        name="resource"
                        control={control}
                        rules={{ required: t('关联工作负载不能为空') }}
                      />
                    </Form.Item>
                    <Form.Item required label={t('触发策略')} showStatusIcon={false}>
                      <Text theme="label" parent="p" reset style={{ marginBottom: '5px' }}>
                        {t(
                          '根据设置的Crontab（Crontab语法格式，例如 "0 23 * * 5"表示每周五23:00）周期性地设置实例数量'
                        )}
                      </Text>
                      <ul className="hpa-edit-strategy-ul">
                        {fields.map((item, index) => {
                          return (
                            <li key={item.id} className="hpa-edit-strategy-li">
                              <Controller
                                render={({ field }) => (
                                  <Input
                                    {...field}
                                    placeholder="Crontab"
                                    className={
                                      errors.strategy && errors.strategy[index] && errors.strategy[index].key
                                        ? 'is-error'
                                        : ''
                                    }
                                  />
                                )}
                                name={`strategy.${index}.key`}
                                control={control}
                                defaultValue={item.key}
                                rules={{
                                  required: t('执行策略不能为空'),
                                  pattern: {
                                    value:
                                      /^(\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\*\/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\*|([0-9]|1[0-9]|2[0-3])|\*\/([0-9]|1[0-9]|2[0-3])) (\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\*\/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\*|([1-9]|1[0-2])|\*\/([1-9]|1[0-2])) (\*|([0-6])|\*\/([0-6]))$/,
                                    message: t('执行策略格式不正确')
                                  }
                                }}
                              />
                              <Text style={{ fontSize: '14px' }}> </Text>
                              <Controller
                                render={({ field }) => (
                                  <InputNumber
                                    {...field}
                                    step={1}
                                    min={0}
                                    unit={t('个实例')}
                                    className={
                                      errors.strategy && errors.strategy[index] && errors.strategy[index].value
                                        ? 'is-error'
                                        : ''
                                    }
                                  />
                                )}
                                name={`strategy.${index}.value`}
                                control={control}
                                defaultValue={item.value}
                                rules={{
                                  required: t('目标实例数不能为空'),
                                  min: {
                                    value: 0,
                                    message: t('目标实例数需大于等于0')
                                  }
                                }}
                              />
                              {strategy && strategy.length > 1 && (
                                <LinkButton
                                  onClick={e => {
                                    e.preventDefault();
                                    remove(index);
                                  }}
                                >
                                  <i className="icon-cancel-icon" />
                                </LinkButton>
                              )}
                              {errors.strategy && errors.strategy[index] && (
                                <Text parent="div" theme="danger" reset>
                                  {errors.strategy[index].key ? errors.strategy[index].key.message : ''}
                                </Text>
                              )}
                            </li>
                          );
                        })}
                      </ul>
                      {/*{errors.strategy && <Text parent="div" theme="danger" reset><Trans>不能有空内容或者格式不正确</Trans></Text>}*/}
                      <LinkButton
                        onClick={e => {
                          e.preventDefault();
                          append({ key: '', value: 0 });
                        }}
                      >
                        <Trans>新增策略</Trans>
                      </LinkButton>
                    </Form.Item>
                  </Form>
                  <Form.Action>
                    <Button htmlType="submit" type="primary">
                      <Trans>保存</Trans>
                    </Button>
                    <Button
                      htmlType="button"
                      onClick={() => {
                        history.back();
                      }}
                    >
                      <Trans>取消</Trans>
                    </Button>
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
