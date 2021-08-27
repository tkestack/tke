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

import { Bubble, Button, Table, TableColumn } from '@tea/component';
import { stylize } from '@tea/component/table/addons/stylize';
import { bindActionCreators, FetchState, isSuccessWorkflow, OperationState, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import {
  ButtonBar,
  FormItem,
  InputField,
  ResourceSelectorGeneric,
  ResourceSelectorInfoRow,
  ResourceSelectorProps,
  TipInfo
} from '../../../../common/components';
import { FormLayout, MainBodyLayout } from '../../../../common/layouts';
import { getWorkflowError } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { validateSecretActions } from '../../../actions/validateSecretActinos';
import { SecretTypeList } from '../../../constants/Config';
import { CreateResource, Namespace, SecretData, SecretEditJSONYaml } from '../../../models';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { reduceNs } from '../../../../../../helpers';

const secretTypeTip = {
  Opaque: t('适用于保存秘钥证书和配置文件，Value将以base64格式编码'),
  'kubernetes.io/dockercfg': t('用于保存私有Docker Registry的认证信息')
};

const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

interface EditSecretPanelState {
  /** 当前选择的命名空间的类别 */
  namespaceMode?: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditSecretPanel extends React.Component<RootProps, EditSecretPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      namespaceMode: 'partitial'
    };
  }

  componentDidMount() {
    let { actions, route } = this.props;
    // 进行ns的拉取
    actions.editSecret.ns.applyFilter({ clusterId: route.queries['clusterId'], regionId: route.queries['rid'] });
  }

  componentWillUnmount() {
    let { actions } = this.props;
    // 清除 secret的编辑信息
    actions.editSecret.clearSecretEdit();
  }

  render() {
    let { actions, subRoot, route } = this.props,
      urlParams = router.resolve(route),
      { modifyResourceFlow, secretEdit } = subRoot,
      { name, v_name, secretType, domain, v_domain, password, v_password, username, v_username, nsType } = secretEdit;

    // 判断当前的secret的类型
    let isOpaque = secretType === 'Opaque',
      isDockercfg = secretType === 'kubernetes.io/dockercfg';

    // 渲染secret的类型 选择项
    let selectSecretType = SecretTypeList.find(item => item.value === secretType);

    // 当前的生效范围
    let isSelectAllNs = nsType === 'all';

    let failed = modifyResourceFlow.operationState === OperationState.Done && !isSuccessWorkflow(modifyResourceFlow);

    return (
      <MainBodyLayout>
        <FormLayout>
          <div className="param-box server-update add">
            <ul className="form-list jiqun fixed-layout" style={{ paddingBottom: '50px' }}>
              <FormItem label={t('名称')}>
                <InputField
                  type="text"
                  placeholder={t('请输入名称')}
                  tipMode="popup"
                  validator={v_name}
                  value={name}
                  onChange={actions.editSecret.inputSecretName}
                  onBlur={actions.validate.secret.validateName}
                />
              </FormItem>
              <FormItem label={t('Secret类型')}>
                <div className="form-unit">
                  <ButtonBar
                    size="m"
                    isNeedPureText={false}
                    list={SecretTypeList}
                    selected={selectSecretType}
                    onSelect={item => {
                      actions.editSecret.selectSecretType(item.value + '');
                    }}
                  />
                  <p className="text-label form-input-help">{secretTypeTip[secretType]}</p>
                </div>
              </FormItem>

              <FormItem label={t('生效范围')}>
                <div className="up-date" style={{ width: '720px' }}>
                  <div className="as-sel-box">
                    <label className="form-ctrl-label" style={{ marginBottom: '0' }}>
                      <input
                        type="radio"
                        className="tc-15-radio"
                        checked={isSelectAllNs}
                        onChange={e => {
                          actions.editSecret.selectNsType('all');
                        }}
                      />
                      <span style={{ verticalAlign: '0' }}>
                        {t('存量所有命名空间（不包括kube-system、kube-public和后续增量命名空间）')}
                      </span>
                    </label>
                  </div>
                </div>
                <div className="up-date" style={{ width: '720px' }}>
                  <div className="as-sel-box">
                    <label className="form-ctrl-label" style={{ marginBottom: '0' }}>
                      <input
                        type="radio"
                        className="tc-15-radio"
                        checked={!isSelectAllNs}
                        onChange={e => {
                          actions.editSecret.selectNsType('specific');
                        }}
                      />
                      <span style={{ verticalAlign: '0' }}>{t('指定命名空间')}</span>
                    </label>
                    {!isSelectAllNs && this._renderSpecificNs()}
                  </div>
                </div>
              </FormItem>

              {isOpaque && <EditSecretDataForOpaque />}

              <FormItem label={t('仓库域名')} isShow={isDockercfg}>
                <Bubble
                  placement="left"
                  content={
                    v_domain.status === 2 ? (
                      <span className="form-input-help" style={{ marginLeft: '0' }}>
                        {v_domain.message}
                      </span>
                    ) : null
                  }
                >
                  <div className={classnames('form-unit', { 'is-error': v_domain.status === 2 })}>
                    <input
                      type="text"
                      className="tc-15-input-text m"
                      placeholder={t('请输入域名或IP')}
                      value={domain}
                      onChange={e => actions.editSecret.inputThirdHubDomain(e.target.value)}
                      onBlur={actions.validate.secret.validateThirdHubDomain}
                    />
                  </div>
                </Bubble>
              </FormItem>

              <FormItem label={t('用户名')} isShow={isDockercfg}>
                <Bubble
                  placement="left"
                  content={
                    v_username.status === 2 ? (
                      <span className="form-input-help" style={{ marginLeft: '0' }}>
                        {v_username.message}
                      </span>
                    ) : null
                  }
                >
                  <div className={classnames('form-unit', { 'is-error': v_username.status === 2 })}>
                    <input
                      type="text"
                      className="tc-15-input-text m"
                      placeholder={t('请输入第三方仓库的用户名')}
                      value={username}
                      onChange={e => actions.editSecret.inputThirdHubUserName(e.target.value)}
                      onBlur={actions.validate.secret.validateThirdHubUsername}
                    />
                  </div>
                </Bubble>
              </FormItem>

              <FormItem label={t('密码')} isShow={isDockercfg}>
                <Bubble
                  placement="left"
                  content={
                    v_password.status === 2 ? (
                      <span className="form-input-help" style={{ marginLeft: '0' }}>
                        {v_password.message}
                      </span>
                    ) : null
                  }
                >
                  <div className={classnames('form-unit', { 'is-error': v_password.status === 2 })}>
                    <input
                      type="password"
                      className="tc-15-input-text m"
                      placeholder={t('请输入第三方仓库的登陆密码')}
                      value={password}
                      onChange={e => actions.editSecret.inputThirdHubPassword(e.target.value)}
                      onBlur={actions.validate.secret.validateThirdHubPassword}
                    />
                  </div>
                </Bubble>
              </FormItem>

              <li className="pure-text-row fixed">
                <Button
                  className="mr10"
                  type="primary"
                  disabled={modifyResourceFlow.operationState === OperationState.Performing}
                  onClick={this._handleSubmit.bind(this)}
                >
                  {failed ? t('重试') : t('创建Secret')}
                </Button>
                <Button onClick={e => router.navigate(Object.assign({}, urlParams, { mode: 'list' }), route.queries)}>
                  {t('取消')}
                </Button>
                <TipInfo
                  isShow={failed}
                  className="error"
                  style={{ display: 'inline-block', marginLeft: '20px', marginBottom: '0px', maxWidth: '750px' }}
                >
                  {getWorkflowError(modifyResourceFlow)}
                </TipInfo>
              </li>
            </ul>
          </div>
        </FormLayout>
      </MainBodyLayout>
    );
  }

  /** 选择部分的ns的相关展示 */
  private _renderSpecificNs() {
    let { actions, subRoot } = this.props,
      { nsListSelection, nsList, nsQuery } = subRoot.secretEdit;

    // 是否需要展示ns的loading态
    let isShowLoadingNs: boolean = nsList.fetched !== true || nsList.fetchState === FetchState.Fetching;

    // 表示ResourceSelector 里要显示和选择的数据类型
    const ResourceSelector = ResourceSelectorGeneric as new () => ResourceSelectorGeneric<Namespace>;

    // 参数的配置
    const selectorProps: ResourceSelectorProps<Namespace> = {
      /** 要供选择的数据 */
      list: nsList.data.records,

      /** 已选中的数据 */
      selection: nsListSelection,

      /** 用户选择发生改变后，应该更新选中的数据状态 */
      onSelectionChanged: selected => actions.editSecret.selectNsList(selected),

      /** 选择器标题 */
      selectorTitle: t('当前集群有以下可用命名空间'),

      /** 如何渲染具体一项的名字 */
      itemNameRender: namespace => {
        return (
          <div>
            <span title={namespace.displayName}>{namespace.displayName}</span>
          </div>
        );
      },

      search: {
        // keyword: nsQuery.keyword,
        // onChange: actions.editSecret.ns.changeKeyword,
        onSearch: actions.editSecret.ns.performSearch,
        placeholder: t('请输入命名空间')
      }
    };

    return (
      <ResourceSelector {...selectorProps} style={{ overflow: 'auto' }}>
        {isShowLoadingNs && <ResourceSelectorInfoRow>{loadingElement}</ResourceSelectorInfoRow>}
        {nsList.data.recordCount <= 0 && <ResourceSelectorInfoRow>{t('暂无可用命名空间')}</ResourceSelectorInfoRow>}
      </ResourceSelector>
    );
  }

  /** 处理请求提交 */
  private _handleSubmit() {
    let { actions, subRoot, route, region, namespaceList } = this.props,
      { resourceInfo, secretEdit } = subRoot;

    actions.validate.secret.validateSecretEdit();

    if (validateSecretActions._validateSecretEdit(secretEdit)) {
      let { name, secretType, data, nsType, nsListSelection, username, password, domain } = secretEdit;

      let isOpaque = secretType === 'Opaque',
        isDockercfg = secretType === 'kubernetes.io/dockercfg';

      // 提交的数据字段
      let secretData = {};

      // 如果是opaque的类型的话，则处理key、value的值
      if (isOpaque) {
        data.forEach(item => {
          secretData[item.keyName] = window.btoa(item.value);
        });
      }

      // 如果是 dockercfg的类型，则data当中是 .dockercfg字段，并且内容是 dockercfg的内容，需要username、domain、password和auth字段
      if (isDockercfg) {
        let finalData = {
          [domain]: {
            username,
            password,
            auth: window.btoa(`${username}:${password}`)
          }
        };
        secretData['.dockercfg'] = window.btoa(JSON.stringify(finalData));
      }

      let finalJSON = '';
      let finalNsList =
        nsType === 'all'
          ? namespaceList.data.records.filter(item => item.name !== 'kube-system' && item.name !== 'kube-public')
          : nsListSelection;
      finalNsList.forEach(item => {
        let jsonData: SecretEditJSONYaml = {
          kind: resourceInfo.headTitle,
          apiVersion: (resourceInfo.group ? resourceInfo.group + '/' : '') + resourceInfo.version,
          metadata: {
            name: name,
            namespace: item.name,
            labels: {
              'qcloud-app': name
            }
          },
          type: secretType,
          data: secretData
        };

        finalJSON += JSON.stringify(jsonData);
      });

      let resource: CreateResource = {
        id: uuid(),
        clusterId: route.queries['clusterId'],
        jsonData: finalJSON
      };

      actions.workflow.applyResource.start([resource], region.selection.value);
      actions.workflow.applyResource.perform();
    }
  }
}

