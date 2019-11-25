import { Switch, Justify, Table, Row, Col, Select } from '@tea/component';
import { bindActionCreators, insertCSS } from '@tencent/qcloud-lib';
import * as React from 'react';
import { connect } from 'react-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';
import { TipInfo } from '../../../../common/components';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

insertCSS(
  'EventActionPanel',
  `
.tc-large-width{
    width: 250px;
}
`
);

interface ResouceEventPanelState {
  /** 是否开启自动刷新 */
  isAutoRenew?: boolean;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceDetailEventActionPanel extends React.Component<RootProps, ResouceEventPanelState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      isAutoRenew: true
    };
  }

  render() {
    let { actions, subRoot } = this.props,
      {
        resourceDetailState,
        resourceInfo,
        detailResourceOption: { detailResourceName, detailResourceSelection, detailResourceList }
      } = subRoot;

    return (
      <Table.ActionPanel>
        <TipInfo>{t('资源事件只保存最近1小时内发生的事件，请尽快查阅。')}</TipInfo>
        <Justify
          left={
            resourceInfo.requestType.useDetailInfo ? (
              <Row>
                <Col style={{ fontSize: '12px' }} className="tea-mb-2n">
                  {t('对象选择')}
                  <Select
                    className="tea-ml-2n"
                    options={resourceInfo.requestType.detailInfoList['event']}
                    value={detailResourceName}
                    onChange={value => actions.resource.initDetailResourceName(value)}
                  />
                  <Select
                    options={detailResourceList}
                    value={detailResourceSelection}
                    onChange={value => actions.resource.selectDetailResouceIns(value)}
                  />
                </Col>
              </Row>
            ) : null
          }
          right={
            <React.Fragment>
              <span
                className="descript-text"
                style={{ display: 'inline-block', verticalAlign: 'middle', marginRight: '10px', fontSize: '12px' }}
              >
                {t('自动刷新')}
              </span>
              <Switch
                value={this.state.isAutoRenew}
                onChange={checked => this._handleSwitch(checked)}
                className="mr20"
              />
            </React.Fragment>
          }
        />
      </Table.ActionPanel>
    );
  }

  private _handleSwitch(isChecked: boolean) {
    let { actions, route } = this.props;

    if (!isChecked) {
      actions.resourceDetail.event.clearPolling();
    } else {
      // 进行事件的拉取
      actions.resourceDetail.event.poll();
    }

    this.setState({ isAutoRenew: isChecked });
  }
}
