import { Identifiable } from '@tencent/qcloud-lib';

export interface SubRouter extends Identifiable, BasicRouter {
  /** 是否有二级导航 */
  sub?: BasicRouter[];
}

export interface BasicRouter {
  /** 路由的名称 */
  name: string;

  /** 路由的路径 */
  path?: string;

  /** 非嵌套路由需要标识一个basic */
  basicUrl?: string;
}

export interface SubRouterFilter {
  /** 模块名称 */
  module: 'cluster' | 'mesh';
  /** sub 一级路由的名称 */
  sub: string;
}
