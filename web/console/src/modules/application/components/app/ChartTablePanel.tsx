import * as React from 'react';
import { connect } from 'react-redux';
import { TablePanel as CTablePanel } from '@tencent/ff-component';
import { LinkButton, emptyTips } from '../../../common/components';
import { SelectMultiple, Card, Modal, Icon, Justify, Row, Col, MediaObject, Tag, Pagination } from '@tea/component';
import { bindActionCreators } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { router } from '../../router';
import { allActions } from '../../actions';
import { Chart } from '../../models';
import { RootProps } from './AppContainer';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), {
    dispatch
  });

interface ChartTableState {
  chartGroupID?: string[];
  selectedChart?: string;
}

interface Props extends RootProps {
  onSelectChart?: Function;
  SelectedChart?: string;
}
@connect(state => state, mapDispatchToProps)
export class ChartTablePanel extends React.Component<Props, ChartTableState> {
  constructor(props, context) {
    super(props, context);
    this.state = {
      chartGroupID: [],
      selectedChart: ''
    };
  }

  render() {
    let {
      chartList: {
        list,
        query: {
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
        <div style={{ backgroundColor: '#f2f2f2', boxShadow: '0 2px 3px 0 rgba(0,0,0,.2)', padding: '10px' }}>
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
                <div
                  style={{ cursor: 'pointer' }}
                  onClick={e => {
                    this.setState({ selectedChart: chart.metadata.name });
                    actions.chart.list.select(chart);

                    const cg = cgMap[chart.spec.chartGroupName];
                    const prj = cg && cg.spec.projects && cg.spec.projects.length > 0 ? cg.spec.projects[0] : '';
                    this.props.onSelectChart && this.props.onSelectChart(chart, prj);
                  }}
                >
                  <Card
                    style={{
                      boxShadow: '0 2px 3px 1px rgba(0,0,0,.2)',
                      border:
                        chart.metadata.name ===
                        (this.state.selectedChart ? this.state.selectedChart : this.props.SelectedChart)
                          ? 'solid #ff9d00'
                          : 'solid #fff'
                    }}
                  >
                    <Card.Body
                      title={
                        <div style={{ cursor: 'pointer' }} onClick={e => {}}>
                          {chart.spec.chartGroupName +
                            ' / ' +
                            chart.spec.name +
                            '(' +
                            (cgMap[chart.spec.chartGroupName]
                              ? typeMap[cgMap[chart.spec.chartGroupName].spec.type]
                              : '-') +
                            ')'}
                        </div>
                      }
                      subtitle={
                        <React.Fragment>
                          <Tag key={chart.lastVersion ? chart.lastVersion.version : 'empty'}>
                            {chart.lastVersion ? chart.lastVersion.version : t('空')}
                          </Tag>
                          <Tag key={'visibility'}>
                            {cgMap[chart.spec.chartGroupName]
                              ? visibilityMap[cgMap[chart.spec.chartGroupName].spec.visibility]
                              : '-'}
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
                </div>
              </Col>
            ))}
          </Row>
          <br />
          <Pagination
            style={{ backgroundColor: '#f2f2f2', marginTop: '20px' }}
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
}
