import * as React from 'react';
import { Table, TableColumn, Text } from '@tea/component';
import { LinkButton } from '../../../../common/components';
import { InstallingHelm } from '../../../models';
import { RootProps } from '../../HelmApp';
import { InstallingStatus, InstallingStatusText } from '../../../constants/Config';
import classNames from 'classnames';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

import { stylize } from '@tea/component/table/addons/stylize';
import { selectable } from '@tea/component/table/addons/selectable';

export class InstallingHelmContent extends React.Component<RootProps, {}> {
  render() {
    const {
      listState: { installingHelmDetail }
    } = this.props;
    return (
      <div>
        <div className="configuration-box" style={{ overflow: 'hidden' }}>
          <div className="version-wrap">{this.renderInstallingHelmTable()}</div>
          <div className="rich-textarea simple-mod">
            <pre
              style={{
                backgroundColor: 'black',
                color: 'white',
                width: '100%',
                height: '100%',
                whiteSpace: 'pre-wrap',
                wordWrap: 'break-word',
                margin: 0
              }}
            >
              {installingHelmDetail ? installingHelmDetail.message : ''}
            </pre>
          </div>
        </div>
      </div>
    );
  }
  delete(helm: InstallingHelm) {
    this.props.actions.helm.ignoreInstallingHelm(helm);
  }
  private renderInstallingHelmTable() {
    let {
      actions,
      listState: { installingHelmList, installingHelmSelection }
    } = this.props;

    const colunms: TableColumn<InstallingHelm>[] = [
      {
        key: 'helmName',
        header: t('Helm名称'),
        width: '40%',
        render: x => {
          return (
            <Text parent="div" overflow>
              {x.name}
            </Text>
          );
        }
      },
      {
        key: 'status',
        header: t('状态'),
        width: '25%',
        render: x => {
          return (
            <div>
              <p
                className={classNames(
                  'text-overflow',
                  InstallingStatusText[x.status] && InstallingStatusText[x.status].classname
                )}
              >
                {InstallingStatusText[x.status] ? InstallingStatusText[x.status].text : '-'}
              </p>
            </div>
          );
        }
      },
      {
        key: 'op',
        header: t('操作'),
        width: '35%',
        render: x => this._renderOperationCell(x)
      }
    ];

    return (
      <Table
        columns={colunms}
        records={installingHelmList.data.records}
        recordKey="id"
        addons={[
          stylize({
            className: 'version-list update-cont hidden-checkbox'
          }),
          selectable({
            // targetColumnKey: 'helmName',
            value: installingHelmSelection ? [installingHelmSelection.id as string] : [],
            rowSelect: true,
            onChange: keys => {
              let record = null;
              if (keys.length === 1) {
                record = installingHelmList.data.records.find(record => record.id === keys[0]);
              } else if (keys.length === 2) {
                let selectKey = keys.find(key => key !== installingHelmSelection.id);
                record = installingHelmList.data.records.find(record => record.id === selectKey);
              }

              if (record) {
                actions.helm.selectInstallingHelm(null);
                actions.helm.selectInstallingHelm(record);
              }
            }
          })
        ]}
      />
    );
  }

  private _renderOperationCell(helm: InstallingHelm) {
    //1：失败；0:成功
    return (
      <div>
        <LinkButton
          isShow={helm.status === InstallingStatus.ERROR}
          onClick={e => {
            e.stopPropagation();
            this.delete(helm);
          }}
        >
          {t('取消安装')}
        </LinkButton>
      </div>
    );
  }
}
