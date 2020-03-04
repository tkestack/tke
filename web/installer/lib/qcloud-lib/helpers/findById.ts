import { Identifiable } from '../types/Identifiable';

export function findById<T extends Identifiable>(collection: T[], id: string | number): T {
  for (let i = 0; i < collection.length; i++) {
    if (collection[i].id === id) {
      return collection[i];
    }
  }
  return null;
}
