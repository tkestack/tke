import { Router } from '../../../helpers/Router';

/**
 * @param sub 二级导航，eg: create、update等
 * @param tab 详情等tab子页
 * 业务侧和平台侧使用不同的URL 路径
 */
const baseURI = window.location.href.includes('/tkestack-project') ? '/tkestack-project/application/list/k8sLog/logagent' : '/tkestack/log';
export const router = new Router(`${baseURI}(/:mode)(/:tab)`, { mode: '', tab: '' });
