import { Router } from '../../../helpers/Router';

/**
 * @param mode 当前的模式，create | update | detail
 * @param type 当前addon的类型  logCollector等
 * @param tab tab页面
 */
export const router = new Router('/tkestack/addon(/:mode)(/:type)(/:tab)', {
  mode: '',
  type: '',
  tab: ''
});
