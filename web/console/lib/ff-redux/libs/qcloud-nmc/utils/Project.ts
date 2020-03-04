const manager = seajs.require('manager');
const appUtil = seajs.require('appUtil');

export interface Project {
  /** 项目 ID */
  id: number;

  /** 项目名称 */
  name: string;
}

let projectListPromise: Promise<Project[]> = null;

/**
 * 获取项目列表
 */
export async function fetchProjectList(): Promise<Project[]> {
  if (!projectListPromise) {
    projectListPromise = new Promise((resolve, reject) => {
      manager.getProjects(result => {
        const projectList: Project[] = result.data.map(x => ({
          id: x.projectId,
          name: x.name
        }));
        resolve(projectList);
      });
    });
  }
  return projectListPromise;
}

/**
 * 获取当前默认的项目ID（存于localStorage中，nmc全局属性）
 * @param skipAll 是否无视 “全部项目”（-1），默认为false
 */
export function getProjectId(skipAll?: boolean): number {
  return +appUtil.getProjectId(skipAll);
}

/**
 * 设置为当前的默认项目（存于localStorage中，nmc全局属性）
 */
export function setProjectId(projectId: nmc.ProjectId) {
  appUtil.setProjectId(projectId);
}
