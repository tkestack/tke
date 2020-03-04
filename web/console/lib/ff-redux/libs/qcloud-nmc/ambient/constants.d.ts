declare namespace nmc {
  interface Constants {
    /** RegionID => RegionKey */
    REGIONMAP: { [id: number]: string };

    /** RegionID => RegionName */
    REGIONNAMES: { [id: number]: string };

    /** 地域顺序 */
    REGIONORDER: number[];
  }
  const Constants: Constants;
}
