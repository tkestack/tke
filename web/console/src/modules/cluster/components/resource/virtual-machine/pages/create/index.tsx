import React, { useEffect } from 'react';
import { TeaFormLayout } from '@src/modules/common/layouts/TeaFormLayout';
import { Form, Input, TextArea, Text, Select, Button, InputNumber, InputAdornment, Segment } from 'tea-component';
import { useRecoilValueLoadable, useRecoilState, useRecoilValue } from 'recoil';
import {
  namespaceListState,
  clusterIdState,
  mirrorListState,
  diskListState,
  diskListValidateState
} from '../../store/creation';
import { DiskPanel } from './diskPanel';
import { virtualMachineAPI } from '@src/webApi';
import { Controller, useForm } from 'react-hook-form';
import { getReactHookFormStatusWithMessage } from '@helper';

export const VMCreatePanel = () => {
  const [clusterId, setClusterId] = useRecoilState(clusterIdState);
  const namespaceListLoadable = useRecoilValueLoadable(namespaceListState);
  const mirrorListLoadable = useRecoilValueLoadable(mirrorListState);
  const diskList = useRecoilValue(diskListState);
  const diskListValidate = useRecoilValue(diskListValidateState);

  const { control, handleSubmit } = useForm({
    mode: 'onBlur',
    defaultValues: {
      name: '',
      description: '',
      namespace: 'default',
      cpu: 1,
      memory: 1,
      mirror: null,
      networkMode: 'Bridge'
    }
  });

  useEffect(() => {
    const url = new URL(location.href);
    const clusterId = url.searchParams.get('clusterId');

    if (clusterId) setClusterId(clusterId);
  }, [location.href]);

  async function create({ name, description, namespace, cpu, memory, mirror, networkMode }) {
    const diskListValidateFailed = diskListValidate.some(item =>
      Object.values(item).some(({ status }) => status === 'error')
    );

    if (diskListValidateFailed) return;

    try {
      await virtualMachineAPI.createVM({
        clusterId,
        namespace,
        vmOptions: {
          name,
          description,
          networkMode,
          mirror: mirrorListLoadable?.contents?.find(({ value }) => value === mirror),
          diskList,
          cpu,
          memory
        }
      });

      history.back();
    } catch (error) {
      console.log('createVm error --->', error);
    }
  }

  return (
    <TeaFormLayout
      title="创建虚拟机"
      footer={
        <>
          <Button type="primary" style={{ marginRight: 20 }} onClick={handleSubmit(create)}>
            创建虚拟机
          </Button>

          <Button onClick={() => history.back()}>取消</Button>
        </>
      }
    >
      <Form>
        <Controller
          control={control}
          name="name"
          rules={{
            required: '虚拟机名称必填！',
            maxLength: 63,
            pattern: /^[a-z]([-a-z0-9]*[a-z0-9])?$/
          }}
          render={({ field, ...others }) => (
            <Form.Item
              label="虚拟机名称"
              extra={`最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾`}
              {...getReactHookFormStatusWithMessage(others)}
            >
              <Input placeholder="请输入虚拟机名称" {...field} />
            </Form.Item>
          )}
        />

        <Controller
          control={control}
          name="description"
          rules={{
            maxLength: 1000
          }}
          render={({ field, ...others }) => (
            <Form.Item label="描述" {...getReactHookFormStatusWithMessage(others)}>
              <TextArea placeholder="请输入描述信息，不能超过1000个字符" {...field} />
            </Form.Item>
          )}
        />

        <Form.Item label="集群">
          <Text reset>{clusterId}</Text>
        </Form.Item>

        <Controller
          control={control}
          name="namespace"
          rules={{ required: '命名空间必选!' }}
          render={({ field, ...others }) => (
            <Form.Item label="命名空间" {...getReactHookFormStatusWithMessage(others)}>
              <Select
                type="simulate"
                searchable
                appearence="button"
                size="s"
                style={{ width: '130px', marginRight: '5px' }}
                options={
                  namespaceListLoadable?.state === 'hasValue'
                    ? namespaceListLoadable?.contents?.map(value => ({ value }))
                    : []
                }
                {...field}
              />
            </Form.Item>
          )}
        />

        <Form.Item label="规格">
          <Controller
            control={control}
            name="cpu"
            render={({ field, ...others }) => (
              <InputAdornment before="CPU" after="核" style={{ marginRight: 30 }}>
                <InputNumber hideButton min={1} step={1} {...field} />
              </InputAdornment>
            )}
          />

          <Controller
            control={control}
            name="memory"
            render={({ field, ...others }) => (
              <InputAdornment before="内存" after="Gi">
                <InputNumber hideButton min={1} step={1} {...field} />
              </InputAdornment>
            )}
          />
        </Form.Item>

        <Controller
          control={control}
          name="mirror"
          rules={{ required: '镜像必选!' }}
          render={({ field, ...others }) => (
            <Form.Item label="镜像" {...getReactHookFormStatusWithMessage(others)}>
              <Select
                type="simulate"
                matchButtonWidth
                style={{ width: 200 }}
                searchable
                appearence="button"
                options={mirrorListLoadable?.state === 'hasValue' ? mirrorListLoadable?.contents : []}
                {...field}
                onChange={value => {
                  field.onChange(value);
                  field.onBlur();
                }}
              />
            </Form.Item>
          )}
        />

        <Form.Item label="磁盘">
          <DiskPanel />
        </Form.Item>

        <Controller
          control={control}
          name="networkMode"
          render={({ field }) => (
            <Form.Item label="网络模式" extra="虚拟机和容器组使用相同IP">
              <Segment options={[{ value: 'Bridge' }]} {...field} />
            </Form.Item>
          )}
        />
      </Form>
    </TeaFormLayout>
  );
};
