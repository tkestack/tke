const manager = seajs.require('manager');

type Callback<T> = (data: T) => any;

function promisify<T>(receive: (cb: Callback<T>) => any): Promise<T> {
  return new Promise<T>(resolve => receive(resolve));
}

export function queryWhitelist(whiteKey: string[]) {
  return new Promise<nmc.WhitelistMap>(resolve => manager.queryWhiteList({ whiteKey }, resolve));
}

export async function isInWhitelist(key: string) {
  const commonData = await promisify(manager.getComData.bind(manager) as typeof manager.getComData);
  if (!commonData.userInfo) return false;
  const uin = String(commonData.userInfo.ownerUin);

  const cached = await promisify(manager.getAllWhiteList.bind(manager) as typeof manager.getAllWhiteList);
  if (cached[key] && cached[key].indexOf(uin) > -1) {
    return true;
  }

  const query = await queryWhitelist([key]);
  return query[key] && query[key].indexOf(uin) > -1;
}
