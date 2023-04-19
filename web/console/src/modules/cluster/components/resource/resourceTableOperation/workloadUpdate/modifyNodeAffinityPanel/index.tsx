import { t } from '@/tencent/tea-app/lib/i18n';
import React, { useState } from 'react';
import { Form, Radio } from 'tea-component';
import { AffinityRulePanel } from './affinityRulePanel';
import { NodeAffinityTypeEnum, TolerationTypeEnum } from './constants';
import { TolerationRulePanel } from './tolerationRulePanel';

export const ModifyNodeAffinityPanel = ({ flag, onSubmit }) => {
  const [state, setState] = useState({
    nodeAffinityType: NodeAffinityTypeEnum.Unset,
    tolerationType: TolerationTypeEnum.UnSet
  });

  function handleSubmit() {}

  return (
    <Form>
      <Form.Item label="节点调度策略">
        <Radio.Group
          value={state.nodeAffinityType}
          onChange={(type: NodeAffinityTypeEnum) => setState(pre => ({ ...pre, nodeAffinityType: type }))}
        >
          <Radio name={NodeAffinityTypeEnum.Unset}>不使用调度策略</Radio>
          <Radio name={NodeAffinityTypeEnum.Rule}>自定义调度规则</Radio>
        </Radio.Group>
      </Form.Item>

      {state.nodeAffinityType === NodeAffinityTypeEnum.Rule && (
        <>
          <Form.Item
            label={t('强制满足条件')}
            tips="调度期间如果满足其中一个亲和性条件则调度到对应node，如果没有节点满足条件则调度失败。"
          >
            <AffinityRulePanel showWeight={false} submiting={flag} onSubmit={handleSubmit} />
          </Form.Item>

          <Form.Item
            label={t('尽量满足条件')}
            tips="调度期间如果满足其中一个亲和性条件则调度到对应node，如果没有节点满足条件则随机调度到任意节点。"
          >
            <AffinityRulePanel submiting={flag} onSubmit={handleSubmit} />
          </Form.Item>
        </>
      )}

      <Form.Item label="容忍调度">
        <Radio.Group
          value={state.tolerationType}
          onChange={(type: TolerationTypeEnum) => setState(pre => ({ ...pre, tolerationType: type }))}
        >
          <Radio name={TolerationTypeEnum.UnSet}>不使用容忍调度</Radio>
          <Radio name={TolerationTypeEnum.Set}>使用容忍调度</Radio>
        </Radio.Group>
      </Form.Item>

      {state.tolerationType === TolerationTypeEnum.Set && (
        <Form.Item>
          <TolerationRulePanel />
        </Form.Item>
      )}
    </Form>
  );
};
