import { Router } from '../../../helpers/Router';

export const router = new Router('/tkestack/notify(/:mode)(/:resourceName)(/:tab)', {
  mode: '',
  resourceName: 'channel',
  tab: ''
});
