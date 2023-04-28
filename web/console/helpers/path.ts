export function getParamByUrl(key: string) {
  const searchParams = new URL(window.location.href)?.searchParams;
  return searchParams?.get(key);
}
