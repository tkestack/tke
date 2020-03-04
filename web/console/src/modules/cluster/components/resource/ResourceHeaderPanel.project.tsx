import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Justify } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { firstRouterNameMap } from '../../../../../config';
import { apiVersion } from '../../../../../config/resource/common';
import { ResourceConfigVersionMap } from '../../../../../config/resourceConfig';
import { allActions } from '../../actions';
import { router } from '../../router';
import { RootProps } from '../ClusterApp';
import { IsInNodeManageDetail } from './resourceDetail/ResourceDetail';

interface ResourceHeaderPanelState {
  /** 当前的游标，是resourceTablePanel 还是 ResourceDetailPanel */
  currentIndex?: number;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceHeaderPanel extends React.Component<RootProps, ResourceHeaderPanelState> {
  constructor(props: RootProps, context) {
    super(props, context);

    // 如果是在list当中，则返回跳回集群列表，如果是detail，则为1
    let urlParmas = router.resolve(props.route);
    this.state = {
      currentIndex: urlParmas['mode'] === 'list' ? 0 : 1
    };
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let oldMode = this.props.subRoot.mode,
      newMode = nextProps.subRoot.mode;

    // 这里是去更新当前的currentIndex，决定路由如何变化
    if (oldMode !== newMode) {
      this.setState({
        currentIndex: newMode === 'list' ? 0 : 1
      });
    }
  }

  render() {
    let { subRoot, route, clusterVersion } = this.props,
      urlParams = router.resolve(route),
      { mode } = subRoot;

    // 回到集群列表即 ''，回到
    let breads = []; //['cluster', mode];

    if (mode === 'detail') {
      let { resourceIns, np } = route.queries;
      let headTitle: string = apiVersion[ResourceConfigVersionMap(clusterVersion)][urlParams['resourceName']].headTitle;
      let nsName: string = IsInNodeManageDetail(urlParams['type']) || np === '' || np === undefined ? '' : `(${np})`;
      let breadName = `${headTitle}:${resourceIns}${nsName}`;
      breads.push(breadName);
    }
    let breadHeader = this._reduceBreadCrumbs(breads);

    return (
      <Justify
        left={
          <React.Fragment>
            <a
              href="javascript:;"
              style={{ marginRight: '5px' }}
              className="back-link"
              onClick={this._handleClickForFirstBread.bind(this, this.state.currentIndex)}
            >
              <i className="btn-back-icon" />
              {t('返回')}
            </a>
            {breadHeader}
          </React.Fragment>
        }
      />
    );
  }

  /** 生成面包屑导航 /集群/集群ID(集群名)/... */
  private _reduceBreadCrumbs(breads: any[]) {
    let { route, cluster } = this.props;

    let content: JSX.Element;
    content = (
      <ol className={'breadcrumb'}>
        {breads.map((bread, index) => {
          let routerName = '';
          if (index === 0) {
            //   routerName = firstRouterNameMap[bread];
            // } else if (index === 1) {
            //   routerName = `${route.queries['clusterId']}(${
            //     cluster.selection ? cluster.selection.spec.displayName : '-'
            //   })`;
            // } else if (index === 2) {
            routerName = bread;
          }

          return (
            <li key={index}>
              {index === breads.length - 1 ? (
                <span>{routerName}</span>
              ) : (
                <a
                  href="javascript:;"
                  onClick={e => {
                    this._handleClickForFirstBread(index);
                  }}
                >
                  {routerName}
                </a>
              )}
            </li>
          );
        })}
      </ol>
    );

    return content;
  }

  /** 面包屑的点击处理 */
  private _handleClickForFirstBread(index: number) {
    let { region, route } = this.props,
      urlParams = router.resolve(route);

    if (index === 0) {
      router.navigate({}, { rid: region.selection.value + '' });
    } else if (index === 1) {
      let queries = JSON.parse(JSON.stringify(Object.assign({}, route.queries, { resourceIns: undefined })));
      let param = JSON.parse(JSON.stringify(Object.assign({}, urlParams, { mode: 'list', tab: undefined })));
      router.navigate(param, queries);
    }
  }
}
