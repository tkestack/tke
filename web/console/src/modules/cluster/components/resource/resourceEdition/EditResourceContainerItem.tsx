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

import { Bubble, Col, InputNumber, Justify, Row, Text } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, LinkButton } from '../../../../common/components';
import { isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateWorkloadActions } from '../../../actions/validateWorkloadActions';
import { RootProps } from '../../ClusterApp';
import { EditResourceContainerAdvancedPanel } from './EditResourceContainerAdvancedPanel';
import { EditResourceContainerEnvItem } from './EditResourceContainerEnvItem';
import { EditResourceContainerLimitItem } from './EditResourceContainerLimitItem';
import { EditResourceContainerMountItem } from './EditResourceContainerMountItem';
import { RegistrySelectDialog, RegistryTagSelectDialog } from './registrySelect';
import { PermissionProvider } from '@common';

interface ContainerItemProps extends RootProps {
  /** 容器的id */
  cKey?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceContainerItem extends React.Component<ContainerItemProps, {}> {
  state: Readonly<{ tags: any[] }> = { tags: null };
  render() {
    let { actions, subRoot, cKey, clusterVersion, cluster } = this.props,
      { workloadEdit, addons } = subRoot,
      { containers, canAddContainer, volumes, isCanUseGpu } = workloadEdit;

    const container = containers.find(item => item.id === cKey);
    // 选择镜像所需的一些信息
    const selectRegistry = {
      id: uuid(),
      cKey: cKey
    };

    // 是否能够新增容器
    let canAdd = canAddContainer,
      editingContainer = containers.find(c => c.status === 'editing');

    canAdd = isEmpty(editingContainer) || validateWorkloadActions._canAddContainer(editingContainer, volumes);

    // 判断是否能够删除容器
    const canDelete = containers.length > 1;

    // 判断是否能够使用gpu
    const hasGPUManager = !!cluster?.selection?.spec?.features?.gpuType;
    const k8sVersion = clusterVersion.split('.');
    const isK8sOk = +k8sVersion[0] >= 1 && +k8sVersion[1] >= 8;
    const canUseGpu = isK8sOk && isCanUseGpu,
      canUseGpuManager = +k8sVersion[0] >= 1 && +k8sVersion[1] >= 10 && hasGPUManager;

    return (
      container && (
        <div className="run-docker-box">
          <Justify
            right={
              <React.Fragment>
                <LinkButton
                  disabled={!canAdd}
                  tip={t('保存')}
                  errorTip={t('请完成待编辑项')}
                  onClick={() => this._handleSaveContainer(cKey)}
                >
                  <i className="icon-submit-gray" />
                </LinkButton>
                <LinkButton
                  disabled={!canDelete}
                  tip={t('删除')}
                  errorTip={t('不可删除，至少创建一个容器')}
                  onClick={() => actions.editWorkload.deleteContainer(cKey)}
                >
                  <i className="icon-cancel-icon" />
                </LinkButton>
              </React.Fragment>
            }
          />
          <div className="edit-param-list">
            <div className="param-box">
              <div className="param-bd">
                <ul className="form-list fixed-layout" style={{ marginTop: '0' }}>
                  <FormItem label={t('名称')}>
                    <div className={classnames('form-unit', { 'is-error': container.v_name.status === 2 })}>
                      <Bubble
                        placement="bottom"
                        content={container.v_name.status === 2 ? container.v_name.message : null}
                      >
                        <input
                          type="text"
                          className="tc-15-input-text m"
                          placeholder={t('请输入容器名称')}
                          value={container.name}
                          onChange={e => actions.editWorkload.updateContainer({ name: e.target.value }, cKey)}
                          onBlur={e => actions.validate.workload.validateContainerName(e.target.value, cKey)}
                        />
                      </Bubble>
                      <p className="text-label">
                        {t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且不能以分隔符开头或结尾')}
                      </p>
                    </div>
                  </FormItem>
                  <FormItem label={t('镜像')}>
                    <div className={classnames('form-unit', { 'is-error': container.v_registry.status === 2 })}>
                      <Bubble
                        placement="bottom"
                        content={container.v_registry.status === 2 ? container.v_registry.message : null}
                      >
                        <input
                          type="text"
                          className="tc-15-input-text m mr10"
                          value={container.registry}
                          onChange={e => {
                            actions.editWorkload.updateContainer({ registry: e.target.value }, cKey);

                            this.setState({ tags: null });
                          }}
                          onBlur={e => actions.validate.workload.validateRegistrySelection(e.target.value, cKey)}
                        />
                      </Bubble>

                      <RegistrySelectDialog
                        onConfirm={({ registry, tags }) => {
                          actions.editWorkload.updateContainer({ registry }, cKey);

                          this.setState({ tags });
                        }}
                      />
                    </div>
                  </FormItem>
                  <FormItem label={t('镜像版本（Tag）')} className="tag-mod">
                    <div className="tc-15-autocomplete xl">
                      <input
                        type="text"
                        className="tc-15-input-text m mr10"
                        value={container.tag}
                        onChange={e => {
                          actions.editWorkload.updateContainer({ tag: e.target.value }, cKey);
                        }}
                      />

                      {this.state.tags && (
                        <RegistryTagSelectDialog
                          tags={this.state.tags}
                          onConfirm={tag => actions.editWorkload.updateContainer({ tag }, cKey)}
                        />
                      )}
                    </div>
                  </FormItem>

                  <EditResourceContainerMountItem cKey={cKey} />

                  <EditResourceContainerLimitItem cKey={cKey} />

                  <PermissionProvider value="platform.cluster.workload.workload_create_gpu">
                    <FormItem label={t('GPU限制')} isShow={canUseGpu || canUseGpuManager}>
                      {canUseGpuManager ? (
                        <React.Fragment>
                          <Row>
                            <Col span={6}>
                              <Text theme="text">{t('卡数:')}</Text>
                              <InputField
                                type="text"
                                className="tc-15-input-text m"
                                style={{ width: '60px' }}
                                value={container.gpuCore}
                                validator={container.v_gpuCore}
                                tipMode="popup"
                                ops={t('个')}
                                onChange={value =>
                                  actions.editWorkload.updateContainer({ gpuCore: value }, container.id + '')
                                }
                                onBlur={value =>
                                  actions.validate.workload.validateGpuCoreLimit(value, container.id + '')
                                }
                              />
                            </Col>
                            <Col span={8}>
                              <Text theme="text">{t('显存:')}</Text>
                              <InputField
                                type="text"
                                className="tc-15-input-text m"
                                style={{ width: '60px' }}
                                value={container.gpuMem}
                                validator={container.v_gpuMem}
                                tipMode="popup"
                                onChange={value =>
                                  actions.editWorkload.updateContainer({ gpuMem: value }, container.id + '')
                                }
                                ops={'*256MiB'}
                                onBlur={value =>
                                  actions.validate.workload.validateGpuMemLimit(value, container.id + '')
                                }
                              />
                            </Col>
                          </Row>
                          <p className="form-input-help">
                            <Trans>卡数只能填写0.1-1或者1的整数倍。 显存须为256MiB整数倍。</Trans>
                            <br />
                          </p>
                        </React.Fragment>
                      ) : (
                        <InputNumber
                          value={container.gpu}
                          min={0}
                          max={100}
                          step={1}
                          unit={t('个')}
                          onChange={value => actions.editWorkload.updateContainer({ gpu: value }, container.id + '')}
                        />
                      )}
                    </FormItem>
                  </PermissionProvider>

                  <EditResourceContainerEnvItem cKey={cKey} />
                </ul>
              </div>

              <hr className="hr-mod" />

              {container.isOpenAdvancedSetting && <EditResourceContainerAdvancedPanel cKey={cKey} />}
              <a
                href="javascript:;"
                className="more-links-btn"
                onClick={() => actions.editWorkload.toggleAdvancedSetting(cKey)}
              >
                {container.isOpenAdvancedSetting ? t('隐藏高级设置') : t('显示高级设置')}
              </a>
            </div>
          </div>
        </div>
      )
    );
  }

  /** 选择保存按钮 */
  private _handleSaveContainer(cKey: string) {
    const { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers, volumes } = workloadEdit;

    const container = containers.find(c => c.id === cKey);
    // 校验container的所有选项
    actions.validate.workload.validateContainer(container);

    if (validateWorkloadActions._validateContainer(container, volumes, containers)) {
      actions.editWorkload.updateContainer({ status: 'edited' }, cKey);
    }
  }
}
