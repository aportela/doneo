import type {
  AxiosError,
  InternalAxiosRequestConfig,
  AxiosHeaders,
  AxiosInstance,
} from "axios";
import type { SessionStoreType } from "../../stores/session";
import type { APIError, APIErrorResponse } from "./errors";

const getContentType = (
  headers: AxiosHeaders | Record<string, unknown> | undefined,
): string | undefined => {
  if (!headers) {
    return undefined;
  }

  if ("get" in headers && typeof headers.get === "function") {
    const ct = headers.get("content-type");
    return typeof ct === "string" ? ct : undefined;
  }

  const raw = (headers as Record<string, unknown>)["content-type"];

  if (Array.isArray(raw)) {
    return raw.find((v): v is string => typeof v === "string");
  }

  return typeof raw === "string" ? raw : undefined;
};

const initializedInstances = new WeakSet<AxiosInstance>();

export const setupInterceptors = (
  sessionStore: SessionStoreType,
  instance: AxiosInstance,
): void => {
  if (initializedInstances.has(instance)) {
    return;
  }

  initializedInstances.add(instance);

  instance.interceptors.request.use((config: InternalAxiosRequestConfig) => {
    if (sessionStore.hasAccessToken && sessionStore.accessToken) {
      config.headers.set?.(
        "Authorization",
        `Bearer ${sessionStore.accessToken}`,
      );
    }
    return config;
  });

  instance.interceptors.response.use(
    (response) => response,
    (error: AxiosError): Promise<APIError> => {
      if (error.response?.status === 401) {
        if (sessionStore.hasAccessToken) {
          sessionStore.removeAccessToken();
        }
      }

      const contentType = getContentType(error.response?.headers);
      const responseData = error.response?.data as APIErrorResponse | undefined;

      const apiError: APIError = {
        ...error,
        isAPIError: Boolean(
          contentType?.toLowerCase().includes("application/json") &&
            responseData?.APIError === true,
        ),
        code: responseData?.code || error.code || "",
        message: responseData?.message || error.message,
      };

      return Promise.reject(apiError);
    },
  );
};
