import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../../common/components';
import { SelectMultiple, Card, Modal, Icon, Justify, Row, Col, MediaObject, Tag, Pagination } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../../router';
import { allActions } from '../../../actions';
import { Chart } from '../../../models';
import { RootProps } from '../ChartApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface ChartTableState {
  chartGroupID?: string[];
}

@connect(state => state, mapDispatchToProps)
export class TablePanel extends React.Component<RootProps, ChartTableState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      chartGroupID: []
    };
  }

  render() {
    let {
      chartList: {
        selection,
        list,
        query: {
          search,
          keyword,
          paging: { pageIndex, pageSize }
        }
      },
      chartGroupList,
      route,
      actions
    } = this.props;
    let finalList = list.data.records;
    let cgMap = {};
    let cgSelect = [];
    chartGroupList.list.data.records.forEach(cg => {
      cgMap[cg.spec.name] = cg;
      cgSelect.push({ text: cg.spec.name + '(' + cg.spec.displayName + ')', value: cg.spec.name });
    });
    // 过滤仓库
    let filterList = finalList;
    if (this.state.chartGroupID.length > 0) {
      filterList = list.data.records.filter(c => {
        return this.state.chartGroupID.indexOf(c.spec.chartGroupName) > -1;
      });
    }
    let filterListConut = filterList.length;
    filterList = filterList.slice((pageIndex - 1) * pageSize, pageIndex * pageSize);

    const typeMap = {
      personal: '个人仓库',
      project: '业务仓库',
      system: '系统仓库'
    };
    const visibilityMap = {
      Public: '公有',
      Private: '私有'
    };

    return (
      <React.Fragment>
        <div style={{ backgroundColor: '#fff', boxShadow: '0 2px 3px 0 rgba(0,0,0,.2)', padding: '20px' }}>
          <Row gap={10}>
            <Col span={16} key={'select'}>
              <SelectMultiple
                searchable={true}
                allowEmpty={true}
                size={'full'}
                value={this.state.chartGroupID}
                options={cgSelect}
                appearence={'filter'}
                placeholder={'仓库过滤'}
                onChange={value => {
                  this.setState({ chartGroupID: value });
                }}
              ></SelectMultiple>
            </Col>
          </Row>
          <Row gap={10}>
            {filterList.map((chart, index) => (
              <Col span={8} key={chart.metadata.name}>
                <Card style={{ boxShadow: '0 2px 3px 1px rgba(0,0,0,.2)' }}>
                  <Card.Body
                    title={
                      <div
                        style={{ cursor: 'pointer' }}
                        onClick={e => {
                          const cg = cgMap[chart.spec.chartGroupName];
                          router.navigate(
                            { mode: 'detail', sub: 'chart' },
                            Object.assign({}, route.queries, {
                              chart: chart.metadata.name,
                              cg: chart.metadata.namespace,
                              chartName: chart.spec.name,
                              cgName: chart.spec.chartGroupName,
                              prj: cg && cg.spec.projects && cg.spec.projects.length > 0 ? cg.spec.projects[0] : ''
                            })
                          );
                        }}
                      >
                        <Justify
                          className={`app-tke-fe-apply__card-hd`}
                          left={
                            <MediaObject
                              media={
                                <img
                                  style={{ height: '48px' }}
                                  src={
                                    chart.lastVersion && chart.lastVersion.icon
                                      ? chart.lastVersion.icon
                                      : '/static/image/chart-icon-48x48.png'
                                  }
                                  alt=""
                                />
                              }
                              align="middle"
                            >
                              {''}
                            </MediaObject>
                          }
                        />
                        {chart.spec.chartGroupName + ' / ' + chart.spec.name}
                      </div>
                    }
                    subtitle={
                      <React.Fragment>
                        <Tag style={{ marginLeft: '-5px' }} theme="primary" key={'visibility'}>
                          {cgMap[chart.spec.chartGroupName]
                            ? visibilityMap[cgMap[chart.spec.chartGroupName].spec.visibility]
                            : '-'}
                        </Tag>
                        <Tag theme="primary" key={'type'}>
                          {cgMap[chart.spec.chartGroupName] ? typeMap[cgMap[chart.spec.chartGroupName].spec.type] : '-'}
                        </Tag>
                        <Tag theme="primary" key={chart.lastVersion ? chart.lastVersion.version : 'version'}>
                          {chart.lastVersion ? chart.lastVersion.version : t('空')}
                        </Tag>
                      </React.Fragment>
                    }
                  >
                    <p
                      title={chart.lastVersion ? chart.lastVersion.description : ''}
                      style={{
                        height: '36px',
                        display: '-webkit-box',
                        WebkitBoxOrient: 'vertical',
                        WebkitLineClamp: 2,
                        overflow: 'hidden',
                        fontSize: '12px',
                        color: '#888'
                      }}
                    >
                      {chart.lastVersion ? chart.lastVersion.description : ''}
                    </p>
                  </Card.Body>
                </Card>
              </Col>
            ))}
          </Row>
          <br />
          <Pagination
            pageIndex={pageIndex}
            pageSize={pageSize}
            recordCount={filterListConut}
            onPagingChange={query => {
              if (query.pageIndex > Math.ceil(filterListConut / query.pageSize)) {
                query.pageIndex = 1;
              }
              actions.chart.list.changePaging(query);
            }}
            pageSizeOptions={[15, 24, 36, 72]}
          />
        </div>
      </React.Fragment>
    );
  }

  /** 渲染操作按钮 */
  _renderOperationCell = (chart: Chart) => {
    let { actions } = this.props;
    return (
      <React.Fragment>
        <LinkButton onClick={() => this._removeChart(chart)}>
          <Trans>删除</Trans>
        </LinkButton>
      </React.Fragment>
    );
  };

  _removeChart = async (chart: Chart) => {
    let { actions, route } = this.props;
    const yes = await Modal.confirm({
      message: t('确定删除模板：') + `${chart.spec.displayName}？`,
      description: <p className="text-danger">{t('删除该Chart后，相关数据将永久删除，请谨慎操作。')}</p>,
      okText: t('删除'),
      cancelText: t('取消')
    });
    if (yes) {
      let urlParam = router.resolve(route);
      actions.chart.list.removeChartWorkflow.start([chart], {
        repoType: urlParam['tab'] || 'all'
      });
      actions.chart.list.removeChartWorkflow.perform();
    }
  };
}
