import * as React from 'react';
import { RootProps } from '../../ClusterApp';
import { OperationState, isSuccessWorkflow } from '@tencent/qcloud-redux-workflow';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { allActions } from '../../../actions';
import { connect } from 'react-redux';
import { WorkflowDialog, TablePanelColumnProps, GridTable, FormPanel } from '../../../../common/components';
import { CreateResource, Resource } from '../../../models';
import { router } from '../../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { resourceConfig } from '../../../../../../config';
import { Text, TableColumn, Button, Justify } from '@tencent/tea-component';
import { stylize } from '@tencent/tea-component/lib/table/addons';
import { downloadCrt } from '../../../../../../helpers';
import { clearNodeSH } from '../../../constants/Config';

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
          <Justify
            left={
              <FormPanel.Switch
                text={'清理节点数据'}
                value={isOpenDownloadButton}
                onChange={() => {
                  this.setState({ isOpenDownloadButton: !isOpenDownloadButton });
                }}
              ></FormPanel.Switch>
            }
            right={
              isOpenDownloadButton && (
                <Button type="link" onClick={e => downloadCrt(clearNodeSH, `clear${Date.now()}.sh`)}>
                  {t('下载脚本')}
                </Button>
              )
            }
          />
          {isOpenDownloadButton && (
            <FormPanel.InlineText parent="p" theme="danger">
              {t('注:清除数据不可恢复，请确认')}
            </FormPanel.InlineText>
          )}
        </div>
      </WorkflowDialog>
    );
  }
}
