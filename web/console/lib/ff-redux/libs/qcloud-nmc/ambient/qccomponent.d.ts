declare namespace nmc {
  export interface BeeInstance {
    $set(props: any): void;
    $set(path: string, props: any): void;
    $get(path: string): void;
    $destroy(): void;
  }

  export interface Bee {
    mount(id: string | HTMLElement, props?: any): BeeInstance;
    utils: {
      ie: number;
    };
  }
  export namespace Bee {}
}
