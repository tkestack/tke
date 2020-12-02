import * as React from 'react';
import { connect } from 'react-redux';

import { Button, Dropdown, List, Modal, Select, Switch, Table, TagSearchBox, Text, Tooltip } from '@tea/component';
// import { TagSearchBox } from '../../../../common/components/tagsearchbox';
import { bindActionCreators, FetchState, insertCSS } from '@tencent/ff-redux';
import { ChartInstancesPanel } from '@tencent/tchart';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Justify } from '@tencent/tea-component/lib/justify';

import { resourceConfig } from '../../../../../../config';
import { dateFormatter, downloadCsv, reduceNs } from '../../../../../../helpers';
import { DisplayFiledProps, ResourceInfo } from '../../../../common/models';
import { includes, isEmpty } from '../../../../common/utils';
import { allActions } from '../../../actions';
import { Resource } from '../../../models';
import { MonitorPanelProps, resourceMonitorFields } from '../../../models/MonitorPanel';
import { router } from '../../../router';
import { RootProps } from '../../ClusterApp';
import { TellIsNeedFetchNS } from '../ResourceSidebarPanel';
import { PlatformContext, IPlatformContext, PlatformTypeEnum } from '@/Wrapper';

interface ResouceActionPanelState {
  /** 是否开启自动刷新 */
  isOpenAutoRenew?: boolean;

  /** searchbox的 */
  searchBoxValues?: any[];

  /** 搜索框当中的搜索的数量 */
  searchBoxLength?: number;

