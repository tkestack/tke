import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Icon, Justify } from '@tencent/tea-component';

import { allActions } from '../actions';
import { router } from '../router';
import { RootProps } from './ProjectApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class SubpageHeadPanel extends React.Component<RootProps, {}> {
  render() {
    let { route, project } = this.props,
      projectName = project.selections[0] ? project.selections[0].spec.displayName : '',
      projectId = project.selections[0] ? project.selections[0].metadata.name : '';

    return (
      <React.Fragment>
        <Justify
          left={
            <React.Fragment>
              <a href="javascript:;" className="back-link" onClick={this.goBack.bind(this)}>
                <Icon type="btnback" />
                {t('返回')}
              </a>
              {projectId && <h2>{`${projectName}(${projectId})`}</h2>}
            </React.Fragment>
          }
        />
      </React.Fragment>
    );
  }

  goBack() {
    router.navigate({});
    // history.back();
  }
}
