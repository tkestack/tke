/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
