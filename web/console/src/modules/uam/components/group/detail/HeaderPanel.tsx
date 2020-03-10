import * as React from 'react';
import { connect } from 'react-redux';
import { RootProps } from '../GroupApp';
import { Justify, Icon } from '@tencent/tea-component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { bindActionCreators } from '@tencent/ff-redux';
import { allActions } from '../../../actions';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class HeaderPanel extends React.Component<RootProps, {}> {

  goBack = () => {
    history.back();
  }

  render() {
    let { route } = this.props;
    let title = route.queries['groupName'];

    return (
      <Justify
        left={
          <React.Fragment>
            <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
              <Icon type="btnback" />
              {t('返回')}
            </a>
            <span className="line-icon">|</span>
            <h2>{title}</h2>
          </React.Fragment>
        }
      />
    );
  }
}
