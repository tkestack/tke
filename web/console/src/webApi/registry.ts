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

export async function fetchRepositoryList() {
  return Request.get<any, any>('/apis/registry.tkestack.io/v1/repositories');
}

export async function fetchNamespaceList() {
  return Request.get<any, any>('/apis/registry.tkestack.io/v1/namespaces');
}

export async function fetchRepoInfo(name: string) {
  return Request.get<any, any>(`/apis/registry.tkestack.io/v1/repositories/${name}`);
}
