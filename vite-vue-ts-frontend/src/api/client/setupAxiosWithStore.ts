import type { AxiosInstance } from "axios";
import type { SessionStoreType } from "../../stores/session";
import { setupInterceptors } from "./interceptors";
import { axiosInstance as baseInstance } from "./axios";

export const setupAxiosWithStore = (store: SessionStoreType): AxiosInstance => {
  const instance: AxiosInstance = baseInstance;
  setupInterceptors(store, instance);
  return instance;
};
