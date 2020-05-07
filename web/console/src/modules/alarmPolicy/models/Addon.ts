export interface AddonStatus {
  [propName: string]: Addon;
}

export interface Addon {
  status: string;
  name: string;
}
