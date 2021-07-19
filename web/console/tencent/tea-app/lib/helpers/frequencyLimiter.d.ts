export declare class FrequencyLimiter {
    private name;
    private timer;
    private times;
    private timeout;
    private limitTimes;
    private tenseLimitTimes;
    private tense;
    private hasOutput;
    constructor({ name, limitTimes, tenseLimitTimes, timeout, }: {
        name?: string;
        limitTimes?: number;
        tenseLimitTimes?: number;
        timeout?: number;
    });
    exec(fn: Function): void;
}
