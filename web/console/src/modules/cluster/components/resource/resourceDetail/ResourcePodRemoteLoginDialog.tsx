import * as React from 'react';
import { bindActionCreators, uuid } from '@tencent/qcloud-lib';
import { Modal, Button, Table, TableColumn, Bubble } from '@tea/component';

import { RootProps } from '../../ClusterApp';
import { connect } from 'react-redux';
import { allActions } from '../../../actions';
import * as classnames from 'classnames';
import { PodContainer } from '../../../models';
import { ContainerStatusMap } from '../../../constants/Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { stylize } from '@tea/component/table/addons/stylize';
const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(state => state, mapDispatchToProps)
export class ResourcePodRemoteLoginDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions, subRoot, route } = this.props,
      { resourceDetailState } = subRoot,
      { podSelection, isShowLoginDialog } = resourceDetailState;

    if (!isShowLoginDialog) {
      return <noscript />;
    }

    const cancel = () => {
      // 关闭 远程登录的弹窗
      actions.resourceDetail.pod.toggleLoginDialog();
      // 置空当前的pod的选项
      actions.resourceDetail.pod.podSelect([]);
    };

    let containers = podSelection[0].spec.containers,
      containerStatus = podSelection[0].status.containerStatuses;

    // 容器登录的web-console的网址
    let loginUrl = '';
    // 容器登录的命名空间
    let namespace: string = podSelection[0] ? podSelection[0].metadata.namespace : 'default';

    let columns: TableColumn<PodContainer>[] = [
      {
        key: 'name',
        header: t('容器名称'),
        width: '45%',
        render: x => (
          <div className="sl-editor-name">
            <span className="text-overflow m-width" title={x.name}>
              {x.name}
            </span>
          </div>
        )
      },
      {
        key: 'status',
        header: t('状态'),
        width: '25%',
        render: x => this._reduceContainerStatus(x, containerStatus)
      },
      {
        key: 'description',
        header: t('操作'),
        width: '25%',
        render: x => (
          <div className="text-left">
            <div className="sl-editor-name">
              <a
                href={
                  loginUrl +
                  '?clusterId=' +
                  route.queries['clusterId'] +
                  '&podId=' +
                  podSelection[0].metadata.name +
                  '&containerName=' +
                  (x.name ? x.name : '') +
                  '&namespace=' +
                  namespace
                }
                target="_blank"
              >
                {t('登录')}
              </a>
            </div>
          </div>
        )
      }
    ];

    let containersLength = containers.length;
    return (
      <Modal visible={true} caption={t('容器登录')} onClose={cancel} disableEscape={true}>
        <Modal.Body>
          <div className="docker-dialog jiqun">
            <div className="act-outline">
              <div className="act-summary">
                <p>
                  <span>
                    <Trans count={containersLength}>
                      该实例下共有<strong className="text-warning">{{ containersLength }}个</strong>容器
                    </Trans>
                  </span>
                </p>
              </div>
              <div className="del-colony-tb">
                <Table
                  columns={columns}
                  records={containers}
                  addons={[
                    stylize({
                      className: 'ovm-dialog-tablepanel',
                      bodyStyle: { overflowY: 'auto', height: 160 }
                    })
                  ]}
                />
              </div>
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }

  /** 处理容器的状态 */
  private _reduceContainerStatus(container: PodContainer, containerStatus: any[]) {
    let finder = containerStatus ? containerStatus.find(c => c.name === container.name) : undefined,
      statusKey = finder && Object.keys(finder.state)[0];

    return (
      <div>
        <span
          className={classnames(
            'text-overflow',
            ContainerStatusMap[statusKey] && ContainerStatusMap[statusKey].classname
          )}
        >
          {ContainerStatusMap[statusKey] ? ContainerStatusMap[statusKey].text : '-'}
        </span>
        {statusKey && statusKey !== 'running' && (
          <Bubble placement="right" content={finder.state[statusKey].reason || null}>
            <div className="tc-15-bubble-icon">
              <i className="tc-icon icon-what" />
            </div>
          </Bubble>
        )}
      </div>
    );
  }
}
