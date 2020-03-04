import { findByCondition } from '../../qcloud-lib';

const manager = seajs.require('manager');
const appUtil = seajs.require('appUtil');
const constants = seajs.require('config/constants');

const { REGIONORDER, REGIONMAP } = constants;

export interface Region {
  /** 区域 ID，如广州为 1 */
  id: number;

  /** 区域 key，如广州为 gz */
  key: string;

  /** 区域名称，如广州为【华南地区（广州）】 */
  name: string;
}

let regionListPromise: Promise<Region[]> = null;

/**
 * 获取当前默认的地域ID（存于localStorage中，nmc全局属性）
 */
export function getRegionId(): number {
  return +appUtil.getRegionId();
}

/**
 * 设置为默认地域（存于localStorage中，nmc全局属性）
 */
export function setRegionId(regionId: nmc.RegionId) {
  appUtil.setRegionId(regionId);
}

/**
 * 获取区域列表
 */
export function fetchRegionList(): Promise<Region[]> {
  if (!regionListPromise) {
    regionListPromise = new Promise((resolve, reject) => {
      manager.queryRegion(map => {
        let list: Region[] = [];
        REGIONORDER.forEach(id => {
          if (map[id]) {
            const key = REGIONMAP[id];
            const name = map[id];
            list.push({ id, key, name });
          }
        });
        resolve(list);
      }, reject);
    });
  }
  return regionListPromise;
}

export async function findRegion(regionId: number | string) {
  const list = await fetchRegionList();
  const region = findByCondition(list, x => x.id === regionId || x.key === regionId);
  return region || null;
}
