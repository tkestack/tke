import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { allActions } from '../../../actions';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { Bubble } from '@tea/component';
import { connect } from 'react-redux';
import { isEmpty } from '../../../../common/utils';
import { FormItem, LinkButton } from '../../../../common/components';
import * as classnames from 'classnames';
import { Resource } from '../../../models';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

/** 创建workload，环境变量，valueFrom的列表 */
export const valueFromList = [
  {
    value: 'configMap',
    label: 'ConfigMap'
  },
  {
    value: 'secret',
    label: 'Secret'
  }
];

const valueFromSelectListStyle: React.CSSProperties = {
  marginRight: '6px',
  minWidth: '120px',
  maxWidth: '120px'
};

interface ContainerEnvItemProps extends RootProps {
  cKey: string;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class EditResourceContainerEnvItem extends React.Component<ContainerEnvItemProps, {}> {
  render() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    let container = containers.find(c => c.id === cKey);
    let envs = container ? container.envs : [],
      valueFrom = container ? container.valueFrom : [];
    let canAdd = isEmpty(envs.filter(x => !x.envName)),
      canAddValueFrom =
        valueFrom.length === 0 || valueFrom.filter(x => x.key !== '' && x.name !== '' && x.aliasName !== '').length > 0
          ? true
          : false;

    return (
      <FormItem label={t('环境变量')} tips={t('设置容器中的变量')}>
        <div className="form-unit is-success">
          {this._renderEnvList()}

          {envs.length > 0 && valueFrom.length > 0 && <hr className="hr-mod" />}

          {this._renderValueFromList()}

          <LinkButton
            className="tea-mr-1n"
            disabled={!canAdd}
            errorTip={t('请先完成待编辑项')}
            onClick={() => {
              actions.editWorkload.addEnv(cKey);
            }}
          >
            {t('新增变量')}
          </LinkButton>

          <LinkButton
            disabled={!canAddValueFrom}
            errorTip={t('请先完成待编辑项')}
            onClick={() => {
              actions.editWorkload.addValueFrom(cKey);
            }}
          >
            {t('引用ConfigMap/Secret')}
          </LinkButton>

          <p className="text-label">{t('变量名只能包含大小写字母、数字及下划线，并且不能以数字开头')}</p>
        </div>
      </FormItem>
    );
  }

