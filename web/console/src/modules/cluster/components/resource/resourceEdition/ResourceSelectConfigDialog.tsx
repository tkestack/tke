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
import * as classnames from 'classnames';
import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Button, Modal, Radio } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem } from '../../../../common/components';
import { FormLayout } from '../../../../common/layouts';
import { initValidator, Validation } from '../../../../common/models';
import { allActions } from '../../../actions';
import { ConfigItems } from '../../../models';
import { RootProps } from '../../ClusterApp';

/** configMap选项列表 all | optional */
const configMapKeyTypeList = [
  {
    value: 'all',
    label: t('全部')
  },
  {
    value: 'optional',
    label: t('指定部分Key')
  }
];

interface ResourceSelectConfigMapDialogState {
  /** 子路径 type 为 all模式下使用 */
  subPath?: string;

  /** subPath的校验 */
  v_subPath?: Validation;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceSelectConfigDialog extends React.Component<RootProps, ResourceSelectConfigMapDialogState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      subPath: '',
      v_subPath: initValidator
    };
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let { isShowConfigDialog, configEdit } = nextProps.subRoot.workloadEdit;

    let oldIsShowDialog = this.props.subRoot.workloadEdit.isShowConfigDialog,
      newIsShowDialog = isShowConfigDialog;

    // 重新打开的时候，需要进行一些操作，如果是指定某些key的话，
    if (oldIsShowDialog === false && newIsShowDialog === true && configEdit.keyType === 'optional') {
      let { volumes, currentEditingVolumeId } = nextProps.subRoot.workloadEdit;
      let currentVolume = volumes.find(v => v.id === currentEditingVolumeId);
      let configKeys = currentVolume.volumeType === 'secret' ? currentVolume.secretKey : currentVolume.configKey;
      nextProps.actions.editWorkload.config.initConfigItemsByVolumes(configKeys);
    }
  }

  render() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { isShowConfigDialog, currentEditingVolumeId, configEdit, volumes } = workloadEdit;

    let currentVolume = volumes.find(v => v.id === currentEditingVolumeId);
    let isSecret = currentVolume && currentVolume.volumeType === 'secret' ? true : false;

    // 如果不需要展示configmap的弹窗选择
    if (!isShowConfigDialog) {
      return <noscript />;
    }

    const commonAction = () => {
      // 关闭窗口
      actions.editWorkload.toggleConfigDialog();
      // 清空configMapItems
      actions.editWorkload.config.initConfigItems([], configEdit.keyType);
    };

    const cancel = () => {
      commonAction();
    };

    const perform = () => {
      if (configEdit.configSelection.length) {
        let configItems: ConfigItems[] = [];
        // 校验所有的内容是否ok，如果是选择全部key的话
        if (configEdit.keyType === 'optional') {
          let isFormOk = this._validateAllConfigItems(configEdit.configItems);

          if (isFormOk) {
            commonAction();

            // 如果都ok，则拼configItems的数据出来
            configItems = configEdit.configItems.map(item => {
              let tmp: ConfigItems = {
                id: uuid(),
                configKey: item.configKey,
                path: item.path,
                mode: item.mode
              };
              return tmp;
            });
          }
        } else {
          commonAction();
        }

        let updateContent = {};
        let name = configEdit.configSelection[0].metadata.name;
        if (isSecret) {
          updateContent['secretName'] = name;
          updateContent['secretKey'] = configItems;
        } else {
          updateContent['configName'] = name;
          updateContent['configKey'] = configItems;
        }
        actions.editWorkload.updateVolume(updateContent, currentEditingVolumeId);
      }
    };

    /** 展示configName */
    let configName = isSecret ? 'Secret' : 'ConfigMap';

    /** 渲染configMap列表 */
    let configOptions = [],
      configListArr = isSecret ? configEdit.secretList.data : configEdit.configList.data;
    configOptions = configListArr.recordCount
      ? configListArr.records.map((configItem, index) => {
          return (
            <option key={index} value={configItem.metadata.name}>
              {configItem.metadata.name}
            </option>
          );
        })
      : [];
    configOptions.unshift(
      <option value="">
        {configListArr.recordCount
          ? t('请选择{{configName}}', { configName })
          : t('无可用{{configName}}', { configName })}
      </option>
    );

    let defaultConfigSelect = isSecret ? currentVolume.secretName : currentVolume.configName;
    if (configEdit.configSelection.length) {
      // 主要是configMap和 secret相互切换的时候，得判断当前的selection是不是在 configList当中
      let hasThisSelection = configListArr.records.find(item => item.id === configEdit.configSelection[0].id)
        ? true
        : false;
      defaultConfigSelect = hasThisSelection ? configEdit.configSelection[0].metadata.name : '';
    }

    let disabledAddItem = configEdit.configItems.length < configEdit.configKeys.length ? false : true;

    return (
      <Modal visible={true} caption={t('设置') + configName} onClose={cancel} disableEscape={true} size={750}>
        <Modal.Body>
          <FormLayout>
            <div className="param-box server-update add">
              <ul className="form-list jiqun fixed-layout">
                <FormItem label={t('选择') + configName}>
                  <select
                    className="tc-15-select m"
                    style={{ marginRight: '6px' }}
                    value={defaultConfigSelect}
                    onChange={e => {
                      this._handleConfigMapSelect(e.target.value, isSecret);
                    }}
                  >
                    {configOptions}
                  </select>
                </FormItem>
                <FormItem label={t('选项')}>
                  <Radio.Group
                    value={configEdit.keyType}
                    onChange={value => actions.editWorkload.config.changeKeyType(value)}
                  >
                    {configMapKeyTypeList.map((item, cIndex) => {
                      return (
                        <Radio key={cIndex} name={item.value}>
                          {item.label}
                        </Radio>
                      );
                    })}
                  </Radio.Group>
                </FormItem>
                <FormItem label="Items" isShow={configEdit.keyType === 'optional'}>
                  {this._renderConfigMapItems()}
                  <Bubble placement="left" content={disabledAddItem ? <p>{configName + t('无更多可用Key')}</p> : null}>
                    <a
                      href="javascript:;"
                      style={{ verticalAlign: 'middle' }}
                      className={classnames('more-links-btn', { disabled: disabledAddItem })}
                      onClick={!disabledAddItem && actions.editWorkload.config.addConfigItem}
                    >
                      {t('添加Item')}
                    </a>
                    {disabledAddItem && <i className="n-error-icon" style={{ verticalAlign: 'middle' }} />}
                  </Bubble>
                </FormItem>
              </ul>
            </div>
          </FormLayout>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={perform}>
            {t('确认')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }

  /** 当configKeyType 为 指定部分Key  */
  private _renderConfigMapItems() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { configEdit } = workloadEdit;

    /** 子路径的相关提示 */
    let pathTips = t('向特定路径挂在，如挂载点是 /data/config，子路径是dev，最终会存储在/data/config/dev下');

    /** 类型为 特定key的时候，渲染items的列表 */
    let itemsOptions = configEdit.configKeys.map((item, index) => {
      return (
        <option key={index} value={item}>
          {item}
        </option>
      );
    });

    let content: any;

    if (configEdit.configItems.length) {
      let isShow = true;

      content = (
        <div>
          {configEdit.configItems.map((item, cIndex) => {
            return (
              <div style={{ marginBottom: '5px' }} key={cIndex}>
                <div
                  style={{ display: 'inline-block', verticalAlign: 'middle' }}
                  className={classnames('form-unit', { 'is-error': item.v_configKey.status === 2 })}
                >
                  <Bubble placement="bottom" content={item.v_configKey.status === 2 ? item.v_configKey.message : null}>
                    <select
                      className="tc-15-select m"
                      style={{ marginRight: '6px', maxWidth: '200px' }}
                      value={item.configKey}
                      onChange={e => {
                        this._handleConfigMapKeySelect(e.target.value, item.id + '');
                      }}
                    >
                      {itemsOptions}
                    </select>
                  </Bubble>
                </div>
                <div
                  style={{ display: 'inline-block', verticalAlign: 'middle' }}
                  className={classnames('form-unit', { 'is-error': item.v_path.status === 2 })}
                >
                  <Bubble placement="bottom" content={item.v_path.status === 2 ? item.v_path.message : null}>
                    <input
                      type="text"
                      style={{ width: '150px', marginRight: '6px' }}
                      className="tc-15-input-text m"
                      placeholder={t('请输入子路径，eg: dev')}
                      value={item.path}
                      onChange={e =>
                        actions.editWorkload.config.updateConfigItems({ path: e.target.value }, item.id + '')
                      }
                      onBlur={e => this._handleItemsSubPathBlur(e.target.value, item.id + '')}
                    />
                  </Bubble>
                </div>
                <div
                  style={{ display: 'inline-block', verticalAlign: 'middle' }}
                  className={classnames('form-unit', { 'is-error': item.v_mode.status === 2 })}
                >
                  <Bubble placement="bottom" content={item.v_mode.status === 2 ? item.v_mode.message : null}>
                    <input
                      type="text"
                      style={{ width: '150px', marginRight: '10px' }}
                      className="tc-15-input-text m"
                      placeholder={t('文件权限，如0644')}
                      value={item.mode}
                      onChange={e =>
                        actions.editWorkload.config.updateConfigItems({ mode: e.target.value }, item.id + '')
                      }
                      onBlur={e => this._handleItemsModeBlur(e.target.value, item.id + '')}
                    />
                  </Bubble>
                </div>
                {cIndex > 0 ? (
                  <a
                    href="javascript:;"
                    onClick={() => {
                      actions.editWorkload.config.deleteConfigMapItem(item.id + '');
                    }}
                  >
                    <i className="icon-cancel-icon" />
                  </a>
                ) : (
                  <Bubble placement="bottom" content={t('不可删除，至少指定一个Key')}>
                    <a href="javascript:;" className="disabled">
                      <i className="icon-cancel-icon" />
                    </a>
                  </Bubble>
                )}
              </div>
            );
          })}
          <p className="text-label">{pathTips}</p>
        </div>
      );
    } else {
      content = <noscript />;
    }

    return content;
  }

  /** 处理选择configMap */
  private _handleConfigMapSelect(name: string, isSecret: boolean = false) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { configEdit } = workloadEdit,
      { configList, secretList } = configEdit;

    let configSelect = [];
    let configMapInfo = isSecret
      ? secretList.data.records.find(item => item.metadata.name === name)
      : configList.data.records.find(item => item.metadata.name === name);
    if (configMapInfo) {
      configSelect.push(configMapInfo);
    }
    actions.editWorkload.config.selectConfig(configSelect);
  }

  /** items的选择 */
  private _handleConfigMapKeySelect(keyName: string, cId: string) {
    // 更新其状态到 store当中
    this.props.actions.editWorkload.config.updateConfigItems({ configKey: keyName }, cId);
  }

  /** item当中blur之后的 */
  private _handleItemsSubPathBlur(path: string, cId: string) {
    // 校验是否合法
    let result = this._validateSubPath(path);

    // 更新其状态到 store当中
    this.props.actions.editWorkload.config.updateConfigItems({ v_path: result }, cId);
  }

  /** mode blur之后，相关的一些操作 */
  private _handleItemsModeBlur(mode: string, cId: string) {
    // 校验是否合法
    let result = this._validateMode(mode);

    // 更新其状态到store当中
    this.props.actions.editWorkload.config.updateConfigItems({ v_mode: result }, cId);
  }

  /** 校验configKey的mode是否正确 */
  private _validateMode(mode: string) {
    let status = 0,
      message = '',
      reg = /^\d+$/;

    if (!reg.test) {
      status = 2;
      message = t('数据格式不正确');
    } else if (+mode < 0 || +mode > 777) {
      status = 2;
      message = t('mode取值范围为 0-0777');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  }

  /** 校验configKey的选择是否正确 */
  private _validateConfigKey(keyName: string, configMapItems: ConfigItems[]) {
    let status = 0,
      message = '';

    if (configMapItems.filter(item => item.configKey === keyName).length > 1) {
      status = 2;
      message = t('Key不可重复');
    } else {
      status = 1;
      message = '';
    }
    return { status, message };
  }

  /** 校验子路径是否正确 */
  private _validateSubPath(path: string) {
    let status = 0,
      message = '',
      reg = /^[\w]/;

    if (!path) {
      status = 2;
      message = t('子路径不能为空');
    } else if (!reg.test(path)) {
      status = 2;
      message = t('子路径格式不正确，且不能包含 ..');
    } else {
      status = 1;
      message = '';
    }

    return { status, message };
  }

  /** 校验所有的configItems是否正确 */
  private _validateAllConfigItems(configItems: ConfigItems[]) {
    let result = true;
    configItems.forEach(item => {
      let configKeyResult = this._validateConfigKey(item.configKey, configItems),
        pathResult = this._validateSubPath(item.path);

      // 更新store的内容
      this.props.actions.editWorkload.config.updateConfigItems(
        { v_configKey: configKeyResult, v_path: pathResult },
        item.id + ''
      );

      result = result && configKeyResult.status === 1 && pathResult.status === 1;
    });
    return result;
  }
}
