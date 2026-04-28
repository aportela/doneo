import { useLocalStorage } from "@vueuse/core";
import { LOCAL_STORAGE_NAMESPACE } from "../constants";

export type StorageValue = string | number | boolean | null | object;

const createStorageEntry = <T extends StorageValue>(
  key: string,
  defaultValue: T,
) => ({
  get(): T {
    const stored = localStorage.getItem(LOCAL_STORAGE_NAMESPACE + key);
    return stored === null ? defaultValue : (stored as T);
  },
  set(value: T) {
    useLocalStorage(LOCAL_STORAGE_NAMESPACE + key, value).value = value;
  },
  remove() {
    localStorage.removeItem(LOCAL_STORAGE_NAMESPACE + key);
  },
});

export { createStorageEntry };
