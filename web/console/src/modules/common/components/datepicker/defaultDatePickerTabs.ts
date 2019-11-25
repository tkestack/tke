import { DatePickerProps } from './';
import { t, Trans } from '@tencent/tea-app/lib/i18n';

export const defaultDatePickerOptions: DatePickerProps = {
  tabs: [
    { from: '%TODAY', to: '%TODAY', label: t('今天') },
    { from: '%TODAY-1', to: '%TODAY-1', label: t('昨天') },
    { from: '%TODAY-7', to: '%TODAY', label: t('近7天') },
    { from: '%TODAY-30', to: '%TODAY', label: t('近30天') }
  ],
  range: {
    min: '%TODAY-90',
    max: '%TODAY'
  }
};
