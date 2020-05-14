import { Router } from '../../../helpers/Router';

/**
 * @param sub 二级导航，eg: create、update等
 * @param tab 详情等tab子页
 */
export const router = new Router('/tkestack/audit(/:sub)(/:tab)', { sub: '', tab: '' });
