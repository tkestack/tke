///<reference path="nmc.d.ts" />

interface DefineInterface {
  (id: string, factory: (require: nmc.Require, exports, module) => any): any;
}

interface Window {
  define: DefineInterface;
}

declare namespace seajs {
  export var require: nmc.Require;
}
