export enum FetcherTrigger {
  /**
   * trigger a load operation
   */
  Start = 'Start',

  /**
   * trigger when load for the tolerance duration
   * */
  Loading = 'Loading',

  /** trigger a receive operation */
  Done = 'Done',

  /** trigger a failed result */
  Fail = 'Fail',

  /** trigger a manual update */
  Update = 'Update',

  Clear = 'Clear'
}
