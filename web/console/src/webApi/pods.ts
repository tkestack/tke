import Request from './request';

interface FetchPodsProps {
  namespace: string;
  query: Record<string, any>;
}
export const fetchPods = ({ namespace, query }: FetchPodsProps) => {
  return Request.get(`/api/v1/namespaces/${namespace}/pods`, {
    params: query
  });
};
