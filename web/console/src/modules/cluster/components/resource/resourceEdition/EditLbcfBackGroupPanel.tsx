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

import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, Button, ContentView, FormItem, Input, Justify, List, Select, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, deepClone, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { resourceConfig } from '../../../../../../config';
import { validateLbcfActions } from '../../../../../modules/cluster/actions/validateLbcfActions';
import { TipInfo } from '../../../../common/components';
import { getWorkflowError } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { CreateResource, LbcfBGJSONYaml, MergeType } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { EditLbcfBackGroupItemPanel } from './EditLbcfBackGroupItemPanel';
import { reduceNs } from '../../../../../../helpers';
import { BackendType } from '@src/modules/cluster/constants/Config';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditLbcfBackGroupPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let {
        actions,
        route,
        subRoot: {
          resourceOption: { ffResourceList }
        }
      } = this.props,
      urlParams = router.resolve(route);
    actions.lbcf.selectLbcfNamespace(route.queries['np']);
    let resourceIns = route.queries['resourceIns'];
    // 这里是从列表页进入的时候，需要去初始化 workloadEdit当中的内容，如果是直接在当前页面刷新的话，会去拉取列表，在fetchResource之后，会初始化
    if (ffResourceList.list.data.recordCount && urlParams['tab'] === 'updateBG') {
      let finder = ffResourceList.list.data.records.find(item => item.metadata.name === resourceIns);
      finder && actions.lbcf.initGameBGEdition(finder.spec.backGroups);
    }
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.lbcf.clearEdition();
  }

  render() {
    let { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { lbcfEdit, modifyMultiResourceWorkflow, updateMultiResource } = subRoot,
      { namespace, lbcfBackGroupEditions } = lbcfEdit;

    let canEdit = lbcfBackGroupEditions.every(item => !item.onEdit);
    let cantDelete = lbcfBackGroupEditions.length === 1;

    let mode = urlParams['tab'] === 'createBG' ? 'create' : 'update';

    let failed =
      mode === 'create'
        ? modifyMultiResourceWorkflow.operationState === OperationState.Done &&
          !isSuccessWorkflow(modifyMultiResourceWorkflow)
        : updateMultiResource.operationState === OperationState.Done && !isSuccessWorkflow(updateMultiResource);
    return (
      <ContentView>
        <ContentView.Body>
          <FormPanel>
            <FormPanel.Item label={t('命名空间')} text>
              {namespace}
            </FormPanel.Item>
            <FormPanel.Item label={t('后端配置')}>
              {lbcfBackGroupEditions.map(backgroup => {
                let { id, name, onEdit } = backgroup;
                if (onEdit) {
                  return <EditLbcfBackGroupItemPanel backGroupId={id + ''} backGroupmode={mode} />;
                } else {
                  return (
                    <FormPanel
                      fixed
                      key={id}
                      isNeedCard={false}
                      labelStyle={{
                        minWidth: 430
                      }}
                    >
                      <FormPanel.Item label={<Text theme="strong">{name}</Text>}>
                        <Justify
                          right={
                            <React.Fragment>
                              <Bubble content={canEdit ? null : t('请先完成待编辑项')}>
                                <Button
                                  icon="pencil"
                                  disabled={!canEdit}
                                  onClick={() => {
                                    actions.lbcf.changeBackgroupEditStatus(id + '', true);
                                  }}
                                />
                              </Bubble>
                              {mode === 'create' ? (
                                <Bubble content={cantDelete ? t('至少保留一项') : null}>
                                  <Button
                                    icon="close"
                                    disabled={cantDelete}
                                    onClick={() => actions.lbcf.deleteLbcfBackGroup(id + '')}
                                  />
                                </Bubble>
                              ) : (
                                <noscript />
                              )}
                            </React.Fragment>
                          }
                        />
                      </FormPanel.Item>
                    </FormPanel>
                  );
                }
              })}
              {mode === 'create' ? (
                <div
                  style={{
                    lineHeight: '44px',
                    border: '1px dashed #ddd',
                    marginTop: '10px',
                    fontSize: '12px',
                    textAlign: 'center'
                  }}
                >
                  <Bubble content={canEdit ? null : t('请先完成待编辑项')}>
                    <Button
                      type="link"
                      disabled={!canEdit}
                      onClick={() => {
                        actions.lbcf.addLbcfBackGroup();
                      }}
                    >
                      {t('添加')}
                    </Button>
                  </Bubble>
                </div>
              ) : (
                <noscript />
              )}
            </FormPanel.Item>
            <FormPanel.Footer>
              <Button
                className="mr10"
                type="primary"
                disabled={modifyMultiResourceWorkflow.operationState === OperationState.Performing}
                onClick={this._handleSubmit.bind(this)}
              >
                {failed ? t('重试') : t('配置')}
              </Button>
              <Button
                onClick={e => {
                  this._cancel();
                }}
              >
                {t('取消')}
              </Button>
              <TipInfo
                isShow={failed}
                type="error"
                style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
              >
                {getWorkflowError(modifyMultiResourceWorkflow)}
              </TipInfo>
            </FormPanel.Footer>
          </FormPanel>
        </ContentView.Body>
      </ContentView>
    );
  }
  private _cancel() {
    let { actions, subRoot, route, region, cluster, clusterVersion } = this.props,
      { modifyMultiResourceWorkflow, updateMultiResource } = subRoot,
      urlParams = router.resolve(route);
    let mode = urlParams['tab'] === 'createBG' ? 'create' : 'update';
    if (mode === 'create') {
      if (modifyMultiResourceWorkflow.operationState === OperationState.Done) {
        actions.workflow.modifyMultiResource.reset();
      }

      if (modifyMultiResourceWorkflow.operationState === OperationState.Started) {
        actions.workflow.modifyMultiResource.cancel();
      }
    } else {
      if (updateMultiResource.operationState === OperationState.Done) {
        actions.workflow.updateMultiResource.reset();
      }

      if (updateMultiResource.operationState === OperationState.Started) {
        actions.workflow.updateMultiResource.cancel();
      }
    }
    router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries);
  }
  /** 处理提交请求 */
  private _handleSubmit() {
    let { actions, subRoot, route, region, cluster, clusterVersion } = this.props,
      {
        resourceInfo,
        lbcfEdit,
        resourceOption: { ffResourceList }
      } = subRoot;
    actions.validate.lbcf.validateGameBGEdit();
    let { resourceIns } = route.queries,
      urlParams = router.resolve(route);
    let backGroupmode = urlParams['tab'] === 'createBG' ? 'create' : 'update';
    if (validateLbcfActions._validateGameBGEdit(lbcfEdit)) {
      let { namespace, lbcfBackGroupEditions } = lbcfEdit;
      let backGroupResourceInfo = resourceConfig(clusterVersion).lbcf_bg;
      let resources = [];

      if (backGroupmode === 'create') {
        lbcfBackGroupEditions.forEach(item => {
          let labelObject = {};
          let { labels, ports, name, backgroupType, serviceName, staticAddress, byName } = item;
          labels.forEach(label => {
            labelObject[label.key] = label.value;
          });
          let jsonData: LbcfBGJSONYaml = {
            kind: 'BackendGroup',
            apiVersion: 'lbcf.tkestack.io/v1beta1', //(resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
            metadata: {
              name: name,
              namespace: reduceNs(namespace)
            },
            spec: {
              lbName: resourceIns,
              pods:
                backgroupType === BackendType.Pods
                  ? {
                      port: { portNumber: +ports[0].portNumber, protocol: ports[0].protocol },
                      byLabel: labels.length
                        ? {
                            selector: labelObject
                          }
                        : undefined,
                      byName: byName.length ? byName.map(name => name.value) : undefined
                    }
                  : undefined,
              service:
                backgroupType === BackendType.Service
                  ? {
                      name: serviceName,
                      port: { portNumber: +ports[0].portNumber, protocol: ports[0].protocol },
                      nodeSelector: labels.length ? labelObject : undefined
                    }
                  : undefined,
              static: backgroupType === BackendType.Static ? staticAddress.map(address => address.value) : undefined
            }
          };
          jsonData = JSON.parse(JSON.stringify(jsonData));
          let resource: CreateResource = {
            id: uuid(),
            resourceInfo: backGroupResourceInfo,
            namespace: namespace,
            mode: backGroupmode,
            clusterId: route.queries['clusterId'],
            jsonData: JSON.stringify(jsonData)
          };
          resources.push(resource);
        });

        actions.workflow.modifyMultiResource.start(resources, region.selection.value);
        actions.workflow.modifyMultiResource.perform();
      } else {
        let finder = ffResourceList.list.data.records.find(item => item.metadata.name === resourceIns);
        let labelArray = {};
        finder &&
          finder.spec.backGroups.forEach(backGroup => {
            let label = deepClone(backGroup.labels);
            for (let key in label) {
              label[key] = null;
            }
            labelArray[backGroup.name] = label;
          });
        lbcfBackGroupEditions.forEach(item => {
          let labelObject = {};
          let { labels, ports, name, backgroupType, serviceName, staticAddress, byName } = item;
          labels.forEach(label => {
            labelObject[label.key] = label.value;
          });
          let jsonData = {};
          if (backgroupType === BackendType.Pods) {
            jsonData = {
              spec: {
                pods: {
                  port: { portNumber: +ports[0].portNumber, protocol: ports[0].protocol },
                  byLabel: {
                    selector: Object.assign({}, labelArray[item.name], labelObject)
                  },
                  byName: byName.length ? byName.map(name => name.value) : null
                }
              }
            };
          } else if (backgroupType === BackendType.Service) {
            jsonData = {
              spec: {
                service: {
                  name: serviceName,
                  port: { portNumber: +ports[0].portNumber, protocol: ports[0].protocol },
                  nodeSelector: Object.assign({}, labelArray[item.name], labelObject)
                }
              }
            };
          } else {
            jsonData = {
              spec: {
                static: staticAddress.map(name => name.value)
              }
            };
          }

          jsonData = JSON.parse(JSON.stringify(jsonData));
          let resource: CreateResource = {
            id: uuid(),
            resourceInfo: backGroupResourceInfo,
            resourceIns: name,
            namespace: namespace,
            mergeType: MergeType.Merge,
            mode: backGroupmode,
            clusterId: route.queries['clusterId'],
            jsonData: JSON.stringify(jsonData)
          };
          resources.push(resource);
        });

        actions.workflow.updateMultiResource.start(resources, region.selection.value);
        actions.workflow.updateMultiResource.perform();
      }
    }
  }
}
