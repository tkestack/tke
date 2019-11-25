import * as React from 'react';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { onChange } from '../../schema/schemaUtil';
import { Radio, SearchBox, Table, Transfer, Text } from '@tencent/tea-component';
import { EditResource } from './EditResource';
import { receiverGroupSchema } from '../../schema/receiverGroupSchema';
import { scrollable, selectable, removeable } from '@tencent/tea-component/lib/table/addons';
import { LinkButton } from '../../../common';
import { router } from '../../router';
export class EditResourceReceiverGroup extends EditResource {
  getState() {
    return {
      ...super.getState(),
      searchValue: ''
    };
  }
  getSchema() {
    return receiverGroupSchema;
  }

  componentDidMount() {
    super.componentDidMount();
    this.props.actions.resource.receiver.fetch();
  }

  renderForm() {
    let resource = receiverGroupSchema;
    resource = this.state.resource;
    const columns = [
      {
        key: 'name',
        header: t('名称'),
        render: x => {
          return (
            <React.Fragment>
              <Text parent="div" className="m-width" overflow>
                <LinkButton
                  title={x.metadata.name}
                  onClick={() => {
                    router.navigate(
                      { mode: 'detail', resourceName: 'receiver' },
                      { ...this.props.route.queries, resourceIns: x.metadata.name }
                    );
                  }}
                  className="tea-text-overflow"
                >
                  {x.metadata.name || '-'}
                </LinkButton>
              </Text>
              <Text parent="div">{x.spec.displayName || '-'}</Text>
            </React.Fragment>
          );
        }
      }
    ];
    return (
      <Form>
        <Form.Item label={t('名称')} required>
          <Input
            placeholder={t('请填写名称')}
            value={resource.properties.spec.properties.displayName.value}
            onChange={onChange(resource.properties.spec.properties.displayName)}
          />
        </Form.Item>
        <Form.Item label={t('接收人')} required>
          <Transfer
            leftCell={
              <Transfer.Cell
                scrollable={false}
                title={t('选择接收人')}
                tip={t('支持按住 shift 键进行多选')}
                header={
                  <SearchBox
                    placeholder={t('搜索')}
                    value={this.state.searchValue}
                    onChange={searchValue => this.setState({ searchValue })}
                  />
                }
              >
                <Table
                  recordKey={x => x.metadata.name}
                  records={this.props.receiver.list.data.records.filter(
                    item =>
                      !this.state.searchValue ||
                      item.spec.displayName.includes(this.state.searchValue) ||
                      item.metadata.name.includes(this.state.searchValue)
                  )}
                  columns={columns}
                  addons={[
                    scrollable({
                      maxHeight: 310
                    }),
                    selectable({
                      value: resource.properties.spec.properties.receivers.value,
                      onChange: onChange(resource.properties.spec.properties.receivers),
                      rowSelect: true
                    })
                  ]}
                />
              </Transfer.Cell>
            }
            rightCell={
              <Transfer.Cell title={`已选择 (${resource.properties.spec.properties.receivers.value.length})`}>
                <Table
                  recordKey={x => x.metadata.name}
                  records={this.props.receiver.list.data.records.filter(i =>
                    resource.properties.spec.properties.receivers.value.includes(i.metadata.name)
                  )}
                  columns={columns}
                  addons={[
                    removeable({
                      onRemove: key => {
                        resource.properties.spec.properties.receivers.value = resource.properties.spec.properties.receivers.value.filter(
                          i => i !== key
                        );
                        this.setState({});
                      }
                    })
                  ]}
                />
              </Transfer.Cell>
            }
          />
        </Form.Item>
      </Form>
    );
  }
}
