import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { allActions } from '../../../actions';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { Bubble } from '@tea/component';
import { connect } from 'react-redux';
import { isEmpty } from '../../../../common/utils';
import { FormItem, LinkButton } from '../../../../common/components';
import * as classnames from 'classnames';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class EditResourceLabelPanel extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { workloadLabels } = workloadEdit;

    // 判断，只有key 和 value 都为非空时，才能进行添加的操作
    let canAdd = isEmpty(workloadLabels.filter(x => !x.labelKey || !x.labelValue));

    return (
      <FormItem label={t('标签')}>
        <div className="form-unit is-success">
          {this._renderLabelList()}
          <div>
            <LinkButton
              disabled={!canAdd}
              errorTip={t('请先完成待编辑项')}
              onClick={() => {
                actions.editWorkload.addLabels();
              }}
            >
              {t('新增变量')}
            </LinkButton>
            <p className="text-label">
              {t('只能包含大小写字母、数字及分隔符"-"、"_"和"."，且必须以大小写字母、数字开头和结尾')}
            </p>
          </div>
        </div>
      </FormItem>
    );
  }

  /** 展示label的选项 */
  private _renderLabelList() {
    let { actions, subRoot } = this.props,
      { workloadEdit } = subRoot,
      { workloadLabels } = workloadEdit;

    return workloadLabels.map((label, index) => {
      let isDisabled = index === 0;
      return (
        <div className="code-list" key={index} style={{ marginBottom: '5px' }}>
          <div
            className={classnames({ 'is-error': label.v_labelKey.status === 2 })}
            style={{ display: 'inline-block' }}
          >
            <Bubble
              placement="bottom"
              content={(() => {
                if (isDisabled) {
                  return t('默认标签不可更改');
                }
                if (label.v_labelKey.status === 2) {
                  return label.v_labelKey.message;
                }
                return null;
              })()}
            >
              <input
                type="text"
                disabled={isDisabled}
                placeholder="Key"
                className="tc-15-input-text m"
                style={{ width: '128px' }}
                value={label.labelKey}
                maxLength={63}
                onChange={e => actions.editWorkload.updateLabels({ labelKey: e.target.value }, label.id + '')}
                onBlur={e => actions.validate.workload.validateAllWorkloadLabelKey()}
              />
            </Bubble>
            <span className="inline-help-text">=</span>
            <Bubble
              placement="bottom"
              content={label.v_labelValue.status === 2 ? <p>{label.v_labelValue.message}</p> : null}
            >
              <input
                type="text"
                placeholder="Value"
                className="tc-15-input-text m"
                style={{ width: '128px', marginLeft: '5px' }}
                value={label.labelValue}
                maxLength={63}
                onChange={e => actions.editWorkload.updateLabels({ labelValue: e.target.value }, label.id + '')}
                onBlur={e => actions.validate.workload.validateAllWorkloadLabelValue()}
              />
            </Bubble>
            <span className="inline-help-text">
              <LinkButton
                errorTip={t('默认标签不可删除')}
                disabled={isDisabled}
                onClick={() => actions.editWorkload.deleteLabels(label.id + '')}
              >
                <i className="icon-cancel-icon" />
              </LinkButton>
            </span>
          </div>
          {/* <span className="inline-help-text">
            <LinkButton
              errorTip={t('默认标签不可删除')}
              disabled={isDisabled}
              onClick={() => actions.editWorkload.deleteLabels(label.id + '')}
            >
              <i className="icon-cancel-icon" />
            </LinkButton>
          </span> */}
        </div>
      );
    });
  }
}
