import * as EventEmitter from "eventemitter3";
interface InvalidateEventArgs {
    source: "accountChanged" | "logout";
}
interface UserEventTypes {
    invalidate: [InvalidateEventArgs];
}
export declare class UserEmitter extends EventEmitter<UserEventTypes> {
}
export declare const userEmitter: UserEmitter;
export {};
