import { type AxiosResponse } from "axios";

import { type ValidAuthTypes } from "./common";

import type { UserInterface } from "./models/user";
import type { ProjectTypeInterface } from "./models/projectType";
import type { ProjectPriorityInterface } from "./models/projectPriority";

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

interface SearchUsersResponse extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    users: UserInterface[];
  };
}

interface GetUserResponse extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    user: UserInterface;
  };
}

interface SearchProjectTypesResponse
  extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    projectTypes: ProjectTypeInterface[];
  };
}

interface GetProjectTypeResponse extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    projectType: ProjectTypeInterface;
  };
}

interface SearchProjectPrioritiesResponse
  extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    projectPriorities: ProjectPriorityInterface[];
  };
}

interface GetProjectPriorityResponse
  extends Omit<DefaultAxiosResponse, "data"> {
  data: {
    projectPriority: ProjectPriorityInterface;
  };
}

export {
  type DefaultAxiosResponse,
  type SignInSucessResponse,
  type GetNewAccessTokenResponse,
  type SearchUsersResponse,
  type GetUserResponse,
  type SearchProjectTypesResponse,
  type GetProjectTypeResponse,
  type SearchProjectPrioritiesResponse,
  type GetProjectPriorityResponse,
};
