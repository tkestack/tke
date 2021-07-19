export declare function getLastProjectId(): number;
export declare function setLastProjectId(projectId: number): void;
export declare function clearLastProjectId(): void;
export declare function getPermitedProjectInfo(): Promise<PermitedProjectInfo>;
export declare const getPermitedProjectList: () => Promise<ProjectItem[]>;
export interface PermitedProjectInfo {
    /**
     * 当前用户是否有具备查看所有项目的权限
     */
    isShowAll: boolean;
    /**
     * 当前用户有权限的项目列表
     */
    projects: ProjectItem[];
}
export interface ProjectItem {
    /**
     * 项目 ID，为 0 表示默认项目
     */
    projectId: number;
    /**
     * 项目名称
     */
    projectName: string;
}
