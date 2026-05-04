import { type AxiosResponse } from "axios";

import { type ValidAuthTypes } from "./common";

import type { ProjectTypeInterface } from "./models/projectType";

interface DefaultAxiosResponse<T = unknown> {
  data: AxiosResponse<T>;
}

interface SignInSucessResponse extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    accessToken: {
      token: string;
      expiresAtTimestamp: number;
    };
    refreshToken: {
      token: string;
      expiresAtTimestamp: number;
    };
    tokenType: ValidAuthTypes;
  };
}

interface GetNewAccessTokenResponse extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    accessToken: {
      token: string;
      expiresAtTimestamp: number;
    };
    tokenType: ValidAuthTypes;
  };
}

interface SearchProjectTypesResponse extends Omit<
  DefaultAxiosResponse,
  "data"
> {
  data: {
    projectTypes: ProjectTypeInterface[];
  };
}

interface GetProjectTypeResponse extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    projectType: ProjectTypeInterface;
  };
}

export {
  type DefaultAxiosResponse,
  type SignInSucessResponse,
  type GetNewAccessTokenResponse,
  type SearchProjectTypesResponse,
  type GetProjectTypeResponse,
};
