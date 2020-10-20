import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../ChartGroupApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {

  render() {
    const { route } = this.props;
    let urlParam = router.resolve(route);
    const { mode } = urlParam;
    return (
      <Justify
        left={
          <h2>
            {mode ? (
              <React.Fragment>
                <a href="javascript:history.go(-1);">
                  <Icon type="btnback" />
                </a>
                <span style={{ marginLeft: '10px' }}>{mode}</span>
              </React.Fragment>
            ) : (
              <Trans>模板仓库管理</Trans>
            )}
          </h2>
        }
      />
    );
  }
}
