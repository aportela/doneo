import type { AxiosError } from "axios";

export type APIErrorResponse = {
  APIError?: boolean;
  code?: string;
  message: string;
};

export type APIError = AxiosError & {
  isAPIError: boolean;
  code: string;
  message: string;
};
