import * as React from 'react';
import { connect } from 'react-redux';
import { Justify, Icon } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { RootProps } from '../RoleApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {

  render() {
    const { route } = this.props;
    let urlParam = router.resolve(route);
    const { sub } = urlParam;
    return (
      <Justify
        left={
          <h2>
            {sub ? (
              <React.Fragment>
                <a href="javascript:history.go(-1);">
                  <Icon type="btnback" />
                </a>
                <span style={{ marginLeft: '10px' }}>{sub}</span>
              </React.Fragment>
            ) : (
              <Trans>角色管理</Trans>
            )}
          </h2>
        }
      />
    );
  }
}
