import React, { useState } from 'react';
import { Form, Radio } from 'tea-component';
import { AffinityRulePanel } from './affinityRulePanel';
import { NodeAffinityTypeEnum, TolerationTypeEnum } from './constants';

export const ModifyNodeAffinityPanel = () => {
  const [state, setState] = useState({
    nodeAffinityType: NodeAffinityTypeEnum.Unset,
    tolerationType: TolerationTypeEnum.UnSet
  });

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

      <Form.Item>
        <AffinityRulePanel />
      </Form.Item>

      <Form.Item label="容忍调度">
        <Radio.Group
          value={state.tolerationType}
          onChange={(type: TolerationTypeEnum) => setState(pre => ({ ...pre, tolerationType: type }))}
        >
          <Radio name={TolerationTypeEnum.UnSet}>不使用容忍调度</Radio>
          <Radio name={TolerationTypeEnum.Set}>使用容忍调度</Radio>
        </Radio.Group>
      </Form.Item>
    </Form>
  );
};
