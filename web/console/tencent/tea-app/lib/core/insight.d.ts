export declare const internalInsightShouldNotUsedByBusiness: InsightCore;
export interface InsightCore {
    EventLevel: {
        Info: 0;
        Warn: 1;
        Error: 2;
    };
    register(name: string, meta: InsightEventMeta): InsightEventBus;
    care<T extends (...args: any) => any>(method: T, options: InsightCareOptions): T;
    watch<T extends Function>(method: T, bus: InsightEventBus, builder: Function): T;
    captureTrace(target: any, defaultMessage?: string, defaultStack?: string): InsightTrace;
    shutdown(): void;
    on(event: string, callback: Function): void;
    off(event: string, callback: Function): void;
    emit(event: string, ...args: any[]): void;
    jsErrorBus: InsightEventBus;
    promiseErrorBus: InsightEventBus;
    detectNetwork(detectList: InsightNetworkDetectTask[]): void;
    stat(name: string, data: any, urgent?: boolean): void;
}
declare type InsightNetworkDetectTask = string | {
    name: string;
    url: string;
    alias: string;
};
interface InsightEventMeta {
    level: InsightEventLevel;
}
declare type InsightEventLevel = 0 | 1 | 2;
interface InsightEventBus {
    name: string;
    push: (event: any, urgent?: boolean) => void;
}
interface InsightTrace {
    message: string;
    stack: string;
}
interface Function {
    __insight_internal__?: boolean;
    __insight_careby__?: Function;
    __insight_careto__?: Function;
    __insight_fillto__?: Function;
    __insight_fillby__?: Function;
}
interface InsightCareOptions {
    capture: InsightCareCapture;
    before?: Function;
    after?: Function;
}
interface InsightCareCapture {
    (trace: InsightTrace, method: Function, error: Error): void;
}
export {};
