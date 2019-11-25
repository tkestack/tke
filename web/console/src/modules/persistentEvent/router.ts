import { Router } from '../../../helpers/Router';

/**
 * @param mode  当前的展示内容类型 list | create | update | detail
 * @param type  三级菜单栏所对应的资源 resource | service ……
 * @param resourceName  资源的名称  deployment 等
 * @param tab   tab页面
 */
export const router = new Router('/tkestack/persistent-event(/:mode)', { mode: '' });
