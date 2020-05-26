export function searchList<T>(list: T[], search: string, ...fields: string[]): T[] {
  if (!search) {
    return list;
  }

  const itemMatch = (x: T) =>
    fields.reduce((matched, field) => matched || String(x[field]).indexOf(search) > -1, false);

  return list.filter(itemMatch);
}