  /** 监控组件属性 */
  monitorPanelProps?: MonitorPanelProps;
}

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourceActionPanel extends React.Component<RootProps, ResouceActionPanelState> {
  static contextType = PlatformContext;
  context: IPlatformContext;

  constructor(props, context) {
    super(props, context);
    this.state = {
      isOpenAutoRenew: false,
      searchBoxValues: [],
      searchBoxLength: 0
    };
  }

  render() {
    let { route, subRoot } = this.props,
      urlParams = router.resolve(route);

    let kind = urlParams['type'],
      resourceName = urlParams['resourceName'];

    let monitorButton = null;
    monitorButton = ['deployment', 'statefulset', 'daemonset', 'tapp'].includes(resourceName) && this._renderMonitorButton();

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <React.Fragment>
              {this._renderCreateButton()}
              {monitorButton}
            </React.Fragment>
          }
          right={
            <React.Fragment>
              {TellIsNeedFetchNS(resourceName) && this._renderNamespaceSelect()}
              {this._renderTagSearchBox()}
              {this._renderAutoRenew()}
              {this._renderManualRenew()}
              {this._renderDownload()}
            </React.Fragment>
          }
        />
        {this.state && this.state.monitorPanelProps && (
          <Modal
            visible={true}
            caption={this.state.monitorPanelProps.title}
            onClose={() => this.setState({ monitorPanelProps: undefined })}
            size={1050}
          >
            <Modal.Body>
              <ChartInstancesPanel
                tables={this.state.monitorPanelProps.tables}
                groupBy={this.state.monitorPanelProps.groupBy}
                instance={this.state.monitorPanelProps.instance}
              >
                {this.state.monitorPanelProps.headerExtraDOM}
              </ChartInstancesPanel>
            </Modal.Body>
          </Modal>
        )}
      </Table.ActionPanel>
    );
  }
  _handleMonitor() {
    let { subRoot, route } = this.props,
      { resourceOption } = subRoot;

    this.setState({
      monitorPanelProps: {
        title: t('工作负载监控'),
        tables: [
          {
            fields: resourceMonitorFields,
            table: 'k8s_workload',
            conditions: [
              ['tke_cluster_instance_id', '=', route.queries.clusterId],
              ['workload_kind', '=', subRoot.resourceInfo.headTitle],
              ['namespace', '=', reduceNs(this.props.route.queries['np'])]
            ]
          }
        ],
        groupBy: [{ value: 'workload_name' }],
        instance: {
          columns: [{ key: 'workload_name', name: t('工作负载名称') }],
          list: resourceOption.ffResourceList.list.data.records.map(ins => ({
            workload_name: ins.metadata.name,
            isChecked:
              !resourceOption.resourceMultipleSelection.length ||
              resourceOption.resourceMultipleSelection.find(item => item.metadata.name === ins.metadata.name)
          }))
        }
      } as MonitorPanelProps
    });
  }

  private _renderMonitorButton() {
    return (
      <Button
        type="primary"
        onClick={() => {
          this._handleMonitor();
        }}
      >
        {t('监控')}
      </Button>
    );
  }

  /** render新建按钮 */
  private _renderCreateButton() {
    let { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const isShow = !isEmpty(resourceInfo) && resourceInfo.actionField && resourceInfo.actionField.create.isAvailable;
    return isShow ? (
      <Button
        type="primary"
        onClick={() => {
          this._handleClickForCreate();
        }}
      >
        {t('新建')}
      </Button>
    ) : (
      <noscript />
    );
  }

  /** action for create button */
  private _handleClickForCreate() {
    let { route } = this.props,
      urlParams = router.resolve(route);
    router.navigate(Object.assign({}, urlParams, { mode: 'create' }), route.queries);
  }

  /** 生成命名空间选择列表 */
  private _renderNamespaceSelect() {
    let { actions, namespaceList, namespaceSelection } = this.props;

    let selectProps = {};

    if (this.context.type === PlatformTypeEnum.Business) {
      const groups = namespaceList.data.records.reduce((gr, { clusterDisplayName, clusterName }) => {
        const value = `${clusterDisplayName}(${clusterName})`;
        return { ...gr, [clusterName]: <Tooltip title={value}>{value}</Tooltip> };
      }, {});

      let options = namespaceList.data.recordCount
        ? namespaceList.data.records.map(item => {
            const text = `${item.clusterDisplayName}-${item.namespace}`;

            return {
              value: item.name,
              text: <Tooltip title={text}>{text}</Tooltip>,
              groupKey: item.clusterName,
              realText: text
            };
          })
        : [{ value: '', text: t('无可用命名空间'), disabled: true }];

      selectProps = {
        groups,
        options,
        filter: (inputValue, { realText }: any) => (realText ? realText.includes(inputValue) : true)
      };
    } else {
      let options = namespaceList.data.recordCount
        ? namespaceList.data.records.map((item, index) => ({
            value: item.name,
            text: item.displayName
          }))
        : [{ value: '', text: t('无可用命名空间'), disabled: true }];

      selectProps = {
        options
      };
    }

    return (
      <div style={{ display: 'inline-block', fontSize: '12px', verticalAlign: 'middle' }}>
        <Text theme="label" verticalAlign="middle">
          {t('命名空间')}
        </Text>
        <Tooltip>
          <Select
            {...selectProps}
            type="simulate"
            searchable
            appearence="button"
            size="s"
            style={{ width: '130px', marginRight: '5px' }}
            value={namespaceSelection}
            onChange={value => {
              actions.namespace.selectNamespace(value);
            }}
            placeholder={namespaceList.data.recordCount ? t('请选择命名空间') : t('无可用命名空间')}
          />
        </Tooltip>
      </div>
    );
  }

  /** 生成搜索框 */
  private _renderTagSearchBox() {
    let { subRoot } = this.props,
      { resourceInfo, resourceOption } = subRoot,
      { ffResourceList } = resourceOption;

    // const defaultValue = [{attr: {key: 'namespace',name: '命名空间'},values: [{name: namespaceSelection}]}];

    // attributes当中的 namepsace列表的values
    // const namespaceValues = namespaceList.data.recordCount? namespaceList.data.records.map((namespace, index) => { return { key: namespace.id, name: namespace.name }; }) : [];

    // tagSearch的过滤选项
    const attributes = [
      {
        type: 'input',
        key: 'resourceName',
        name: t('名称')
      }
    ];

    // 这里是因为展示命名空间的话，不需要展示namespace
    // let isNeedFetchNamespace = TellIsNeedFetchNS(resourceName);
    // if (isNeedFetchNamespace) {
    //   let tmp = {
    //     type: 'single',
    //     key: 'namespace',
    //     name: '命名空间',
    //     values: namespaceValues
    //   };

    //   attributes.push(tmp);
    // }

    // 受控展示的values
    // const values = resourceQuery.search ? this.state.searchBoxValues : isNeedFetchNamespace ? defaultValue : [];
    const values = ffResourceList.query.search ? this.state.searchBoxValues : [];

    const isShow = !isEmpty(resourceInfo) && resourceInfo.actionField && resourceInfo.actionField.search.isAvailable;

    return isShow ? (
      <div style={{ width: 350, display: 'inline-block' }}>
        <TagSearchBox
          className="myTagSearchBox"
          attributes={attributes}
          value={values}
          onChange={tags => {
            this._handleClickForTagSearch(tags);
          }}
        />
      </div>
    ) : (
      <noscript />
    );
  }

  /** 搜索框的操作，不同的搜索进行相对应的操作 */
  private _handleClickForTagSearch(tags) {
    let { actions, subRoot } = this.props,
      { resourceOption } = subRoot,
      { ffResourceList } = resourceOption;

    // 这里是控制tagSearch的展示
    this.setState({
      searchBoxValues: tags,
      searchBoxLength: tags.length
    });

    // 如果检测到 tags的长度变化，并且key为 resourceName 去掉了，则清除搜索条件
    if (
      tags.length === 0 ||
      (tags.length === 1 && ffResourceList.query.search && tags[0].attr && tags[0].attr.key !== 'resourceName')
    ) {
      actions.resource.changeKeyword('');
      actions.resource.performSearch('');
    }

    tags.forEach(tagItem => {
      let attrKey = tagItem.attr ? tagItem.attr.key : null;
      if (attrKey === 'resourceName' || attrKey === null) {
        // let oldSerach = resourceQuery.search;
        let search = tagItem.values[0].name;
        actions.resource.changeKeyword(search);
        actions.resource.performSearch(search);
      }
    });
  }

  /** 生成自动刷新按钮 */
  private _renderAutoRenew() {
    let { subRoot } = this.props,
      { resourceInfo } = subRoot;

    const isShow = !isEmpty(resourceInfo) && resourceInfo.actionField && resourceInfo.actionField.autoRenew.isAvailable;
    return isShow ? (
      <span>
        <span
          className="descript-text"
          style={{ display: 'inline-block', verticalAlign: 'middle', marginRight: '10px', fontSize: '12px' }}
        >
          {t('自动刷新')}
        </span>
        <Switch
          value={this.state.isOpenAutoRenew}
          onChange={checked => {
            this.setState({ isOpenAutoRenew: !this.state.isOpenAutoRenew });
          }}
          className="mr20"
        />
      </span>
    ) : (
      <noscript />
    );
  }

  /** 生成手动刷新按钮 */
  private _renderManualRenew() {
    let { actions, subRoot, namespaceSelection } = this.props,
      { resourceOption, resourceInfo } = subRoot,
      { ffResourceList } = resourceOption;

    let loading = ffResourceList.list.loading || ffResourceList.list.fetchState === FetchState.Fetching;
    const isShow =
      !isEmpty(resourceInfo) && resourceInfo.actionField && resourceInfo.actionField.manualRenew.isAvailable;
    return isShow ? (
      <Button
        icon="refresh"
        disabled={loading}
        onClick={e => {
          actions.resource.fetch();
        }}
        title={t('刷新')}
      />
    ) : (
      <noscript />
    );
  }

  /** 生成自动下载按钮 */
  private _renderDownload() {
    let { subRoot } = this.props,
      { resourceOption, resourceInfo } = subRoot,
      { ffResourceList } = resourceOption;

    let loading = ffResourceList.list.loading || ffResourceList.list.fetchState === FetchState.Fetching;
    const isShow = !isEmpty(resourceInfo) && resourceInfo.actionField && resourceInfo.actionField.download.isAvailable;
    return isShow ? (
      <Button
        icon="download"
        disabled={loading}
        title={t('导出全部')}
        onClick={() => this.downloadHandle(ffResourceList.list.data.records)}
      />
    ) : (
      <noscript />
    );
  }

  /** 导出数据 */
  private downloadHandle(resourceList: Resource[]) {
    let { clusterVersion, subRoot } = this.props,
      { resourceName } = subRoot;

    let resourceInfo: ResourceInfo = resourceConfig(clusterVersion)[resourceName];
    let rows = [],
      head = [];

    // 这里是去处理head当中显示的内容
    let headKeys = [],
      displayKeys = Object.keys(resourceInfo.displayField);

    displayKeys.forEach(item => {
      if (item !== 'operator') {
        let displayField: DisplayFiledProps = resourceInfo.displayField[item];
        headKeys.push(displayField.headTitle);
      }
    });
    head = headKeys;
    // 这里是去处理rows当中的信息
    resourceList.forEach((resource: Resource) => {
      // 每一行的数据
      let row = [];
      let rowInfos: DisplayFiledProps[] = [];
      displayKeys.forEach(item => {
        if (item !== 'operator') {
          rowInfos.push(resourceInfo.displayField[item]);
        }
      });

      // 获取最终的展示数据
      rowInfos.forEach(item => {
        let showData: any = [];
        item.dataField.forEach(field => {
          let dataFieldIns = field.split('.');
          let data: any = this._getFinalData(dataFieldIns, resource);
          // 如果返回的为 ''，即找不到这个对象，则使用配置文件中设定的默认值
          showData.push(data === '' ? item.noExsitedValue : data);
        });

        showData = showData.length === 1 ? showData[0] : showData;

        let content;
        if (item.dataFormat === 'text' || item.dataFormat === 'status' || item.dataFormat === 'mapText') {
          content = showData;
        } else if (item.dataFormat === 'labels') {
          content = this._reduceLabelsForData(showData);
        } else if (item.dataFormat === 'time') {
          content = dateFormatter(new Date(showData), 'YYYY-MM-DD HH:mm:ss');
        } else if (item.dataFormat === 'ip') {
          content =
            typeof showData === 'string'
              ? showData
              : t('负载均衡IP：') + showData[0] + '\n' + t('服务IP：') + showData[1];
        } else if (item.dataFormat === 'replicas') {
          content = showData[0] + '、' + showData[1];
        } else if (item.dataFormat === 'ingressType') {
          let ingressId = showData['kubernetes.io/ingress.qcloud-loadbalance-id'] || '-';
          content = ingressId + '\n' + t('应用型负载均衡');
        } else if (item.dataFormat === 'ingressRule') {
          content = this._reduceIngressRule(showData, resource);
        } else {
          content = showData;
        }

        row.push(content);
      });

      rows.push(row);
    });
    downloadCsv(rows, head, 'tke_' + resourceName + '_' + new Date().getTime() + '.csv');
  }

  /** 获得labels的最终展示 */
  private _reduceLabelsForData(labels) {
    let showData = '',
      keys;

    // 如果不是数组，showData就是Labels本身
    if (typeof labels === 'string') {
      showData = labels;
    } else {
      keys = Object.keys(labels);
      keys.forEach((item, index) => {
        showData += item + ':' + labels[item];
        if (index !== keys.length - 1) {
          showData += '\n';
        }
      });
    }
    return showData;
  }

  /** 获得ingress的后端服务的信息 */
  private _reduceIngressRule(showData: any, resource: Resource) {
    let data;

    let httpRules =
        showData['kubernetes.io/ingress.http-rules'] && showData['kubernetes.io/ingress.http-rules'] !== 'null'
          ? JSON.parse(showData['kubernetes.io/ingress.http-rules'])
          : [],
      httpsRules =
        showData['kubernetes.io/ingress.https-rules'] && showData['kubernetes.io/ingress.https-rules'] !== 'null'
          ? JSON.parse(showData['kubernetes.io/ingress.https-rules'])
          : [];

    httpRules = httpRules.map(item => Object.assign({}, item, { protocol: 'http' }));
    httpsRules = httpsRules.map(item => Object.assign({}, item, { protocol: 'https' }));

    const getDomain = rule => {
      return `${rule.protocol}://${
        rule.host ? rule.host : resource.status.loadBalancer.ingress ? resource.status.loadBalancer.ingress[0].ip : '-'
      }${rule.path}`;
    };

    let finalRules = [...httpRules, ...httpsRules];

    data = finalRules
      .map(item => {
        return getDomain(item) + '-->' + item.backend.serviceName + ':' + item.backend.servicePort;
      })
      .join('\n');

    return data;
  }

  /** 获取最终展示的数据 */
  private _getFinalData(dataFieldIns, resource: Resource) {
    let result = resource;

    for (let index = 0; index < dataFieldIns.length; index++) {
      // 如果result不为一个 Object，则遍历结束
      if (typeof result !== 'object') {
        break;
      }
      result = result[dataFieldIns[index]]; // 这里做一下处理，防止因为配错找不到
    }

    // 返回空值，是因为如果不存在值，则使用配置文件的默认值
    return result || '';
  }
}