/** secret 内容的展示 key value的形式 */

@connect(state => state, mapDispatchToProps)
class EditSecretDataForOpaque extends React.Component<RootProps, {}> {
  render() {
    let { actions } = this.props;

    return (
      <FormItem className="vm" label={t('内容')}>
        {this._renderSecretData()}
        <a href="javascript:;" className="more-links-btn" onClick={actions.editSecret.addSecretData}>
          {t('添加变量')}
        </a>
      </FormItem>
    );
  }

  /** 展示变量的设置 */
  private _renderSecretData() {
    let { subRoot, actions } = this.props,
      { data } = subRoot.secretEdit;

    let columns: TableColumn<SecretData>[] = [
      {
        key: 'keyName',
        header: t('变量名'),
        width: '40%',
        render: x => (
          <div className={x.v_keyName.status === 2 ? 'is-error' : ''}>
            <Bubble placement="bottom" content={x.v_keyName.status === 2 ? x.v_keyName.message : null}>
              <input
                type="text"
                className="tc-15-input-text m"
                value={x.keyName + ''}
                onChange={e => actions.editSecret.updateSecretData({ keyName: e.target.value }, x.id + '')}
                onBlur={e => actions.validate.secret.validateKeyName(e.target.value, x.id + '')}
              />
            </Bubble>
          </div>
        )
      },
      {
        key: 'equal',
        header: ' ',
        width: '10%',
        render: x => (
          <div>
            <p>=</p>
          </div>
        )
      },
      {
        key: 'value',
        header: t('变量值'),
        width: '40%',
        render: x => (
          <div className={x.v_value.status === 2 ? 'is-error' : ''}>
            <Bubble placement="bottom" content={x.v_value.status === 2 ? x.v_value.message : null}>
              <textarea
                className="tc-15-input-text m"
                style={{ maxWidth: '260px', overflowY: 'visible' }}
                value={x.value + ''}
                onChange={e => actions.editSecret.updateSecretData({ value: e.target.value }, x.id + '')}
                onBlur={e => actions.validate.secret.validateKeyValue(e.target.value, x.id + '')}
              />
            </Bubble>
          </div>
        )
      },
      {
        key: '',
        header: '',
        width: '10%',
        render: x => (
          <div>
            {data.length > 1 ? (
              <a
                href="javascript:;"
                onClick={() => {
                  actions.editSecret.deleteSecretData(x.id + '');
                }}
              >
                <i className="icon-cancel-icon" />
              </a>
            ) : (
              <Bubble placement="bottom" content={t('不可删除，至少设置一项')}>
                <a href="javascript:;" className="disabled">
                  <i className="icon-cancel-icon" />
                </a>
              </Bubble>
            )}
          </div>
        )
      }
    ];

    return (
      <Table
        columns={columns}
        records={data}
        addons={[
          stylize({
            style: { overflow: 'visible', maxWidth: '550px' }
          })
        ]}
      />
    );
  }
}
