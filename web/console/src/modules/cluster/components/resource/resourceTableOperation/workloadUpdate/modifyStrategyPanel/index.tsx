import { t } from '@/tencent/tea-app/lib/i18n';
import { TeaFormLayout } from '@src/modules/common/layouts/TeaFormLayout';
import React, { useEffect } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { Button, Form, InputNumber, Radio, Select } from 'tea-component';
import { IModifyPanelProps, WorkloadKindEnum } from '../constants';
import {
  RegistryUpdateTypeEnum,
  RollingUpdateTypeEnum,
  getUpdateTypeOptionsForKind,
  updateStrategyOptions
} from './constants';

export const ModifyStrategyPanel = ({ kind, resource, title, baseInfo, onCancel, onUpdate }: IModifyPanelProps) => {
  const isDeployment = kind === WorkloadKindEnum.Deployment;
  const isStatefulSet = kind === WorkloadKindEnum.StatefulSet;
  const isDaemonSet = kind === WorkloadKindEnum.DaemonSet;

  const { handleSubmit, watch, control, reset } = useForm({
    mode: 'onBlur',

    defaultValues: {
      updateType: RegistryUpdateTypeEnum.RollingUpdate,
      updateInterval: 0,
      updateStrategy: RollingUpdateTypeEnum.UserDefined,
      maxSurge: 25,
      maxUnavailable: 25,
      batchSize: 1,
      partition: 0
    }
  });

  const updateTypeWatch = watch('updateType');
  const updateStrategyWatch = watch('updateStrategy');

  useEffect(() => {
    if (!resource) return;

    const updateType = isDeployment ? resource?.spec?.strategy?.type : resource?.spec?.updateStrategy?.type;

    // 如果只有maxSurge 有值，而maxUnavailable为0，则为启动新的pod，停止旧的pod
    const minReadySeconds = resource?.spec?.minReadySeconds ?? 0,
      partition = resource?.spec?.updateStrategy?.rollingUpdate?.partition ?? 0;

    let maxSurge = resource?.spec?.strategy?.rollingUpdate?.maxSurge ?? 0,
      maxUnavailable = resource?.spec?.strategy?.rollingUpdate?.maxUnavailable ?? 0,
      rollingUpdateStrategy = RollingUpdateTypeEnum.CreatePod,
      batchSize = 1;

    if (maxSurge === 0 && Number.isInteger(maxUnavailable)) {
      rollingUpdateStrategy = RollingUpdateTypeEnum.DestroyPod;
      batchSize = maxUnavailable;
      maxSurge = '25%';
      maxUnavailable = '25%';
    } else if (maxUnavailable === 0 && Number.isInteger(maxSurge)) {
      rollingUpdateStrategy = RollingUpdateTypeEnum.CreatePod;
      batchSize = maxSurge;
      maxSurge = '25%';
      maxUnavailable = '25%';
    } else {
      rollingUpdateStrategy = RollingUpdateTypeEnum.UserDefined;
    }

    if (isDaemonSet) {
      maxUnavailable = resource?.spec?.updateStrategy?.rollingUpdate?.maxUnavailable ?? 0;
    }

    const newConfig = {
      updateType,
      updateInterval: minReadySeconds,
      updateStrategy: rollingUpdateStrategy,
      maxSurge: parseInt(maxSurge),
      maxUnavailable: parseInt(maxUnavailable),
      batchSize,
      partition
    };

    reset(newConfig);
  }, [resource, isDaemonSet, isDeployment, isStatefulSet, reset]);

  function onSubmit({ updateType, updateInterval, updateStrategy, maxSurge, maxUnavailable, batchSize, partition }) {
    const isRollingUpdate = updateType === RegistryUpdateTypeEnum.RollingUpdate;

    const isUserDefined = updateStrategy === RollingUpdateTypeEnum.UserDefined;
    const isCreatePod = updateStrategy === RollingUpdateTypeEnum.CreatePod;

    const data = {
      spec: {
        minReadySeconds: isStatefulSet ? undefined : isRollingUpdate ? updateInterval : 0,

        strategy: isDeployment
          ? {
              type: updateType,
              rollingUpdate: isRollingUpdate
                ? {
                    maxSurge: isUserDefined ? `${maxSurge}%` : isCreatePod ? batchSize : 0,

                    maxUnavailable: isUserDefined ? `${maxUnavailable}%` : isCreatePod ? 0 : batchSize
                  }
                : null
            }
          : undefined,

        updateStrategy: isDeployment
          ? undefined
          : {
              type: updateType,
              rollingUpdate: isRollingUpdate
                ? isStatefulSet
                  ? { partition }
                  : maxUnavailable
                  ? { maxUnavailable }
                  : undefined
                : null
            }
      }
    };

    onUpdate(data);
  }

  return (
    <TeaFormLayout
      title={title}
      footer={
        <>
          <Button type="primary" style={{ marginRight: 10 }} onClick={handleSubmit(onSubmit)}>
            设置更新策略
          </Button>

          <Button onClick={onCancel}>取消</Button>
        </>
      }
    >
      {baseInfo}

      <hr />

      <Form>
        <Controller
          name="updateType"
          control={control}
          render={({ field }) => (
            <Form.Item
              label={t('更新方式')}
              extra={
                updateTypeWatch === RegistryUpdateTypeEnum.RollingUpdate
                  ? '对实例进行逐个更新，这种方式可以让您不中断业务实现对服务的更新'
                  : '直接关闭所有实例，启动相同数量的新实例'
              }
            >
              <Select appearance="button" options={getUpdateTypeOptionsForKind(kind)} {...field} />
            </Form.Item>
          )}
        />

        {(isDeployment || isDaemonSet) && updateTypeWatch === RegistryUpdateTypeEnum.RollingUpdate && (
          <Controller
            name="updateInterval"
            control={control}
            render={({ field }) => (
              <Form.Item label={t('更新间隔')} extra={t('')}>
                <InputNumber hideButton unit="秒" size="l" min={0} step={1} {...field} />
              </Form.Item>
            )}
          />
        )}

        {isDeployment && updateTypeWatch === RegistryUpdateTypeEnum.RollingUpdate && (
          <Controller
            name="updateStrategy"
            control={control}
            render={({ field }) => (
              <Form.Item
                label={t('更新策略')}
                extra={
                  updateStrategyWatch === RollingUpdateTypeEnum.CreatePod
                    ? '请确认集群有足够的CPU和内存用于启动新的Pod, 否则可能导致集群崩溃'
                    : ''
                }
              >
                <Radio.Group value={field.value} onChange={field.onChange}>
                  {updateStrategyOptions.map(({ text, value }) => (
                    <Radio key={value} name={value}>
                      {text}
                    </Radio>
                  ))}
                </Radio.Group>
              </Form.Item>
            )}
          />
        )}

        {updateTypeWatch === RegistryUpdateTypeEnum.RollingUpdate && (
          <Form.Item label={t('策略配置')}>
            <Form>
              {isDeployment && updateStrategyWatch === RollingUpdateTypeEnum.UserDefined && (
                <Controller
                  name="maxSurge"
                  control={control}
                  render={({ field }) => (
                    <Form.Item label={t('MaxSurge')} extra={t('允许超出所需规模的最大Pod数量')}>
                      <InputNumber hideButton unit="%" size="l" min={0} step={1} {...field} />
                    </Form.Item>
                  )}
                />
              )}

              {((isDeployment && updateStrategyWatch === RollingUpdateTypeEnum.UserDefined) || isDaemonSet) && (
                <Controller
                  name="maxUnavailable"
                  control={control}
                  render={({ field }) => (
                    <Form.Item label={t('MaxUnavailable')} extra={t('允许最大不可用的Pod数量')}>
                      <InputNumber hideButton unit="%" size="l" min={0} step={1} {...field} />
                    </Form.Item>
                  )}
                />
              )}

              {isDeployment && updateStrategyWatch !== RollingUpdateTypeEnum.UserDefined && (
                <Controller
                  name="batchSize"
                  control={control}
                  render={({ field }) => (
                    <Form.Item label={t('Pods')} extra={t('Pod将批量启动或停止')}>
                      <InputNumber hideButton size="l" min={0} step={1} {...field} />
                    </Form.Item>
                  )}
                />
              )}

              {isStatefulSet && (
                <Controller
                  name="partition"
                  control={control}
                  render={({ field }) => (
                    <Form.Item label="Partition">
                      <InputNumber hideButton size="l" min={0} step={1} {...field} />
                    </Form.Item>
                  )}
                />
              )}
            </Form>
          </Form.Item>
        )}
      </Form>
    </TeaFormLayout>
  );
};
