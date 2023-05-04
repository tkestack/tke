import { t } from '@/tencent/tea-app/lib/i18n';
import { zodResolver } from '@hookform/resolvers/zod';
import { TeaFormLayout } from '@src/modules/common/layouts/TeaFormLayout';
import React, { useEffect } from 'react';
import { Controller, FormProvider, useForm } from 'react-hook-form';
import { Button, Form, Radio } from 'tea-component';
import { IModifyPanelProps, WorkloadKindEnum } from '../constants';
import { AffinityRulePanel } from './affinityRulePanel';
import {
  AffinityTypeEnum,
  NodeAffinityFormType,
  NodeAffinityTypeEnum,
  TolerationEffectEnum,
  TolerationOperatorEnum,
  TolerationTypeEnum,
  defaultNodeAffinityFormData,
  nodeAffinitySchema
} from './constants';
import { TolerationRulePanel } from './tolerationRulePanel';

export const ModifyNodeAffinityPanel = ({ title, baseInfo, onCancel, kind, onUpdate, resource }: IModifyPanelProps) => {
  const isCronjob = kind === WorkloadKindEnum.Cronjob;

  const useFormReturn = useForm<NodeAffinityFormType>({
    mode: 'onBlur',
    defaultValues: defaultNodeAffinityFormData,
    resolver: zodResolver(nodeAffinitySchema)
  });

  const { handleSubmit, control, watch, reset } = useFormReturn;

  useEffect(() => {
    if (!resource) return;

    let originNodeAffinityInfo = resource?.spec?.template?.spec?.affinity?.nodeAffinity;

    let tolerationInfo = resource?.spec?.template?.spec?.tolerations;

    if (isCronjob) {
      originNodeAffinityInfo = resource?.spec?.jobTemplate?.spec?.template?.spec?.affinity?.nodeAffinity;

      tolerationInfo = resource?.spec?.jobTemplate?.spec?.template?.spec?.tolerations;
    }

    const forceAffinityRules =
      originNodeAffinityInfo?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms?.map(item => ({
        weight: 1,
        subRules:
          item?.matchExpressions?.map(r => ({
            key: r?.key ?? '',
            operator: r?.operator,
            value: r?.values?.join(';') ?? ''
          })) ?? []
      })) ?? [];

    const attemptAffinityRules =
      originNodeAffinityInfo?.preferredDuringSchedulingIgnoredDuringExecution?.map(item => ({
        weight: item?.weight ?? '',
        subRules:
          item?.preference?.matchExpressions?.map(r => ({
            key: r?.key ?? '',
            operator: r?.operator,
            value: r?.values?.join(';') ?? ''
          })) ?? []
      })) ?? [];

    const data: NodeAffinityFormType = {
      nodeAffinityType: originNodeAffinityInfo ? NodeAffinityTypeEnum.Rule : NodeAffinityTypeEnum.Unset,
      affinityRules: {
        force: forceAffinityRules,
        attempt: attemptAffinityRules
      },

      tolerationType: tolerationInfo ? TolerationTypeEnum.Set : TolerationTypeEnum.UnSet,

      tolerationRules:
        tolerationInfo?.map(item => ({
          key: item?.key ?? '',
          operator: item?.operator,
          value: item?.value ?? '',
          effect: item?.effect ?? TolerationEffectEnum.All,
          time: item?.tolerationSeconds ?? 0
        })) ?? []
    };

    reset(data);
  }, [resource, isCronjob, reset]);

  function onSubmit({ nodeAffinityType, affinityRules, tolerationType, tolerationRules }: NodeAffinityFormType) {
    const templateContent = {
      spec: {
        tolerations:
          tolerationType === TolerationTypeEnum.Set
            ? tolerationRules.map(({ key, operator, value, effect, time }) => ({
                key: key || undefined,
                operator,
                value: operator === TolerationOperatorEnum.Exists ? undefined : value,
                effect: effect === TolerationEffectEnum.All ? undefined : effect,
                tolerationSeconds: effect === TolerationEffectEnum.NoExecute ? time : undefined
              }))
            : null,

        affinity: {
          nodeAffinity:
            nodeAffinityType === NodeAffinityTypeEnum.Unset
              ? null
              : {
                  requiredDuringSchedulingIgnoredDuringExecution: affinityRules.force.length
                    ? {
                        nodeSelectorTerms: affinityRules.force.map(({ subRules }) => ({
                          matchExpressions: subRules.map(({ key, operator, value }) => ({
                            key,
                            operator,
                            values: value ? value.split(';') : undefined
                          }))
                        }))
                      }
                    : null,

                  preferredDuringSchedulingIgnoredDuringExecution: affinityRules.attempt.length
                    ? affinityRules.attempt.map(({ weight, subRules }) => ({
                        weight,
                        preference: {
                          matchExpressions: subRules.map(({ key, operator, value }) => ({
                            key,
                            operator,
                            values: value ? value.split(';') : undefined
                          }))
                        }
                      }))
                    : null
                }
        }
      }
    };

    const jsonData = {
      spec: {
        template: isCronjob ? undefined : templateContent,
        jobTemplate: isCronjob
          ? {
              spec: {
                template: templateContent
              }
            }
          : undefined
      }
    };

    onUpdate(jsonData);
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

      <FormProvider {...useFormReturn}>
        <Form>
          <Controller
            control={control}
            name="nodeAffinityType"
            render={({ field }) => (
              <Form.Item label="节点调度策略">
                <Radio.Group {...field}>
                  <Radio name={NodeAffinityTypeEnum.Unset}>不使用调度策略</Radio>
                  <Radio name={NodeAffinityTypeEnum.Rule}>自定义调度规则</Radio>
                </Radio.Group>
              </Form.Item>
            )}
          />

          {watch('nodeAffinityType') === NodeAffinityTypeEnum.Rule && (
            <>
              <Form.Item
                label={t('强制满足条件')}
                tips="调度期间如果满足其中一个亲和性条件则调度到对应node，如果没有节点满足条件则调度失败。"
              >
                <AffinityRulePanel subName={AffinityTypeEnum.Force} />
              </Form.Item>

              <Form.Item
                label={t('尽量满足条件')}
                tips="调度期间如果满足其中一个亲和性条件则调度到对应node，如果没有节点满足条件则随机调度到任意节点。"
              >
                <AffinityRulePanel subName={AffinityTypeEnum.Attempt} />
              </Form.Item>
            </>
          )}

          <Controller
            control={control}
            name="tolerationType"
            render={({ field }) => (
              <Form.Item label="容忍调度">
                <Radio.Group {...field}>
                  <Radio name={TolerationTypeEnum.UnSet}>不使用容忍调度</Radio>
                  <Radio name={TolerationTypeEnum.Set}>使用容忍调度</Radio>
                </Radio.Group>
              </Form.Item>
            )}
          />

          {watch('tolerationType') === TolerationTypeEnum.Set && (
            <Form.Item>
              <TolerationRulePanel />
            </Form.Item>
          )}
        </Form>
      </FormProvider>
    </TeaFormLayout>
  );
};
