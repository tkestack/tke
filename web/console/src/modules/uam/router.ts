import { Router } from '../../../helpers/Router';

/**
 * @param module   二级菜单栏的一级导航
 */
export const router = new Router('/tkestack/uam(/:module)(/:sub)', { module: '', sub: '' });
