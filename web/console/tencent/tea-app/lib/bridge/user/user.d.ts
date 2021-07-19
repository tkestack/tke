import { AppUserData } from "./current";
import { PermitedProjectInfo, ProjectItem } from "./project";
import { UserEmitter } from "./event";
export declare const user: AppUserBridge;
/**
 * `app.user` 导出的接口
 */
export interface AppUserBridge extends UserEmitter {
    /**
     * 获取当前登录用户的信息
     *
     * @example
      ```js
      const user = await app.user.current();
      ```
     */
    current(): Promise<AppUserData>;
    /**
     * 获取当前登录用户的反 CSRF 凭据，
     * 用于同域网络请求中防御跨站脚本攻击
     */
    getAntiCSRFToken(): string;
    /**
     * 检查用户是否在指定的白名单中
     *
     * @returns 如在白名单内，返回当前 ownerUin，否则返回 0
     *
     * **注意**：这个方法不会缓存白名单查询结果，如果不想重复查询，业务请自行缓存
     *
     * @example
     ```js
      // 如在白名单内，返回当前 ownerUin，否则返回 0
      const ownerUin = await app.user.checkWhitelist('CLB_NEW_CONSOLE');
      if (ownerUin) {
        // ...白名单操作
      }
      ```
     */
    checkWhitelist(key: string): Promise<number>;
    /**
     * 批量检查用户是否在白名单中
     *
     * @param keys 要检查的白名单键值
     *
     * @returns 返回一个对象，其 key 值为传入的白名单 key 值。对于在白名单的用户，对应的值为用户 ownerUin，否则为 0
     *
     * **注意**：这个方法不会缓存白名单查询结果，如果不想重复查询，业务请自行缓存
     */
    checkWhitelistBatch(keys: string[]): Promise<{
        [key: string]: number;
    }>;
    /**
     * 获取用户最后使用的 regionId
     *
     * - 如果不存在，则返回 `-1`
     * - 该数据基于用户当前的 ownerUin 存储在 localStorage 中
     * - 基于上一点，该数据可能会被其它业务修改，所以使用前应该先校验合法性
     */
    getLastRegionId(): number;
    /**
     * 设置用户最后一次访问的 regionId
     *
     * 该数据可以被其它业务使用
     */
    setLastRegionId(regionId: number): void;
    /**
     * 获取用户最后使用的 projectId
     *
     * - 如果不存在，则返回 `-1`
     * - 该数据基于用户当前的 ownerUin 存储在 localStorage 中
     * - 基于上一点，该数据可能会被其它业务修改，所以使用前应该先校验合法性
     */
    getLastProjectId(): number;
    /**
     * 设置用户最后一次访问的 projectId
     *
     * 该数据可以被其它业务使用
     */
    setLastProjectId(projectId: number): void;
    /**
     * 获取用户有权限的项目信息
     */
    getPermitedProjectInfo(): Promise<PermitedProjectInfo>;
    /**
     * 获取用户有权限的项目列表
     */
    getPermitedProjectList(): Promise<ProjectItem[]>;
    /**
     * 清除用户当前登录态，并弹出登录对话框
     */
    login(): void;
    /**
     * 清除用户当前登录态
     */
    logout(): void;
}
export { AppUserData, AppUserIdentityInfo } from "./current";
export { PermitedProjectInfo, ProjectItem } from "./project";
