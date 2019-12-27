import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { selectable, removeable } from '@tea/component/table/addons';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { allActions } from '../../actions';
import {
  Modal,
  Button,
  Form,
  Input,
  RadioGroup,
  Radio,
  Transfer,
  SearchBox,
  Table,
  Justify,
  Select
} from '@tea/component';
const { useState, useEffect } = React;
// interface StrategyActionPanelState {
//   name: string;
//   effect: string;
//   selectService: string;
//   action: string[];
//   resource: string;
//   modalVisible?: boolean;
//   disabled?: boolean;
//   description?: string;
//   filterKeyword?: string;
//   messages?: any;
//   [props: string]: any;
// }

const columns = [
  {
    key: 'name',
    header: t('Action名称'),
    render: cvm => <p>{cvm.name}</p>
  },
  {
    key: 'description',
    header: t('描述'),
    width: 100,
    render: cvm => {
      return cvm.description;
    }
  }
];

function SourceTable({ dataSource, action, onChange }) {
  return (
    <Table
      records={dataSource}
      recordKey="name"
      columns={columns}
      addons={[
        selectable({
          value: action,
          onChange,
          rowSelect: true
        })
      ]}
    />
  );
}

function TargetTable({ dataSource, onRemove }) {
  return <Table records={dataSource} recordKey="name" columns={columns} addons={[removeable({ onRemove })]} />;
}

