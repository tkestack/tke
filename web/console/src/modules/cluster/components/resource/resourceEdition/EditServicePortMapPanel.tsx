import * as React from 'react';
import { Bubble, Table, TableColumn } from '@tea/component';
import { uuid } from '@tencent/qcloud-lib';
import { FormItem, HeadBubble, InputField, TablePanelColumnProps } from '../../../../common/components';
import { PortMap } from '../../../models';
import { ProtocolList } from '../../../constants/Config';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { stylize } from '@tea/component/table/addons/stylize';
interface EditServicePortMapPanelProps {
  /** 添加端口映射的操作 */
  addPortMap: () => void;

  /** 删除端口映射的操作 */
  deletePortMap: (pId: string) => void;

  /** 当前的communicationType */
  communicationType: string;

  /** 端口映射的数组 */
  portsMap: PortMap[];

  /** 更新端口映射的配置 */
  updatePortMap: (obj: any, pId: string) => void;

  /** 校验端口协议是否正确 */
  validatePortProtocol: (value: any, pId: string) => void;

  /** 校验容器端口 */
  validateTargetPort: (value: any, pId: string) => void;

  /** 校验主机端口 */
  validateNodePort: (value: any, pId: string) => void;

  /** 校验服务端口 */
  validateServicePort: (value: any, pId: string) => void;

  /** 是否展示该部分内容 */
  isShow?: boolean;
}

export class EditServicePortMapPanel extends React.Component<EditServicePortMapPanelProps, {}> {
  render() {
    let { addPortMap, isShow = true } = this.props;

    return isShow ? (
      <FormItem className="vm" label={t('端口映射')}>
        {this._renderPortMap()}
        <a href="javascript:;" className="more-links-btn" onClick={addPortMap}>
          {t('添加端口映射')}
        </a>
      </FormItem>
    ) : (
      <noscript />
    );
  }

  /** 根据当前的访问方式更改服务端口的展示内容 */
  private _renderPortHeader() {
    let { communicationType } = this.props;

    let content: any;

    switch (communicationType) {
      case 'LoadBalancer':
      case 'SvcLBTypeInner':
        content = [
          <span key={uuid()} style={{ fontWeight: 'normal', display: 'block' }}>
            {t('集群外通过负载均衡域名或IP+服务端口访问服务')}
          </span>,
          <span key={uuid()} style={{ fontWeight: 'normal', display: 'block' }}>
            {t('集群内通过服务名+服务端口访问服务')}
          </span>
        ];
        break;

      case 'NodePort':
        content = (
          <span style={{ fontWeight: 'normal', display: 'block' }}>
            {t('可通过云服务器IP+主机端口访问服务，主机端口不填默认自动分配')}
          </span>
        );
        break;

      default:
        content = (
          <span style={{ fontWeight: 'normal', display: 'block' }}>{t('可通过服务名+服务端口可直接访问该服务')}</span>
        );
        break;
    }

    return content;
  }

  /** 展示端口映射，主机端口访问  多 一个 主机端口栏 */
  private _renderPortMap() {
    let {
      portsMap,
      communicationType,
      deletePortMap,
      updatePortMap,
      validatePortProtocol,
      validateTargetPort,
      validateServicePort,
      validateNodePort
    } = this.props;

    /**
     * 渲染协议列表
     */
    let protocolOptions = ProtocolList.map(item => (
      <option key={uuid()} value={item.value}>
        {item.label}
      </option>
    ));

    let columns: TablePanelColumnProps<PortMap>[] = [
      {
        key: 'protocol',
        width: '15%',
        header: t('协议'),
        headerTips: t('使用公网/内网负载均衡时，TCP和UDP协议不能混合使用'),
        render: x => (
          <div className={x.v_protocol.status === 2 ? 'is-error' : ''}>
            <Bubble
              placement="bottom"
              content={x.v_protocol.status === 2 ? <p className="form-input-help">{x.v_protocol.message}</p> : null}
            >
              <select
                className="tc-15-select m"
                style={{ minWidth: '90px' }}
                value={x.protocol}
                onChange={e => updatePortMap({ protocol: e.target.value }, x.id + '')}
                onBlur={e => validatePortProtocol(e.target.value, x.id + '')}
              >
                {protocolOptions}
              </select>
            </Bubble>
          </div>
        )
      },
      {
        key: 'targetPort',
        header: t('容器端口'),
        headerTips: t('端口范围1~65535'),
        render: x => (
          <div className={x.v_targetPort.status === 2 ? 'is-error' : ''}>
            <Bubble placement="bottom" content={x.v_targetPort.status === 2 ? x.v_targetPort.message : null}>
              <input
                type="text"
                className="tc-15-input-text m"
                placeholder={t('容器内应用程序监听的端口')}
                value={x.targetPort + ''}
                onChange={e => updatePortMap({ targetPort: e.target.value }, x.id + '')}
                onBlur={e => validateTargetPort(e.target.value, x.id + '')}
              />
            </Bubble>
          </div>
        )
      },
      {
        key: 'nodePort',
        header: t('主机端口'),
        headerTips: t('端口范围30000~32767，不填自动分配'),
        render: x => (
          <div className={x.v_nodePort.status === 2 ? 'is-error' : ''}>
            <Bubble placement="bottom" content={x.v_nodePort.status === 2 ? x.v_nodePort.message : null}>
              <input
                type="text"
                className="tc-15-input-text m"
                placeholder={t('范围：30000~32767')}
                value={x.nodePort + ''}
                onChange={e => updatePortMap({ nodePort: e.target.value }, x.id + '')}
                onBlur={e => validateNodePort(e.target.value, x.id + '')}
              />
            </Bubble>
          </div>
        )
      },
      {
        key: 'port',
        header: t('服务端口'),
        headerTips: this._renderPortHeader(),
        render: x => (
          <div className={x.v_port.status === 2 ? 'is-error' : ''}>
            <Bubble placement="bottom" content={x.v_port.status === 2 ? x.v_port.message : null}>
              <input
                type="text"
                className="tc-15-input-text m"
                placeholder={t('建议与容器端口一致')}
                value={x.port + ''}
                onChange={e => updatePortMap({ port: e.target.value }, x.id + '')}
                onBlur={e => validateServicePort(e.target.value, x.id + '')}
              />
            </Bubble>
          </div>
        )
      },
      {
        key: '',
        header: '',
        width: '10%',
        render: x => (
          <div>
            {portsMap.length > 1 ? (
              <a
                href="javascript:;"
                onClick={() => {
                  deletePortMap(x.id + '');
                }}
              >
                <i className="icon-cancel-icon" />
              </a>
            ) : (
              <Bubble placement="bottom" content={t('不可删除，至少指定一个端口映射')}>
                <a href="javascript:;" className="disabled">
                  <i className="icon-cancel-icon" />
                </a>
              </Bubble>
            )}
          </div>
        )
      }
    ];

    // nodePort只有在 loadBalancer 和 NodePort 类型下才有
    if (communicationType !== 'NodePort') {
      let columnIndex = columns.findIndex(item => item.key === 'nodePort');
      columns.splice(columnIndex, 1);
    }

    return (
      <Table
        columns={columns}
        records={portsMap}
        addons={[
          stylize({
            style: { overflow: 'visible', maxWidth: communicationType !== 'NodePort' ? '680px' : '850px' }
          })
        ]}
      />
    );
  }
}
