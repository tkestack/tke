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

import React, { useState } from 'react';
import { Table, Text, Bubble, Icon, Form, Input, TableColumn, Button } from '@tencent/tea-component';
import { t } from '@tencent/tea-app/lib/i18n';
import { v4 as uuidv4 } from 'uuid';
import { Validation } from '@src/modules/common';

interface IItem {
  id: string;
  key: string;
  value: string;
}

interface IAsUserExtraInputProps {
  data?: IItem[];

  onChange(data: IItem[]): void;
}

interface IValidator {
  id: string;

  key: {
    status: 'error' | 'success';
    message?: string;
  };

  value: {
    status: 'error' | 'success';
    message?: string;
  };
}

export const AsUserExtraInput = ({ data = [], onChange }: IAsUserExtraInputProps) => {
  const [validators, setValidators] = useState<IValidator[]>([]);

  function handleChange(id: string, { key, value }: { key?: string; value?: string }) {
    const newData = data.map(item => (item.id === id ? { ...item, key, value } : { ...item }));

    onChange(newData);
  }

  function handleDel(id: string) {
    const newData = data.filter(item => item.id !== id);

    onChange(newData);
  }

  function handleAdd() {
    const newData = [
      ...data,
      {
        id: uuidv4(),
        key: '',
        value: ''
      }
    ];

    onChange(newData);
  }

  function handleBlur() {
    const validators = data.map<IValidator>(item => ({
      id: item.id,
      key: item.key ? { status: 'success' } : { status: 'error', message: t('变量名不能为空') },
      value: item.value ? { status: 'success' } : { status: 'error', message: t('变量值不能为空') }
    }));

    setValidators(validators);
  }

  function getValidatorById(id: string, propName: 'key' | 'value') {
    return validators?.find(item => item.id === id)?.[propName] ?? {};
  }

  const columns: TableColumn[] = [
    {
      key: 'key',
      width: '35%',
      header: (
        <>
          <Text>{t('变量名')}</Text>
          <Bubble placement="top-start" content={t('变量名可以重复')}>
            <Icon type="info" />
          </Bubble>
        </>
      ),

      render: ({ key, id }) => (
        <Form.Item label="" {...getValidatorById(id, 'key')}>
          <Input value={key} onChange={key => handleChange(id, { key })} onBlur={() => handleBlur()} />
        </Form.Item>
      )
    },

    {
      key: 'equal',
      width: '7%',
      header: '',
      render: () => <Text>=</Text>
    },

    {
      key: 'value',
      width: '48%',
      header: t('变量值'),
      render: ({ value, id }) => (
        <Form.Item label="" {...getValidatorById(id, 'value')}>
          <Input value={value} onChange={value => handleChange(id, { value })} onBlur={() => handleBlur()} />
        </Form.Item>
      )
    },

    {
      key: 'action',
      width: '10%',
      header: '',
      render: ({ id }) => <Icon type="close" onClick={() => handleDel(id)} />
    }
  ];

  return (
    <>
      <Table columns={columns} records={data} />

      <Button type="link" onClick={handleAdd} style={{ marginTop: 15 }}>
        {t('手动增加')}
      </Button>
    </>
  );
};
