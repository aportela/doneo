export interface APIErrorDetails {
  method: string;
  url: string;

  http: {
    code: number;
    status: string;
    contentType?: string;
  };

  request: {
    params?: unknown;
    data?: unknown;
  };

  response: unknown;
}
