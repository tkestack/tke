import * as React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';
import { YamlEditorPanel } from '../YamlEditorPanel';
import { FetchState } from '@tencent/qcloud-redux-fetcher';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Card, Select, Row, Col } from '@tencent/tea-component';

// 加载中的样式
const loadingElement: JSX.Element = (
  <div>
    <i className="n-loading-icon" />
    &nbsp; <span className="text">{t('加载中...')}</span>
  </div>
);

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceYamlPanel extends React.Component<RootProps, {}> {
  componentDidMount() {
    let { actions, subRoot } = this.props,
      { resourceOption } = subRoot,
      { resourceSelection } = resourceOption;

    // 这里是从 列表跳入到详情页的时候，进行yaml的拉取
    resourceSelection.length && actions.resourceDetail.fetchResourceYaml.fetch();
  }

  componentWillReceiveProps(nextProps: RootProps) {
    let newResourceSelection = nextProps.subRoot.resourceOption.resourceSelection,
      oldResourceSelection = this.props.subRoot.resourceOption.resourceSelection;

    // 在详情页面直接刷新的，由于请求顺序的原因，此时并没有resourceSelection，会导致请求yaml失败
    if (oldResourceSelection.length === 0 && newResourceSelection.length === 1) {
      this.props.actions.resourceDetail.fetchResourceYaml.fetch();
    }
  }

  render() {
    let { subRoot, actions } = this.props,
      {
        resourceDetailState,
        resourceOption,
        resourceInfo,
        detailResourceOption: { detailResourceName, detailResourceList, detailResourceSelection }
      } = subRoot,
      { resourceList } = resourceOption,
      { yamlList } = resourceDetailState;

    const yamlData = yamlList.data.recordCount ? yamlList.data.records[0] : t('暂无YAML配置');

    let isNeedLoading =
      resourceList.fetched !== true ||
      resourceList.fetchState === FetchState.Fetching ||
      yamlList.fetchState === FetchState.Fetching;

    return (
      <Card>
        <Card.Body>
          {isNeedLoading ? (
            loadingElement
          ) : (
            <React.Fragment>
              {resourceInfo.requestType.useDetailInfo && (
                <Row>
                  <Col className={'tea-mb-2n'}>
                    {t('对象选择')}
                    <Select
                      className="tea-ml-2n"
                      options={resourceInfo.requestType.detailInfoList['yaml']}
                      value={detailResourceName}
                      onChange={value => actions.resource.initDetailResourceName(value)}
                    />
                    <Select
                      className="tea-ml-2n"
                      options={detailResourceList}
                      value={detailResourceSelection}
                      onChange={value => actions.resource.selectDetailResouceIns(value)}
                    />
                  </Col>
                </Row>
              )}
              <YamlEditorPanel config={yamlData} readOnly={true} />
            </React.Fragment>
          )}
        </Card.Body>
      </Card>
    );
  }
}