  /** 展示环境变量的选项 */
  private _renderEnvList() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers } = workloadEdit;

    let container = containers.find(c => c.id === cKey);
    let envs = container ? container.envs : [];
    return envs.map((env, index) => {
      return (
        <div className="code-list" key={index}>
          <div className={classnames({ 'is-error': env.v_envName.status === 2 })} style={{ display: 'inline-block' }}>
            <Bubble placement="bottom" content={env.v_envName.status === 2 ? <p>{env.v_envName.message}</p> : null}>
              <input
                type="text"
                placeholder={t('变量名')}
                className="tc-15-input-text m"
                style={{ width: '128px' }}
                value={env.envName}
                onChange={e => actions.editWorkload.updateEnv({ envName: e.target.value }, cKey, env.id + '')}
                onBlur={e => actions.validate.workload.validateAllEnvName(container)}
                maxLength={63}
              />
            </Bubble>
            <span className="inline-help-text">=</span>
            <textarea
              placeholder={t('变量值')}
              className="tc-15-input-text m"
              style={{
                maxWidth: '260px',
                minHeight: '30px',
                minWidth: '128px',
                overflowY: 'visible',
                marginLeft: '5px'
              }}
              value={env.envValue}
              maxLength={63}
              onChange={e => actions.editWorkload.updateEnv({ envValue: e.target.value }, cKey, env.id + '')}
            />
            <span className="inline-help-text">
              <LinkButton onClick={() => actions.editWorkload.deleteEnv(cKey, env.id + '')}>
                <i className="icon-cancel-icon" />
              </LinkButton>
            </span>
          </div>
        </div>
      );
    });
  }

  /** 展示valueFrom的选项 */
  private _renderValueFromList() {
    let { actions, cKey, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { containers, configEdit } = workloadEdit,
      { configList, secretList } = configEdit;

    let container = containers.find(c => c.id === cKey);
    let valueFrom = container ? container.valueFrom : [];

    /** 渲染valueFrom的类型 */
    let configTypeOptions = valueFromList.map(type => {
      return (
        <option key={uuid()} value={type.value}>
          {type.label}
        </option>
      );
    });
    configTypeOptions.unshift(
      <option key={uuid()} value="">
        {t('请选择类型')}
      </option>
    );

    return valueFrom.map((item, index) => {
      let vId = item.id + '';

      // let aliasInput = this._renderAliasInput(item, vId);
      return (
        <div className="code-list" key={index} style={{ marginBottom: '5px' }}>
          <select
            className="tc-15-select m"
            style={valueFromSelectListStyle}
            value={item.type}
            onChange={e => {
              actions.editWorkload.updateValueFrom({ type: e.target.value }, cKey, vId);
              this._handleResourceName('', cKey, vId);
            }}
          >
            {configTypeOptions}
          </select>

          <select
            className="tc-15-select m"
            disabled={item.type === ''}
            style={valueFromSelectListStyle}
            value={item.name}
            onChange={e => {
              this._handleResourceName(e.target.value, cKey, vId);
            }}
          >
            {this._reduceResourceListOptions(
              item.type === 'configMap' ? configList.data.records : secretList.data.records
            )}
          </select>

          <div className={classnames({ 'is-error': item.v_key.status === 2 })} style={{ display: 'inline-block' }}>
            <Bubble placement="bottom" content={item.v_key.status === 2 ? item.v_key.message : null}>
              <select
                className="tc-15-select m"
                disabled={item.type === ''}
                style={valueFromSelectListStyle}
                value={item.key}
                onChange={e => {
                  actions.editWorkload.updateValueFrom({ key: e.target.value }, cKey, vId);
                  actions.validate.workload.validateValueFromKey(e.target.value, cKey, vId);
                }}
              >
                {this._rendeceResourceKeyOptions(item.type, item.name)}
              </select>
            </Bubble>
          </div>

          <Trans>
            <span className="inline-help-text text-label" style={{ margin: '0 5px 0 0' }}>
              以
            </span>
            <div
              className={classnames({ 'is-error': item.v_aliasName.status === 2 })}
              style={{ display: 'inline-block' }}
            >
              <Bubble placement="bottom" content={item.v_aliasName.status === 2 ? item.v_aliasName.message : null}>
                <input
                  type="text"
                  className="tc-15-input-text m"
                  style={{ maxWidth: '120px' }}
                  placeholder={t('请输入别名')}
                  value={item.aliasName}
                  onChange={e => actions.editWorkload.updateValueFrom({ aliasName: e.target.value }, cKey, vId)}
                  onBlur={e => actions.validate.workload.validateValueFromAlias(e.target.value, cKey, vId)}
                />
              </Bubble>
            </div>
            <span className="inline-help-text text-label">为别名</span>
          </Trans>

          <span className="inline-help-text">
            <LinkButton onClick={() => actions.editWorkload.deleteValueFrom(cKey, vId)}>
              <i className="icon-cancel-icon" />
            </LinkButton>
          </span>
        </div>
      );
    });
  }

  // private _renderAliasInput(item, vId) {
  //   let { actions, cKey, subRoot } = this.props;
  //   return (
  //     <div className={classnames({ 'is-error': item.v_aliasName.status === 2 })} style={{ display: 'inline-block' }}>
  //       <Bubble placement="bottom" content={item.v_aliasName.status === 2 ? item.v_aliasName.message : null}>
  //         <input
  //           type="text"
  //           className="tc-15-input-text m"
  //           style={{ maxWidth: '120px' }}
  //           placeholder={t('请输入别名')}
  //           value={item.aliasName}
  //           onChange={e => actions.editWorkload.updateValueFrom({ aliasName: e.target.value }, cKey, vId)}
  //           onBlur={e => actions.validate.workload.validateValueFromAlias(e.target.value, cKey, vId)}
  //         />
  //       </Bubble>
  //     </div>
  //   );
  // }

  /** 选择具体的资源的操作 */
  private _handleResourceName(resourceName: string, cKey: string, vId: string) {
    let { actions } = this.props;
    actions.editWorkload.updateValueFrom({ name: resourceName, key: '' }, cKey, vId);
  }

  /** 渲染resourceList的options选项 */
  private _reduceResourceListOptions(list: Resource[]) {
    let configListOptions = list.map(item => {
      return (
        <option key={uuid()} value={item.metadata.name}>
          {item.metadata.name}
        </option>
      );
    });
    configListOptions.unshift(
      <option key={uuid()} value="">
        {list.length ? t('请选择资源') : t('列表为空')}
      </option>
    );

    return configListOptions;
  }

  /** 渲染resourceKey的options选项 */
  private _rendeceResourceKeyOptions(resourceType: string, name: string) {
    let { configList, secretList } = this.props.subRoot.workloadEdit.configEdit;

    let list = resourceType === 'configMap' ? configList.data.records : secretList.data.records,
      finder = list.find(item => item.metadata.name === name);
    let dataKeys = finder && finder.data ? Object.keys(finder.data) : [];
    let configKeyOptions = dataKeys.map(item => {
      return (
        <option key={uuid()} value={item}>
          {item}
        </option>
      );
    });
    configKeyOptions.unshift(
      <option key={uuid()} value="">
        {dataKeys.length ? t('选择Key') : t('列表为空')}
      </option>
    );

    return configKeyOptions;
  }
}
