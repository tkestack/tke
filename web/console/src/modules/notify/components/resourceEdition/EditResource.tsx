import * as React from 'react';
import { uuid } from '@tencent/qcloud-lib';
import { t } from '@tencent/tea-app/lib/i18n';
import { Card } from '@tea/component/card';
import { Form } from '@tea/component/form';
import { Input } from '@tea/component/input';
import { Button } from '@tea/component/button';
import { router } from '../../router';
import { channelSchema } from '../../schema/channelSchema';
import { getState, onChange, schemaObjToJSON } from '../../schema/schemaUtil';
import { resourceConfig } from '../../../../../config';
import { Text, Switch, InputNumber } from '@tencent/tea-component';
import { RootProps } from '../NotifyApp';
import { OperationState } from '@tencent/qcloud-redux-workflow';
import { CreateResource, TipInfo, getWorkflowError } from '../../../common';
const rc = resourceConfig();

interface Props extends RootProps {
  instance?;
}
interface State {
  resource?;
  error?;
  searchValue?;
}

export class EditResource extends React.Component<Props, State> {
  constructor(props) {
    super(props);
    let schema = this.getSchema();
    let state = this.getState();
    state.resource = props.instance
      ? getState(schema, this, this.props.instance, this.props.instance)
      : getState(schema, this);
    this.state = state;
  }

  componentDidMount() {
    this.props.actions.workflow.modifyResource.reset();
  }

  getState() {
    return {
      resource: undefined,
      error: undefined
    };
  }

  getSchema(): any {
    return channelSchema;
  }

  render() {
    let { modifyResourceFlow } = this.props;
    return (
      <Card>
        <Card.Body>
          {this.renderForm()}
          <div className="tea-form-operate">
            <Button type="primary" onClick={this.submit}>
              {t('保存')}
            </Button>
            <Button type="weak" onClick={this.back}>
              {t('取消')}
            </Button>
            <TipInfo
              className="tea-mt-1n"
              style={{ marginBottom: '0' }}
              isShow={modifyResourceFlow.operationState === OperationState.Done && modifyResourceFlow.results[0].error}
              type="error"
            >
              {getWorkflowError(modifyResourceFlow)}
            </TipInfo>
          </div>
        </Card.Body>
      </Card>
    );
  }

  renderForm() {
    return <React.Fragment />;
  }

  renderFields(obj) {
    return Object.keys(obj.properties).map(key => (
      <Form.Item label={key} key={key} required={obj.properties[key].required}>
        {obj.properties[key].type === 'boolean' ? (
          <Switch value={obj.properties[key].value} onChange={onChange(obj.properties[key])} />
        ) : obj.properties[key].type === 'number' ? (
          <InputNumber value={obj.properties[key].value} onChange={onChange(obj.properties[key])} />
        ) : (
          <Input
            type={key === 'password' ? 'password' : undefined}
            multiline={key === 'body'}
            value={obj.properties[key].value}
            onChange={onChange(obj.properties[key])}
          />
        )}
      </Form.Item>
    ));
  }

  submit = () => {
    let { actions, route } = this.props;
    let urlParams = router.resolve(route);

    let resourceInfo = rc[urlParams.resourceName] || rc.channel;
    let mode = this.props.instance ? 'modify' : 'create';
    let json = schemaObjToJSON(this.state.resource);
    if (this.props.instance) {
      json.metadata.resourceVersion = this.props.instance.metadata.resourceVersion;
    }
    let jsonData = JSON.stringify(json);

    let resource: CreateResource = {
      id: uuid(),
      resourceInfo,
      mode,
      namespace: (json.metadata && json.metadata.namespace) || route.queries['np'] || 'default',
      jsonData,
      clusterId: route.queries['clusterId'],
      resourceIns: mode === 'modify' ? this.props.instance.metadata.name : ''
    };

    actions.workflow.modifyResource.start([resource]);
    actions.workflow.modifyResource.perform();
  };
  back = () => {
    let { route } = this.props;
    let urlParams = router.resolve(route);
    router.navigate(
      Object.assign({}, urlParams, {
        mode: urlParams.mode === 'update' ? 'detail' : 'list'
      }),
      Object.assign({}, route.queries)
    );
  };
}
