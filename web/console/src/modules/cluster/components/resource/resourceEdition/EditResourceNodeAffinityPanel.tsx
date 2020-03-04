import * as React from 'react';
import { connect } from 'react-redux';

import { Bubble, ExternalLink, Radio } from '@tea/component';
import { bindActionCreators, deepClone, FetchState } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { FormItem, InputField, LinkButton } from '../../../../common/components';
import { isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { affinityRuleOperator, affinityType } from '../../../constants/Config';
import { MatchExpressions } from '../../../models/WorkloadEdit';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class EditResourceNodeAffinityPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, subRoot, route } = this.props;
    let { computer } = subRoot.computerState;
    let { clusterId, rid } = route.queries;
    computer.list.fetched !== true && actions.computer.applyFilter({ clusterId, regionId: +rid });
  }

  _renderComputerList() {
    let { route, subRoot } = this.props,
      { computer } = subRoot.computerState,
      { nodeSelection, v_nodeSelection } = subRoot.workloadEdit;
    let content: JSX.Element[] = [];

    if (computer.list.fetched !== true || computer.list.fetchState === FetchState.Fetching) {
      // do something
    }

    computer.list.data.records.forEach((computer, index) => {
      let item: JSX.Element;
      try {
        let readyCondition = computer.status.conditions.filter(item => item.type === 'Ready')[0];
        let isComputerRunning = readyCondition['status'] === 'True';

        item = (
          <label key={index + 'label'} className="form-ctrl-label" style={{ display: 'block', margin: 10 }}>
            <Bubble placement="top" content={!isComputerRunning ? t('节点状态不正常') : null}>
              <input
                disabled={!isComputerRunning}
                type="checkbox"
                className="tc-15-checkbox"
                checked={nodeSelection.findIndex(node => node.metadata.name === computer.metadata.name) !== -1}
                style={{ verticalAlign: 'middle' }}
                onChange={e => this._handleNodeSelection(e.target.value)}
                value={computer.metadata.name}
              />
              <span>{`${computer.metadata.name}(${computer.metadata.role})`}</span>
            </Bubble>
          </label>
        );
      } catch (error) {}
      item && content.push(item);
    });
    if (computer.list.data.recordCount === 0) {
      content.push(
        <div style={{ fontSize: '11px', marginBottom: '-5px', marginTop: '2px' }}>
          <strong>{t('该集群无可用节点')}</strong>
          {/* { <p className='text-danger'></p> */}
        </div>
      );
    }
    v_nodeSelection.status === 2 &&
      content.push(
        <p className="text-danger" style={{ fontSize: '11px' }}>
          {v_nodeSelection.message}
        </p>
      );
    return content;
  }

  _handleNodeSelection(value) {
    let { subRoot, actions } = this.props,
      { computer } = subRoot.computerState,
      nodeSelection = deepClone(subRoot.workloadEdit.nodeSelection);
    let index = nodeSelection.findIndex(node => node.metadata.name === value);
    if (index !== -1) {
      nodeSelection.splice(index, 1);
    } else {
      let item = computer.list.data.records.find(computer => computer.metadata.name === value);
      item && nodeSelection.push(item);
    }
    actions.editWorkload.selectNodeSelector(nodeSelection);
  }

  _renderAffinityRuleList(type: string) {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { nodeAffinityType, nodeAffinityRule, nodeSelection } = workloadEdit;

    // 获得当前的pod节点调度的操作符
    let operatorTypeList = affinityRuleOperator.map((item, index) => (
      <option value={item.value} key={index}>
        {item.value}
      </option>
    ));
    //获取当前渲染的条件组
    let items: MatchExpressions[];
    if (type === 'preferred') {
      items = nodeAffinityRule.preferredExecution[0].preference.matchExpressions;
    } else if (type === 'required') {
      items = nodeAffinityRule.requiredExecution[0].matchExpressions;
    }
    return items.map((item, index) => {
      let defaultAffinityOperator = item.operator !== '' ? item.operator : 'In',
        finder = affinityRuleOperator.find(item => item.value === defaultAffinityOperator),
        operatorTip = finder ? finder.tip : '';
      let isNeedValues = item.operator !== 'Exists' && item.operator !== 'DoesNotExist';
      return (
        <div className="code-list mb-10" key={index}>
          <InputField
            type="input"
            validator={item.v_key}
            value={item.key}
            placeholder="Label Key"
            style={{ width: '100px', marginRight: '10px' }}
            tipMode="popup"
            onChange={value => {
              actions.editWorkload.updateAffinityRule(type, item.id + '', { key: value });
            }}
            onBlur={e => {
              actions.validate.workload.validateNodeAffinityRuleKey(type, item.id + '');
            }}
          />
          <div className={'mr10'} style={{ display: 'inline-block', verticalAlign: 'top' }}>
            <Bubble
              placement="bottom"
              content={
                <p className="form-input-help" style={{ marginTop: '0' }}>
                  {operatorTip}
                </p>
              }
            >
              <select
                className="tc-15-select s"
                style={{ minWidth: '50px', maxWidth: '110px', height: '30px' }}
                value={defaultAffinityOperator}
                onChange={e => {
                  actions.editWorkload.updateAffinityRule(type, item.id + '', {
                    operator: e.target.value,
                    values: isNeedValues ? item.values : ''
                  });
                  actions.validate.workload.validateNodeAffinityRuleValue(type, item.id + '');
                }}
              >
                {operatorTypeList}
              </select>
            </Bubble>
          </div>
          <InputField
            className="mr10 tc-15-input-text m"
            type="input"
            validator={item.v_values}
            value={item.values}
            placeholder={t("多个Label Value请以 ';' 分隔符隔开")}
            style={{ width: '210px' }}
            tipMode="popup"
            disabled={item.operator === 'Exists' || item.operator === 'DoesNotExist'}
            disabeldTip={t('DoesNotExist,Exists操作符不需要填写value')}
            onChange={value => {
              actions.editWorkload.updateAffinityRule(type, item.id + '', { values: value });
            }}
            onBlur={e => {
              actions.validate.workload.validateNodeAffinityRuleValue(type, item.id + '');
            }}
          />
          <span className="inline-help-text" style={{ height: '30px' }}>
            <LinkButton onClick={() => actions.editWorkload.deleteAffinityRule(type, item.id + '')}>
              <i className="icon-cancel-icon" />
            </LinkButton>
          </span>
        </div>
      );
    });
  }
  render() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { nodeAffinityType, nodeAffinityRule, nodeSelection, v_nodeSelection } = workloadEdit;
    let requiredCanAdd = isEmpty(nodeAffinityRule.requiredExecution[0].matchExpressions.filter(x => !x.key));
    let preferredCanAdd = isEmpty(
      nodeAffinityRule.preferredExecution[0].preference.matchExpressions.filter(x => !x.key)
    );
    return (
      <FormItem label={t('节点调度策略')}>
        <Radio.Group
          value={workloadEdit.nodeAffinityType}
          onChange={value => actions.editWorkload.selectNodeSelectType(value)}
        >
          <Radio style={{ lineHeight: '18px' }} name={affinityType.unset}>
            {t('不使用调度策略')}
          </Radio>
          {
            /// #if tke
            <Radio style={{ lineHeight: '18px' }} name={affinityType.node}>
              {t('指定节点调度')}
            </Radio>
            /// #endif
          }

          <Radio style={{ lineHeight: '18px' }} name={affinityType.rule}>
            {t('自定义调度规则')}
          </Radio>
        </Radio.Group>
        <p className="text-label">{t('可根据调度规则，将Pod调度到符合预期的Label的节点中。')}</p>
        {nodeAffinityType === affinityType.node ? (
          <div
            className={'form-unit specific'}
            style={{
              backgroundColor: '#f2f2f2',
              padding: '5px 10px 15px',
              maxWidth: '350px',
              maxHeight: '150px',
              minHeight: '50px',
              overflowY: 'auto'
            }}
          >
            {this._renderComputerList()}
          </div>
        ) : nodeAffinityType === affinityType.rule ? (
          <div className="up-date">
            <div className="as-sel-box">
              <ul>
                <FormItem
                  label={t('强制满足条件')}
                  tips={t('调度期间如果满足亲和性条件则调度到对应node，如果没有节点满足条件则调度失败。')}
                >
                  {this._renderAffinityRuleList('required')}
                  <div>
                    <LinkButton
                      disabled={!requiredCanAdd}
                      tipDirection={'top'}
                      errorTip={t('请先完成待编辑项')}
                      onClick={() => actions.editWorkload.addAffinityRule('required')}
                    >
                      {t('新增规则')}
                    </LinkButton>
                  </div>
                </FormItem>
              </ul>
            </div>
            <div className="as-sel-box">
              <ul>
                <FormItem
                  label={t('尽量满足条件')}
                  tips={t('调度期间如果满足亲和性条件则调度到对应node，如果没有节点满足条件则随机调度到任意节点。')}
                >
                  {this._renderAffinityRuleList('preferred')}
                  <div>
                    <LinkButton
                      disabled={!preferredCanAdd}
                      tipDirection={'top'}
                      errorTip={t('请先完成待编辑项')}
                      onClick={() => actions.editWorkload.addAffinityRule('preferred')}
                    >
                      {t('新增规则')}
                    </LinkButton>
                  </div>
                </FormItem>
              </ul>
            </div>
          </div>
        ) : (
          <noscript />
        )}
      </FormItem>
    );
  }
}