export const StrategyActionPanel = () => {
  const state = useSelector(state => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { strategyList, categoryList } = state;
  const categoryListRecords = categoryList.data.records;

  const [modalVisible, setModalVisible] = useState(false);
  const [modalBtnDisabled, setModalBtnDisabled] = useState(true);
  const [formParamsValue, setFormParamsValue] = useState({
    name: '',
    effect: '',
    resource: '',
    selectService: '',
    description: ''
  });
  const [actionParamsValue, setActionParamsValue] = useState({
    filterKeyword: '',
    action: []
  });
  const [messages, setMessages] = useState({
    name: '',
    effect: '',
    action: '',
    resource: '',
    description: ''
  });

  // 初始化创建弹窗中的服务选项
  useEffect(() => {
    if (!formParamsValue.selectService && categoryListRecords) {
      setFormParamsValue({
        ...formParamsValue,
        selectService: categoryListRecords[0].name
      });
    }
  }, [formParamsValue, categoryListRecords]);

  // 校验form表单
  useEffect(() => {
    const { name, effect, resource } = formParamsValue;
    const { action } = actionParamsValue;
    let disabled = false;
    if (!name || !effect || !action.length || !resource) {
      disabled = true;
    }
    Object.keys(messages).forEach(item => {
      if (messages[item]) {
        disabled = true;
      }
    });
    setModalBtnDisabled(disabled);
  }, [formParamsValue, messages, actionParamsValue]);

  // 获取当前的actionList
  const options = [],
    categoryActions = {};

  categoryListRecords &&
    categoryListRecords.forEach(item => {
      options.push({
        value: item.metadata.name,
        text: item.Spec.displayName
      });
      categoryActions[item.metadata.name] = Object.values(item.Spec.actions);
    });
  const actionList = categoryActions[formParamsValue.selectService] || [];

  return (
    <React.Fragment>
      <Table.ActionPanel>
        <Justify
          left={
            <Button type="primary" onClick={_open}>
              {t('新建')}
            </Button>
          }
          right={
            <React.Fragment>
              <SearchBox
                value={strategyList.query.keyword || ''}
                onChange={actions.strategy.changeKeyword}
                onSearch={actions.strategy.performSearch}
                onClear={() => {
                  actions.strategy.performSearch('');
                }}
                placeholder={t('请输入策略名称')}
              />
            </React.Fragment>
          }
        />
      </Table.ActionPanel>
      <Modal visible={modalVisible} caption={t('新建策略')} size="l" onClose={_close}>
        <Modal.Body>
          <Form className="add-strategy-form">
            <Form.Item
              label={t('策略名称')}
              required
              status={messages.name ? 'error' : name ? 'success' : undefined}
              message={messages.name ? t(messages.name) : t('长度需要小于256个字符')}
            >
              <Input
                placeholder={t('请输入策略名称')}
                style={{ width: '350px' }}
                defaultValue={name}
                onChange={value => {
                  let msg = '';
                  if (!value) {
                    msg = '请输入策略名称';
                  } else if (value.length > 255) {
                    msg = '长度需要小于256个字符';
                  }
                  setFormParamsValue({ ...formParamsValue, name: value });
                  setMessages({ ...messages, name: msg });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('效果')}
              required
              status={messages.effect ? 'error' : formParamsValue.effect ? 'success' : undefined}
              message={messages.effect ? t(messages.effect) : t('请选择效果')}
            >
              <RadioGroup
                value={formParamsValue.effect}
                onChange={value => {
                  setFormParamsValue({ ...formParamsValue, effect: value });
                  setMessages({ ...messages, effect: value ? '' : '请选择效果' });
                }}
              >
                <Radio name="allow">
                  <Trans>允许</Trans>
                </Radio>
                <Radio name="deny">
                  <Trans>拒绝</Trans>
                </Radio>
              </RadioGroup>
            </Form.Item>
            <Form.Item label="服务" required>
              <Select
                type="native"
                size="m"
                options={options}
                value={formParamsValue.selectService}
                onChange={value => {
                  // 调用action获取接口
                  setFormParamsValue({ ...formParamsValue, selectService: value });
                }}
                placeholder="请选择服务"
              />
            </Form.Item>
            <Form.Item label={t('操作')} required message={t('请选择Action')}>
              <Transfer
                leftCell={
                  <Transfer.Cell
                    title={t('请选择Action')}
                    tip={t('支持按住 shift 键进行多选')}
                    header={
                      <SearchBox
                        value={actionParamsValue.filterKeyword}
                        onChange={value => setActionParamsValue({ ...actionParamsValue, filterKeyword: value })}
                      />
                    }
                  >
                    <SourceTable
                      dataSource={actionList.filter(i => i.name.includes(actionParamsValue.filterKeyword))}
                      action={actionParamsValue.action}
                      onChange={keys => {
                        setActionParamsValue({ ...actionParamsValue, action: keys });
                        setMessages({ ...messages, action: keys.length ? '' : '请选择Action' });
                      }}
                    />
                  </Transfer.Cell>
                }
                rightCell={
                  <Transfer.Cell title={t(`已选择 (${actionParamsValue.action.length})个`)}>
                    <TargetTable
                      dataSource={actionList.filter(i => actionParamsValue.action.includes(i.name))}
                      onRemove={key => {
                        const keys = actionParamsValue.action.filter(i => i !== key);
                        setActionParamsValue({ ...actionParamsValue, action: keys });
                        setMessages({ ...messages, action: keys.length ? '' : '请选择Action' });
                      }}
                    />
                  </Transfer.Cell>
                }
              />
            </Form.Item>
            <Form.Item
              label={t('资源')}
              required
              status={messages.resource ? 'error' : formParamsValue.resource ? 'success' : undefined}
              message={
                messages.resource
                  ? t(messages.resource)
                  : t(
                      '采用分段式描述方式：key1:val1/key2:val2/*，支持*模糊匹配语法，如cluster:cls-123/deployment:deploy-123/*'
                    )
              }
            >
              <Input
                placeholder={t('eg. cluster:cls-123/deployment:deploy-123/*')}
                style={{ width: '350px' }}
                defaultValue={formParamsValue.resource}
                onChange={value => {
                  let msg = '';
                  if (!value) {
                    msg = '请输入资源名称';
                  }
                  setFormParamsValue({ ...formParamsValue, resource: value });
                  setMessages({ ...messages, resource: msg });
                }}
              />
            </Form.Item>
            <Form.Item
              label={t('描述')}
              status={messages.description ? 'error' : undefined}
              message={t('描述不能超过255个字符')}
            >
              <Input
                multiline
                placeholder={t('介绍一下这个策略')}
                style={{ width: '350px' }}
                onChange={value => {
                  let msg = '';
                  if (value && value.length > 255) {
                    msg = '描述不能超过255个字符';
                  }
                  setFormParamsValue({ ...formParamsValue, description: value });
                  setMessages({ ...messages, description: msg });
                }}
              />
            </Form.Item>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Form.Action>
            <Button disabled={modalBtnDisabled} type="primary" onClick={_onSubmit}>
              <Trans>保存</Trans>
            </Button>
            <Button onClick={_close}>
              <Trans>取消</Trans>
            </Button>
          </Form.Action>
        </Modal.Footer>
      </Modal>
    </React.Fragment>
  );

  function _close() {
    setModalVisible(false);
  }
  function _open() {
    actions.strategy.getCategories.fetch();
    setModalVisible(true);
  }
  function _onSubmit() {
    const { name, effect, resource, selectService, description } = formParamsValue;
    const { action } = actionParamsValue;
    const strategyInfo = {
      id: uuid(),
      spec: {
        displayName: name,
        category: selectService,
        description,
        statement: {
          resources: [resource],
          effect,
          actions: action
        }
      }
    };

    console.log(strategyInfo);
    actions.strategy.addStrategy.start([strategyInfo]);
    actions.strategy.addStrategy.perform();
    setModalVisible(false);
    setModalBtnDisabled(true);
    setFormParamsValue({ name: '', effect: '', resource: '', selectService: '', description: '' });
    setActionParamsValue({ filterKeyword: '', action: [] });
  }
};
