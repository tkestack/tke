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

import { Button, ContentView, Form, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid, deepClone } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { validateLbcfActions } from '../../../../../modules/cluster/actions/validateLbcfActions';
import { TipInfo } from '../../../../common/components';
import { getWorkflowError } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { k8sVersionList, LbcfArgsConfig, LbcfConfig } from '../../../constants/Config';
import { CreateResource, LbcfLBJSONYaml } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { reduceNs } from '../../../../../../helpers';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditLbcfPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, route } = this.props;
    actions.lbcf.selectLbcfNamespace(route.queries['np']);
    actions.lbcf.driver.applyFilter({
      clusterId: route.queries['clusterId'],
      namespace: route.queries['np']
    });
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.lbcf.clearEdition();
  }

  render() {
    let { actions, subRoot, route, namespaceList } = this.props,
      urlParams = router.resolve(route),
      { lbcfEdit, modifyResourceFlow } = subRoot,
      { name, v_name, namespace, v_namespace, config, args, driver, v_args, v_config, v_driver } = lbcfEdit;

    /** 渲染namespace列表 */
    let namespaceOptions = namespaceList.data.recordCount
      ? namespaceList.data.records.map((item, index) => {
          return {
            value: item.name,
            text: item.name
          };
        })
      : [];

    let failed = modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow);

    return (
      <ContentView>
        <ContentView.Body>
          <FormPanel>
            <FormPanel.Item
              label={t('名称')}
              errorTipsStyle="Icon"
              validator={v_name}
              input={{
                value: name,
                placeholder: t('请输入名称'),
                onChange: actions.lbcf.inputLbcfName,
                onBlur: e => actions.validate.lbcf.validateLbcfName()
              }}
              message={t('最长63个字符，只能包含小写字母、数字及分隔符("-")，且必须以小写字母开头，数字或小写字母结尾')}
            />
            <FormPanel.Item
              label={t('命名空间')}
              errorTipsStyle="Icon"
              validator={v_namespace}
              select={{
                value: namespace,
                options: namespaceOptions,
                onChange: value => {
                  actions.lbcf.selectLbcfNamespace(value);
                  actions.namespace.selectNamespace(value);
                }
              }}
            />
            <FormPanel.Item
              label={t('负载均衡类型')}
              errorTipsStyle="Icon"
              validator={v_driver}
              select={{
                model: driver,
                action: actions.lbcf.driver
              }}
            />
            <FormPanel.Item label={t('负载均衡配置')} validator={v_config} errorTipsStyle={'Message'}>
              <FormPanel isNeedCard={false} fixed>
                <React.Fragment>
                  {config.map((kv, index) => {
                    return (
                      <div key={index} style={{ marginBottom: 10 }}>
                        <FormPanel.Input
                          value={config[index].key}
                          onChange={value => {
                            let newConfig = deepClone(config);
                            newConfig[index].key = value;
                            actions.lbcf.selectConfig(newConfig);
                          }}
                        />
                        <FormPanel.InlineText style={{ marginLeft: 5, marginRight: 5 }}>=</FormPanel.InlineText>
                        <FormPanel.Input
                          value={config[index].value}
                          onChange={value => {
                            let newConfig = deepClone(config);
                            newConfig[index].value = value;
                            actions.lbcf.selectConfig(newConfig);
                          }}
                        />
                        <Button
                          icon="close"
                          onClick={() => {
                            let newConfig = deepClone(config);
                            newConfig.splice(index, 1);
                            actions.lbcf.selectConfig(newConfig);
                          }}
                        />
                      </div>
                    );
                  })}

                  <div>
                    <Button
                      type="link"
                      onClick={() => {
                        let newConfig = deepClone(config);
                        newConfig.push({ key: '', value: '' });
                        actions.lbcf.selectConfig(newConfig);
                      }}
                    >
                      {t('新增配置')}
                    </Button>
                  </div>
                </React.Fragment>
              </FormPanel>
            </FormPanel.Item>
            <FormPanel.Item label={t('负载均衡属性')} validator={v_args} errorTipsStyle={'Message'}>
              <FormPanel isNeedCard={false} fixed>
                <React.Fragment>
                  {args.map((kv, index) => {
                    return (
                      <div key={index} style={{ marginBottom: 10 }}>
                        <FormPanel.Input
                          value={args[index].key}
                          onChange={value => {
                            let newArgs = deepClone(args);
                            newArgs[index].key = value;
                            actions.lbcf.selectArgs(newArgs);
                          }}
                        />
                        <FormPanel.InlineText style={{ marginLeft: 5, marginRight: 5 }}>=</FormPanel.InlineText>
                        <FormPanel.Input
                          value={args[index].value}
                          onChange={value => {
                            let newArgs = deepClone(args);
                            newArgs[index].value = value;
                            actions.lbcf.selectArgs(newArgs);
                          }}
                        />
                        <Button
                          icon="close"
                          onClick={() => {
                            let newArgs = deepClone(args);
                            newArgs.splice(index, 1);
                            actions.lbcf.selectArgs(newArgs);
                          }}
                        />
                      </div>
                    );
                  })}

                  <div>
                    <Button
                      type="link"
                      onClick={() => {
                        let newArgs = deepClone(args);
                        newArgs.push({ key: '', value: '' });
                        actions.lbcf.selectArgs(newArgs);
                      }}
                    >
                      {t('新增属性')}
                    </Button>
                  </div>
                </React.Fragment>
              </FormPanel>
            </FormPanel.Item>

            <FormPanel.Footer>
              <Button
                className="mr10"
                type="primary"
                disabled={modifyResourceFlow.operationState === OperationState.Performing}
                onClick={this._handleSubmit.bind(this)}
              >
                {failed ? t('重试') : t('创建负载均衡')}
              </Button>
              <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                {t('取消')}
              </Button>
              <TipInfo
                isShow={failed}
                type="error"
                style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
              >
                {getWorkflowError(modifyResourceFlow)}
              </TipInfo>
            </FormPanel.Footer>
          </FormPanel>
        </ContentView.Body>
      </ContentView>
    );
  }

  /** 处理提交请求 */
  private _handleSubmit() {
    let { actions, subRoot, route, region, cluster } = this.props,
      { resourceInfo, mode, lbcfEdit } = subRoot,
      { config, args, driver } = lbcfEdit;

    actions.validate.lbcf.validateLbcfEdit();
    if (validateLbcfActions._validateLbcfEdit(lbcfEdit)) {
      let lbcfConfig = {};
      let lbcfAttribute = {};
      config.forEach(c => {
        if (c.key === 'loadBalancerID' && c.value === '') return;
        lbcfConfig[c.key] = c.value;
      });
      args.forEach(a => {
        lbcfAttribute[a.key] = a.value;
      });
      let {
        name,
        namespace
        // clbId, createLbWay, vpcId
      } = lbcfEdit;
      let jsonData: LbcfLBJSONYaml = {
        kind: 'LoadBalancer',
        apiVersion: 'lbcf.tkestack.io/v1beta1', //(resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
        metadata: {
          name: name,
          namespace: reduceNs(namespace)
        },
        spec: {
          lbDriver: driver.selection ? driver.selection.metadata.name : '',
          lbSpec: lbcfConfig,
          attributes: lbcfAttribute
        }
      };
      jsonData = JSON.parse(JSON.stringify(jsonData));
      let resource: CreateResource = {
        id: uuid(),
        resourceInfo,
        namespace: namespace || 'default',
        mode,
        clusterId: route.queries['clusterId'],
        jsonData: JSON.stringify(jsonData)
      };
      actions.workflow.modifyResource.start([resource], region.selection && region.selection.value);
      actions.workflow.modifyResource.perform();
    }
  }
}
