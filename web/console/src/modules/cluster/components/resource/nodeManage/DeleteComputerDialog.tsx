import * as React from 'react';
import { connect } from 'react-redux';

import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { TableColumn, Text } from '@tencent/tea-component';
import { stylize } from '@tencent/tea-component/lib/table/addons';

import { resourceConfig } from '../../../../../../config';
import { downloadCrt } from '../../../../../../helpers';
import { Clip, GridTable, WorkflowDialog } from '../../../../common/components';
import { allActions } from '../../../actions';
import { clearNodeSH } from '../../../constants/Config';
import { CreateResource, Resource } from '../../../models';
import { RootProps } from '../../ClusterApp';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class DeleteComputerDialog extends React.Component<RootProps, {}> {
  state = {
    isOpenDownloadButton: false
  };
  render() {
    let { isOpenDownloadButton } = this.state;
    let { actions, route, subRoot, region, clusterVersion } = this.props,
      {
        computerState: {
          deleteComputer,
          computerPodList,
          computerPodQuery,
          computer: { selections }
        }
      } = subRoot;
    let resourceIns = selections[0] && selections[0].spec.machineName ? selections[0].spec.machineName : '';
    let nodeName = selections[0] && selections[0].metadata.name;
    let resourceInfo = resourceConfig(clusterVersion).machines;
    // 需要提交的数据
    let resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      clusterId: route.queries['clusterId'],
      resourceIns
    };
    let colunms: TableColumn<Resource>[] = [
      {
        key: 'name',
        header: t('实例（Pod）名称'),
        width: '55%',
        render: x => (
          <Text parent="div" overflow>
            <span title={x.metadata.name}>{x.metadata.name}</span>
          </Text>
        )
      },
      {
        key: 'namespace',
        header: t('所属集群空间'),
        width: '45%',
        render: x => (
          <Text parent="div" overflow>
            <span title={x.metadata.namespace}>{x.metadata.namespace}</span>
          </Text>
        )
      }
    ];
    let computerPodCount = computerPodList.data.recordCount;
    // 这里主要是考虑在更新实例数量的时候，会调用删除接口删除hpa，不应该展示出dialog
    return (
      <WorkflowDialog
        caption={t('您确定要删除节点：{{nodeName}}吗？', {
          nodeName
        })}
        workflow={deleteComputer}
        action={actions.workflow.deleteComputer}
        params={region.selection ? region.selection.value : ''}
        targets={[resource]}
        isDisabledConfirm={resourceIns ? false : true}
      >
        <div className="docker-dialog jiqun">
          {
            <div className="act-outline">
              <div className="act-summary">
                <p>
                  <Trans count={computerPodCount}>
                    <span>
                      节点包含<strong className="text-warning">{{ computerPodCount }}个</strong>实例
                    </span>
                  </Trans>
                </p>
              </div>
              <div className="del-colony-tb">
                <GridTable
                  columns={colunms}
                  emptyTips={<div className="text-center">{t('节点的实例（Pod）列表为空')}</div>}
                  listModel={{
                    list: computerPodList,
                    query: computerPodQuery
                  }}
                  actionOptions={actions.computerPod}
                  addons={[
                    stylize({
                      className: 'ovm-dialog-tablepanel',
                      bodyStyle: { overflowY: 'auto', height: 160, minHeight: 100 }
                    })
                  ]}
                  isNeedCard={false}
                />
              </div>
            </div>
          }
          <Text parent="p" className="tea-mt-1n">
            <Text>{t('如需清理节点上的数据，您可以通过以下脚本进行手动清理，')}</Text>
            <Text theme="danger">{t('数据清理后不可恢复，如节点存在混用情况，请谨慎执行。')}</Text>
          </Text>
          <div className="rich-textarea hide-number">
            <Clip target={'#certificationAuthority'} className="copy-btn">
              {t('复制')}
            </Clip>
            <a
              href="javascript:void(0)"
              onClick={e => downloadCrt(clearNodeSH, `clear${Date.now()}.sh`)}
              className="copy-btn"
              style={{ right: '50px' }}
            >
              {t('下载')}
            </a>
            <div className="rich-content" contentEditable={false}>
              <p
                className="rich-text"
                id="certificationAuthority"
                style={{
                  width: '432px',
                  whiteSpace: 'pre-wrap',
                  overflow: 'auto',
                  height: '300px'
                }}
              >
                {clearNodeSH}
              </p>
            </div>
          </div>
        </div>
      </WorkflowDialog>
    );
  }
}
