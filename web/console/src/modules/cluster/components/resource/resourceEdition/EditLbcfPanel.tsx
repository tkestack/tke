import * as React from 'react';
import { connect } from 'react-redux';

import { Button, ContentView, Form, Text } from '@tea/component';
import { FormPanel } from '@tencent/ff-component';
import { bindActionCreators, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
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
  }

  componentWillUnmount() {
    let { actions } = this.props;
    actions.lbcf.clearEdition();
  }

  render() {
    let { actions, subRoot, route, namespaceList, cluster } = this.props,
      urlParams = router.resolve(route),
      { lbcfEdit, modifyResourceFlow } = subRoot,
      { name, v_name, namespace, v_namespace, config, args } = lbcfEdit;

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
            <FormPanel.Item label={t('负载均衡配置')}>
              <FormPanel isNeedCard={false} fixed>
                <FormPanel.Item
                  keyvalue={{
                    options: LbcfConfig,
                    value: config,
                    onChange: (values, option) => {
                      actions.lbcf.selectConfig(values);
                      if (option.value === 'loadBalancerType') {
                        let loadBalancerType = values.find(v => v.key === 'loadBalancerType');
                        if (loadBalancerType.value === 'INTERNAL') {
                          if (!values.find(v => v.key === 'subnetID')) {
                            values.push({
                              key: 'subnetID',
                              value: ''
                            });
                            actions.lbcf.selectConfig(values.slice(0));
                          }
                        }
                      }
                      if (option.value === 'listenerProtocol') {
                        let listenerProtocol = values.find(v => v.key === 'listenerProtocol');
                        if (listenerProtocol) {
                          switch (listenerProtocol.value) {
                            case 'HTTP':
                              //
                              if (!values.find(v => v.key === 'domain')) {
                                values.push({
                                  key: 'domain',
                                  value: ''
                                });
                              }
                              if (!values.find(v => v.key === 'url')) {
                                values.push({
                                  key: 'url',
                                  value: ''
                                });
                              }
                              actions.lbcf.selectConfig(values.slice(0));
                              break;
                            case 'HTTPS':
                              if (!values.find(v => v.key === 'domain')) {
                                values.push({
                                  key: 'domain',
                                  value: ''
                                });
                              }
                              if (!values.find(v => v.key === 'url')) {
                                values.push({
                                  key: 'url',
                                  value: ''
                                });
                              }
                              actions.lbcf.selectConfig(values.slice(0));
                              if (!args.find(a => a.key === 'listenerCertID')) {
                                args.push({
                                  key: 'listenerCertID',
                                  value: ''
                                });
                                actions.lbcf.selectArgs(args.slice(0));
                              }

                              break;
                          }
                        }
                      }
                    }
                  }}
                />
              </FormPanel>
            </FormPanel.Item>
            <FormPanel.Item label={t('负载均衡属性')}>
              <FormPanel isNeedCard={false} fixed>
                <FormPanel.Item
                  keyvalue={{
                    options: LbcfArgsConfig,
                    value: args,
                    onChange: values => actions.lbcf.selectArgs(values)
                  }}
                />
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
      { config, args } = lbcfEdit;

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
          lbDriver: 'lbcf-clb-driver',
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
