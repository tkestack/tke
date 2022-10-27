import { Request } from './request';

export async function modifyNamespaceDisplayName({
  name,

  displayName
}: {
  name: string;

  displayName: string;
}) {
  return Request.patch(
    `/apis/registry.tkestack.io/v1/namespaces/${name}`,
    {
      spec: {
        displayName
      }
    },
    {
      headers: {
        'Content-Type': 'application/merge-patch+json'
      }
    }
  );
}
