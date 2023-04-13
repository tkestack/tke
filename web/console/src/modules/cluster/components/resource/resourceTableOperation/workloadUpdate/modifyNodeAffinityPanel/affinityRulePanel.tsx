import React, { useState } from 'react';
import { Button, Form, Input, InputNumber, Select } from 'tea-component';

export const AffinityRulePanel = () => {
  const [] = useState([
    {
      weight: 1,
      subRules: [
        {
          key: '',
          operator: '',
          value: ''
        }
      ]
    }
  ]);

  return (
    <Form>
      <Form.Item label="权重">
        <InputNumber hideButton size="l" min={0} step={1} />
      </Form.Item>

      <Form.Item label="条件">
        <Input />
        <Select appearance="button" style={{ margin: '0 20px' }} />
        <Input />
        <Button type="icon" icon="delete" />
      </Form.Item>
    </Form>
  );
};
