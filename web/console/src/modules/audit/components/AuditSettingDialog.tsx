/**
 * 审计配置
 */
import * as React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Form, Modal, Button, Input, InputNumber } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { Form as FinalForm, Field } from 'react-final-form';
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

const getStatus = meta => {
  console.log('meta@getStatus = ', meta);
  if (meta.active && meta.validating) {
    return 'validating';
  }
  if (!meta.touched) {
    return null;
  }
  return meta.error ? 'error' : 'success';
};

const required = value => (value ? undefined : '内容不能为空');

const composeValidators = (...validators) => value =>
  validators.reduce((error, validator) => error || validator(value), undefined);

const regexURL = new RegExp(
  /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)/
);
const isURL = value => {
  return regexURL.test(value);
};
const mustBeURL = value => (!isURL(value) ? '请输入有效的URL地址格式' : undefined);

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

  const [auditSetting, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    const _getStoreConfig = async () => {
      const result = await getStoreConfig();
      dispatch({ type: 'elasticSearch', payload: result.elasticSearch });
    };
    _getStoreConfig();
  }, []);

  const submit = async values => {
    try {
      // let payload = this.stateToPayload(values);
      // let response = await createRule(clusterName, payload);
    } catch (err) {
      // message.error(err)
    }
  };

  /**
   * 检测连接
   * 成功的话更新测试链接状态，成功则完成按钮可用，否则禁用对话框的完成按钮
   */
  const connect = async () => {
    try {
      dispatch({ type: 'connectionState', payload: ConnectionState.CONNECTING });
      let response = await configTest({
        ...auditSetting.elasticSearch,
        password: Base64.encode(auditSetting.elasticSearch.password)
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

  return (
    <FinalForm
      onSubmit={submit}
      initialValues={{ ...auditSetting.elasticSearch }}
      render={({ form, handleSubmit }) => (
        <form id="auditSettingForm" onSubmit={handleSubmit}>
          <Form>
            <Field name="address" validate={composeValidators(required, mustBeURL)}>
              {({ input, meta, ...rest }) => (
                <Form.Item
                  label={'ES地址'}
                  required
                  status={getStatus(meta)}
                  message={getStatus(meta) === 'error' ? meta.error : '例如，http://10.0.0.1:9200'}
                >
                  <Input
                    size="full"
                    value={input.value}
                    onChange={value => {
                      if (input.onChange) {
                        input.onChange(value);
                      }
                      dispatch({ type: 'address', payload: value });
                    }}
                  />
                </Form.Item>
              )}
            </Field>
            <Field name="indices" validate={required}>
              {({ input, meta, ...rest }) => (
                <Form.Item
                  label={'索引'}
                  required
                  status={getStatus(meta)}
                  message={getStatus(meta) === 'error' && meta.error}
                >
                  <Input
                    size="full"
                    value={input.value}
                    onChange={value => {
                      dispatch({ type: 'indices', payload: value });
                      if (input.onChange) {
                        input.onChange(value);
                      }
                    }}
                  />
                </Form.Item>
              )}
            </Field>
            <Field name="reserveDays" validate={required}>
              {({ input, meta, ...rest }) => (
                <Form.Item
                  label={'保留数据时间'}
                  required
                  status={getStatus(meta)}
                  message={getStatus(meta) === 'error' && meta.error}
                >
                  <InputNumber
                    min={1}
                    max={30}
                    value={input.value}
                    onChange={value => dispatch({ type: 'reserveDays', payload: value })}
                  />
                </Form.Item>
              )}
            </Field>
            <Field name="username">
              {({ input, meta, ...rest }) => (
                <Form.Item label={'用户名'} message={'仅需要用户验证的ES需要输入用户名和密码'}>
                  <Input value={input.value} onChange={value => dispatch({ type: 'username', payload: value })} />
                </Form.Item>
              )}
            </Field>
            <Field name="password">
              {({ input, meta, ...rest }) => (
                <Form.Item label={'密码'}>
                  <Input
                    type="password"
                    value={input.value}
                    onChange={value => dispatch({ type: 'password', payload: value })}
                  />
                </Form.Item>
              )}
            </Field>
            <Form.Item message={getConnectionStatusMessage(auditSetting.connectionState)}>
              <Button
                type="primary"
                onClick={() => {
                  connect();
                }}
              >
                检测连接
              </Button>
            </Form.Item>
          </Form>
        </form>
      )}
    />
  );
};

export const AuditSettingDialog = (props: { isShowing: boolean; toggle: () => void }) => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);
  const { isShowing, toggle } = props;
  const [editorState, setEditorState] = useState({ elasticSearch: {}, connectionState: ConnectionState.READY });

  // const getTips = () => (
  //   <Alert>
  //     审计模块为平台提供了操作记录，管理员可以在运维中心里查询审计日志。TKEStack将在ES里生成名为auditevent的index存储审计记录。
  //   </Alert>
  // );

  const handleEditorChanged = ({ elasticSearch, connectionState }) => {
    setEditorState({ elasticSearch, connectionState });
  };

  const handleSubmitItem = async () => {
    try {
      let response = await updateStoreConfig({
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
        {/*{getTips()}*/}
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
