import { useLocalStorage } from "@vueuse/core";
import { LOCAL_STORAGE_NAMESPACE } from "../../constants";

export type StorageValue = string | number | boolean | null | object;

const createStorageEntry = <T extends StorageValue>(
  key: string,
  defaultValue: T,
) => {
  const storage = useLocalStorage<T>(
    LOCAL_STORAGE_NAMESPACE + key,
    defaultValue,
  );

  return {
    get(): T {
      return storage.value;
    },
    set(value: T) {
      storage.value = value;
    },
    remove() {
      storage.value = defaultValue;
    },
  };
};

export { createStorageEntry };
