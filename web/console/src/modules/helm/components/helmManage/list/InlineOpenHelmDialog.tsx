import * as React from 'react';
import { RootProps } from '../../HelmApp';
import { Button, Modal } from '@tea/component';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
export class InlineOpenHelmDialog extends React.Component<RootProps, {}> {
  render() {
    let { actions } = this.props;
    const cancel = () => {};

    const submit = () => {};
    return (
      <Modal visible={true} caption={t('确认开通Helm应用管理？')} onClose={cancel} disableEscape={true}>
        <Modal.Body>
          <div style={{ fontSize: '14px', lineHeight: '20px' }}>
            <div className="docker-dialog jiqun">
              <p>
                {t(
                  '开通Helm应用，将在集群内安装Helm tiller组件，占用集群计算资源，目前您所选的集群尚未开通，是否立即开通？'
                )}
              </p>
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button type="primary" onClick={submit}>
            {t('确定')}
          </Button>
          <Button onClick={cancel}>{t('取消')}</Button>
        </Modal.Footer>
      </Modal>
    );
  }
}
