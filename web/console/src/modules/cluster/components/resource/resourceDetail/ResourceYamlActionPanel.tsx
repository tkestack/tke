import * as React from 'react';
import { Button, Bubble, Justify, Table } from '@tea/component';
import { bindActionCreators } from '@tencent/qcloud-lib';
import { connect } from 'react-redux';
import { allActions } from '../../../actions';
import { RootProps } from '../../ClusterApp';
import { router } from '../../../router';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

const mapDispatchToProps = dispatch =>
  Object.assign({}, bindActionCreators({ actions: allActions }, dispatch), { dispatch });

@connect(
  state => state,
  mapDispatchToProps
)
export class ResourceYamlActionPanel extends React.Component<RootProps, {}> {
  render() {
    let { route } = this.props,
      urlParams = router.resolve(route);

    let disableBtn = false,
      errorTip = '';

    let ns = route.queries['np'],
      resourceIns = route.queries['resourceIns'];
    if (ns === 'kube-system') {
      disableBtn = true;
      errorTip = t('当前命名空间下的资源不可编辑');
    } else if (urlParams['resourceName'] === 'svc' && ns !== 'kube-system') {
      disableBtn = resourceIns === 'kubernetes';
      errorTip = t('系统默认的Service不可编辑');
    }

    return (
      <Table.ActionPanel>
        <Justify
          left={
            <Bubble placement="left" content={disableBtn ? errorTip : null}>
              <Button
                type="primary"
                disabled={disableBtn}
                onClick={() => {
                  if (!disableBtn) {
                    router.navigate(
                      Object.assign({}, urlParams, { mode: 'modify' }),
                      Object.assign({}, route.queries, { resourceIns })
                    );
                  }
                }}
              >
                {t('编辑YAML')}
              </Button>
            </Bubble>
          }
        />
      </Table.ActionPanel>
    );
  }
}
