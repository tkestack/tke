import * as React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Button, Modal, Form, Input, Table } from '@tea/component';
import { bindActionCreators, uuid } from '@tencent/ff-redux';
import { t, Trans } from '@tencent/tea-app/lib/i18n';
import { useForm, useField } from 'react-final-form-hooks';
import { allActions } from '../../../actions';
import { getStatus } from '../../../../common/validate';
import { VALIDATE_PASSWORD_RULE } from '@src/modules/uam/constants/Config';

export function PasswordModifyDialog(props) {
  const state = useSelector((state) => state);
  const dispatch = useDispatch();
  const { actions } = bindActionCreators({ actions: allActions }, dispatch);

  const { isShowing, toggle, user } = props;

  function onSubmit(values, form) {
    console.log('PasswordModifyDialog submit values:', values);
    const { password } = values;
    const extraObj = { ...user.spec.extra, policies: Object.keys(JSON.parse(user.spec.extra.policies)).join(',') };
    actions.user.updateUser.fetch({
      noCache: true,
      data: {
        user: {
          metadata: {
            name: user.metadata.name,
            resourceVersion: user.metadata.resourceVersion,
          },
          spec: { ...user.spec, hashedPassword: btoa(password), extra: extraObj },
        },
      },
    });
    setTimeout(form.reset);
    toggle();
  }
  const { form, handleSubmit, validating, submitting } = useForm({
    onSubmit,
    /**
     * 默认为 shallowEqual
     * 如果初始值有多层，会导致重渲染，也可以使用 `useEffect` 设置初始值：
     * useEffect(() => form.initialize({ }), []);
     */
    initialValuesEqual: () => true,
    initialValues: { password: '', rePassword: '' },
    validate: ({ password, rePassword }) => {
      const errors = {
        password: undefined,
        rePassword: undefined,
      };
      if (!password) {
        errors.password = t('请输入密码');
      } else if (!VALIDATE_PASSWORD_RULE.pattern.test(password)) {
        errors.password = VALIDATE_PASSWORD_RULE.message;
      }

      if (!rePassword) {
        errors.rePassword = t('请再次输入密码');
      } else if (!VALIDATE_PASSWORD_RULE.pattern.test(rePassword)) {
        errors.rePassword = VALIDATE_PASSWORD_RULE.message;
      } else if (password !== rePassword) {
        errors.rePassword = t('两次输入密码需一致');
      }
      return errors;
    },
  });
  const password = useField('password', form);
  const rePassword = useField('rePassword', form);

  return (
    <Modal
      visible={isShowing}
      caption={t('修改密码')}
      onClose={() => {
        toggle();
        setTimeout(form.reset);
      }}
    >
      <Modal.Body>
        <form onSubmit={handleSubmit}>
          <Form>
            <Form.Item
              label={t('用户密码')}
              required
              status={getStatus(password.meta, validating)}
              message={
                getStatus(password.meta, validating) === 'error' ? password.meta.error : VALIDATE_PASSWORD_RULE.message
              }
            >
              <Input {...password.input} type="password" size="l" autoComplete="off" placeholder={t('请输入密码')} />
            </Form.Item>
            <Form.Item
              label={t('确认密码')}
              required
              status={getStatus(rePassword.meta, validating)}
              message={
                getStatus(rePassword.meta, validating) === 'error'
                  ? rePassword.meta.error
                  : VALIDATE_PASSWORD_RULE.message
              }
            >
              <Input
                {...rePassword.input}
                type="password"
                size="l"
                autoComplete="off"
                placeholder={t('请再次输入密码')}
              />
            </Form.Item>
          </Form>
          <Form.Action style={{ textAlign: 'center' }}>
            <Button type="primary" htmlType="submit" loading={submitting} disabled={validating}>
              <Trans>确定</Trans>
            </Button>
            <Button
              type="weak"
              htmlType="reset"
              onClick={() => {
                toggle();
                setTimeout(form.reset);
              }}
            >
              <Trans>取消</Trans>
            </Button>
          </Form.Action>
        </form>
      </Modal.Body>
    </Modal>
  );
}
