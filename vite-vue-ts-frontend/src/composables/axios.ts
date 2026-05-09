import axios from "axios";
import { AxiosError } from "axios";
import { SERVER_API_BASE_PATH } from "../constants";
import { type APIErrorDetails } from "../types/apiErrorDetails";
import { useSessionStore } from "../stores/session";

interface AxiosAPIError extends AxiosError {
  isAPIError: boolean;
  customAPIErrorDetails: APIErrorDetails;
}

const axiosInstance = axios.create({
  baseURL: SERVER_API_BASE_PATH,
});

axiosInstance.interceptors.request.use(
  (config) => {
    const sessionStore = useSessionStore();
    if (sessionStore.hasAccessToken) {
      config.headers["Authorization"] = `Bearer ${sessionStore.accessToken}`;
      config.withCredentials = true;
    }
    return config;
  },
  (error) => {
    throw error;
  },
);

axiosInstance.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (!error) {
      throw new Error("Unknown API error");
    } else {
      if (!error.response) {
        error.response = {
          status: 0,
          statusText: "undefined",
        };
      }
      const sessionStore = useSessionStore();
      if (error.response?.status === 401) {
        if (sessionStore.hasAccessToken) {
          sessionStore.removeAccessToken();
        }
      }
      error.isAPIError =
        error.response.headers["content-type"] == "application/json" &&
        error.response.data.APIError === true;
      error.customAPIErrorDetails = {
        method: error.config?.method || "N/A",
        url: error.request?.responseURL || "N/A",
        httpCode: error.response?.status || "N/A",
        httpStatus: error.response?.statusText || "Unknown error",
        contentType: error.response.headers["content-type"],
        request: {
          params: {
            query: error.config.params || null,
            data: error.config.data || null,
          },
        },
        response: error.response.data || null,
      };
      throw error;
    }
  },
);

const bgDownload = async (url: string, fileName: string = "fileName") => {
  const startTime = Date.now();

  let urlBlob: string | null = null;
  let tmpLink: HTMLAnchorElement | null = null;

  try {
    const response = await axiosInstance.get(url, {
      responseType: "blob",
    });

    const blob = response.data;

    urlBlob = URL.createObjectURL(blob);

    tmpLink = document.createElement("a");
    tmpLink.href = urlBlob;
    tmpLink.download = fileName;

    document.body.appendChild(tmpLink);
    tmpLink.click();

    return {
      fileName,
      url,
      mimeType: blob.type,
      length: blob.size,
      msTime: Date.now() - startTime,
    };
  } finally {
    if (tmpLink) {
      document.body.removeChild(tmpLink);
    }

    if (urlBlob) {
      URL.revokeObjectURL(urlBlob);
    }
  }
};

export { axiosInstance, type AxiosAPIError, bgDownload };
