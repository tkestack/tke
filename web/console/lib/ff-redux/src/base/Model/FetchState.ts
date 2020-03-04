export enum FetchState {
  /** indicates the data is up to date and ready to use */
  Ready = 'Ready',

  /** indicates the data is out of date, and the new data is fetching */
  Fetching = 'Fetching',

  /**
   * indicates the data is out of date, and the new data fetches failed
   */
  Failed = 'Failed'
}
