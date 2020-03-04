/// <reference path="manager.d.ts" />
/// <reference path="constants.d.ts" />
declare namespace nmc {
  interface NetResponse {
    code: number;
    data: any;
  }

  interface NetSendConfig {
    /**
     * GET / POST
     * */
    method?: string;

    /**
     * URL
     * */
    url?: string;
  }

  interface NetSendOption {
    /**
     * Data to send
     * */
    data?: {
      [key: string]: any;
    };

    /**
     * Callback to receive data
     * */
    cb?: (response: NetResponse) => any;

    /**
     * if true, a loading will be display on top
     * */
    global?: boolean;
  }

  interface Net {
    send(config: NetSendConfig, option?: NetSendOption): void;
  }
}
