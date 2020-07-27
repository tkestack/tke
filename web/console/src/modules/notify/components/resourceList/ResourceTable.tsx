import * as React from 'react';
import { Button, Text } from '@tea/component';
import { RootProps } from '../NotifyApp';
import { LinkButton } from '../../../common/components';
import { router } from '../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Resource } from '../../../common';
import { resourceConfig } from '../../../../../config';
import { dateFormatter } from '../../../../../helpers';
import { TablePanelColumnProps, TablePanel } from '@tencent/ff-component';

const rc = resourceConfig();

interface Props extends RootProps {
  onlyTable?: boolean;
  bordered?: boolean;
  resourceName?: string;
}

export class ResourceTable extends React.Component<Props, {}> {
  render() {
    return this._renderTablePanel();
  }

  getColumns(): TablePanelColumnProps<Resource>[] {
    return [];
  }

  protected renderNameColumn(id, displayName, resourceName?) {
    const { route } = this.props;

    let urlParams = router.resolve(route);
    return (
      <React.Fragment>
        <Text parent="div" className="m-width" overflow>
          <LinkButton
            title={id}
            onClick={() => {
              router.navigate(
                {
                  ...urlParams,
                  mode: 'detail',
                  resourceName: resourceName || this.props.resourceName || urlParams.resourceName || 'channel'
                },
                { ...route.queries, resourceIns: id }
              );
            }}
            className="tea-text-overflow"
          >
            {id || '-'}
          </LinkButton>
        </Text>
        {displayName && <Text parent="div">{displayName}</Text>}
      </React.Fragment>
    );
  }

  private _renderTablePanel() {
    const { actions, route } = this.props;

    let urlParams = router.resolve(route);

    let resource = this.props[this.props.resourceName || urlParams['resourceName']] || this.props.channel;
    const columns: TablePanelColumnProps<Resource>[] = [
      {
        key: 'name',
        header: t('名称'),
        render: x => this.renderNameColumn(x.metadata.name, x.spec.displayName)
      },
      ...this.getColumns(),
      {
        key: 'creationTimestamp',
        header: t('创建时间'),
        width: '15%',
        render: x => {
          let time: any = '-';
          if (x.metadata.creationTimestamp) {
            time = dateFormatter(new Date(x.metadata.creationTimestamp), 'YYYY-MM-DD HH:mm:ss');
          }
          return <Text>{time}</Text>;
        }
      }
    ];

    let emptyTips: JSX.Element = (
      <div className="text-center">
        <Trans>
          {(rc[urlParams.resourceName] || rc.channel).headTitle}为空，您可以
          <Button
            type="link"
            onClick={() => {
              this._handleCreate();
            }}
          >
            新建{(rc[urlParams.resourceName] || rc.channel).headTitle}
          </Button>
        </Trans>
      </div>
    );

    return (
      <TablePanel
        left={
          !this.props.onlyTable && (
            <React.Fragment>
              <Button type="primary" onClick={() => this.handleCreate()}>
                {t('新建')}
              </Button>
              <Button
                disabled={resource.selections.length === 0}
                onClick={() => {
                  let { route } = this.props;
                  let urlParams = router.resolve(route);
                  resource.selections.forEach(resource => {
                    resource.resourceInfo = rc[urlParams.resourceName] || rc.channel;
                    resource.isSpetialNamespace = true;
                  });
                  actions.workflow.deleteResource.start(resource.selections);
                }}
              >
                {t('删除')}
              </Button>
            </React.Fragment>
          )
        }
        isNeedCard={!this.props.onlyTable}
        bordered={this.props.bordered}
        columns={columns}
        emptyTips={emptyTips}
        model={resource}
        action={actions.resource[urlParams['resourceName']]}
        getOperations={!this.props.onlyTable && (record => this.getOperations(record))}
        selectable={
          !this.props.onlyTable && {
            value: resource.selections.map(item => item.id as string),
            onChange: keys => {
              actions.resource[urlParams['resourceName']].selects(
                resource.list.data.records.filter(item => keys.indexOf(item.id as string) !== -1)
              );
            }
          }
        }
        isNeedPagination={false}
      />
    );
  }

  private handleCreate() {
    let { route } = this.props;
    let urlParams = router.resolve(route);
    let rid = route.queries['rid'];
    router.navigate({ ...urlParams, mode: 'create' }, { rid });
  }

  private _handleDeleteResource(resource: Resource) {
    let { actions, route } = this.props;
    let urlParams = router.resolve(route);
    resource.resourceInfo = rc[urlParams.resourceName] || rc.channel;
    resource.isSpetialNamespace = true;
    actions.resource[urlParams['resourceName']].selects([resource]);
    actions.workflow.deleteResource.start([resource]);
  }

  private _handleCreate() {
    let { route } = this.props;
    let urlParams = router.resolve(route);
    router.navigate({ ...urlParams, mode: 'create' });
  }
  private getOperations(resource: Resource) {
    const renderDeleteButton = () => {
      return (
        <LinkButton key={'delete'} onClick={() => this._handleDeleteResource(resource)}>
          {t('删除')}
        </LinkButton>
      );
    };

    return [renderDeleteButton()];
  }
}
