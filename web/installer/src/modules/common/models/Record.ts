export interface Record<T> {
    record: T;
    auth?: {
        isAuthorized?: boolean;
        isLoginedSec?: boolean;
        message?: string;
        redirect?: string;
    }
};