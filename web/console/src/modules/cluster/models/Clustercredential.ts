export interface Clustercredential {
  name: string;
  clusterName: string;
  caCert: string;
  token?: string;
  clientKey?: string;
  clientCert?: string;
}
