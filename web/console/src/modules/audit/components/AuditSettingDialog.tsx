/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
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
 * 审计配置
 */
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Form, Modal, Button, Input, InputNumber } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { useForm, useFieldArray, Controller, NestedValue } from 'react-hook-form';
import { isEqual, set } from 'lodash';
import { Base64 } from 'js-base64';

import { allActions } from '../actions';
import { configTest, updateStoreConfig, getStoreConfig } from '../WebAPI';

const { useState, useEffect, useReducer } = React;

enum ConnectionState {
  READY = 'ready',
  CONNECTING = 'connecting',
  SUCCESSFUL = 'successful',
  UNSUCCESSFUL = 'unsuccessful'
}

const initialState = {
  // 连接参数配置
  elasticSearch: {
    address: '',
    indices: '',
    password: '',
    reserveDays: 7,
    username: ''
  },
  // 记录连接状态
  connectionState: ConnectionState.READY
};

export const AuditSetting = props => {
  const { onChange } = props;

  const reducer = (state, action) => {
    const { type, payload } = action;
    let nextState;

    switch (type) {
      case 'elasticSearch':
        nextState = set(Object.assign({}, state), 'elasticSearch', { ...payload, username: '', password: '' });
        break;
      case 'address':
        nextState = set(Object.assign({}, state), 'elasticSearch.address', payload);
        break;
      case 'indices':
        nextState = set(Object.assign({}, state), 'elasticSearch.indices', payload);
        break;
      case 'password':
        nextState = set(Object.assign({}, state), 'elasticSearch.password', payload);
        break;
      case 'reserveDays':
        nextState = set(Object.assign({}, state), 'elasticSearch.reserveDays', payload);
        break;
      case 'username':
        nextState = set(Object.assign({}, state), 'elasticSearch.username', payload);
        break;
      case 'connectionState':
        nextState = set(Object.assign({}, state), 'connectionState', payload);
        break;
      default:
        nextState = state;
    }
    if (!isEqual(nextState, state)) {
      onChange(nextState);
    }
    return nextState;
  };

  const {
    control,
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors }
  } = useForm({ mode: 'onTouched' });
  const [auditSetting, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    const _getStoreConfig = async () => {
      const result = await getStoreConfig();
      dispatch({ type: 'elasticSearch', payload: result.elasticSearch });
      setValue('address', result.elasticSearch.address);
      setValue('indices', result.elasticSearch.indices);
      setValue('reserveDays', result.elasticSearch.reserveDays);
    };
    _getStoreConfig();
  }, []);

  const onSubmit = async data => {
    console.log(data);
  };

  /**
   * 检测连接
   * 成功的话更新测试链接状态，成功则完成按钮可用，否则禁用对话框的完成按钮
   */
  const onConnect = async data => {
    console.log('data@onConnect = ', data);
    try {
      dispatch({ type: 'connectionState', payload: ConnectionState.CONNECTING });
      const response = await configTest({
        ...data,
        // ...auditSetting.elasticSearch,
        password: Base64.encode(data.password)
        // password: Base64.encode(auditSetting.elasticSearch.password)
      });
      if (response.code === 0) {
        dispatch({ type: 'connectionState', payload: ConnectionState.SUCCESSFUL });
      } else {
        dispatch({ type: 'connectionState', payload: ConnectionState.UNSUCCESSFUL });
      }
    } catch (err) {
      dispatch({ type: 'connectionState', payload: ConnectionState.UNSUCCESSFUL });
    }
  };

  const getConnectionStatusMessage = state => {
    let result = '';
    switch (state) {
      case 'ready':
        result = '';
        break;
      case 'connecting':
        result = '连接中...';
        break;
      case 'successful':
        result = '连接成功!';
        break;
      case 'unsuccessful':
        result =
          '连接失败！请检查ElasticSearch相关配置，注意开启了用户验证的话需要输入用户名和密码，请注意global集群与ES的连通性。';
        break;
      default:
        result = '';
        break;
    }
    return result;
  };

  const { address, indices, password, reserveDays, username } = auditSetting.elasticSearch;

  return (
    <form id="auditSettingForm" onSubmit={handleSubmit(onSubmit)}>
      <Form>
        <Controller
          name="address"
          control={control}
          rules={{
            required: true,
            pattern:
              /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)/
          }}
          render={({ field: { onChange, onBlur, value } }) => (
            <Form.Item
              label={'ES地址'}
              required
              status={errors.address ? 'error' : undefined}
              message={
                (errors.address?.type === 'required' && '地址不能为空') ||
                (errors.address?.type === 'pattern' && 'Elasticsearch地址格式不正确')
              }
            >
              <Input
                size="full"
                value={value}
                onBlur={onBlur}
                onChange={value => {
                  if (onChange) {
                    onChange(value);
                  }
                  dispatch({ type: 'address', payload: value });
                }}
              />
            </Form.Item>
          )}
        />
        <Controller
          name="indices"
          control={control}
          rules={{ required: true }}
          render={({ field: { onChange, onBlur, value } }) => (
            <Form.Item
              label={'索引'}
              required
              status={errors.indices ? 'error' : undefined}
              message={errors.indices?.type === 'required' && '索引名不能为空'}
            >
              <Input
                size="full"
                value={value}
                onBlur={onBlur}
                onChange={value => {
                  if (onChange) {
                    onChange(value);
                  }
                  dispatch({ type: 'indices', payload: value });
                }}
              />
            </Form.Item>
          )}
        />
        <Controller
          name="reserveDays"
          control={control}
          rules={{ required: true }}
          render={({ field: { onChange, value } }) => (
            <Form.Item label={'保留数据时间'} required>
              <InputNumber
                min={1}
                max={30}
                value={value}
                onChange={value => {
                  if (onChange) {
                    onChange(value);
                  }
                  dispatch({ type: 'reserveDays', payload: value });
                }}
              />
            </Form.Item>
          )}
        />
        <Controller
          name="username"
          control={control}
          defaultValue=""
          render={({ field: { onChange, value } }) => (
            <Form.Item label={'用户名'} message={'仅需要用户验证的ES需要输入用户名和密码'}>
              <Input
                value={value}
                onChange={value => {
                  if (onChange) {
                    onChange(value);
                  }
                  dispatch({ type: 'username', payload: value });
                }}
              />
            </Form.Item>
          )}
        />
        <Controller
          name="password"
          control={control}
          defaultValue=""
          render={({ field: { onChange, value } }) => (
            <Form.Item label={'密码'}>
              <Input
                type="password"
                value={value}
                onChange={value => {
                  if (onChange) {
                    onChange(value);
                  }
                  dispatch({ type: 'password', payload: value });
                }}
              />
            </Form.Item>
          )}
        />
        <Form.Item message={getConnectionStatusMessage(auditSetting.connectionState)}>
          <Button type="primary" onClick={handleSubmit(onConnect)}>
            检测连接
          </Button>
        </Form.Item>
      </Form>
    </form>
  );
};

export const AuditSettingDialog = (props: { isShowing: boolean; toggle: () => void }) => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { isShowing, toggle } = props;
  const [editorState, setEditorState] = useState({ elasticSearch: {}, connectionState: ConnectionState.READY });

  const handleEditorChanged = ({ elasticSearch, connectionState }) => {
    setEditorState({ elasticSearch, connectionState });
  };

  const handleSubmitItem = async () => {
    try {
      const response = await updateStoreConfig({
        ...editorState.elasticSearch,
        password: Base64.encode(editorState.elasticSearch['password'])
      });
      if (response && response.code === 0) {
        // code = 0, message = ''
        toggle();
      }
    } catch (err) {}
  };

  return (
    <Modal visible={isShowing} caption={'审计配置'} onClose={toggle}>
      <Modal.Body>
        <AuditSetting onChange={handleEditorChanged} />
      </Modal.Body>
      <Modal.Footer>
        <Button
          type="primary"
          onClick={handleSubmitItem}
          disabled={editorState.connectionState !== ConnectionState.SUCCESSFUL}
        >
          完成
        </Button>
        <Button onClick={toggle}>取消</Button>
      </Modal.Footer>
    </Modal>
  );
};
