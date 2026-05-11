import type { UserResponse } from "../../users/types/dto";

export type Token = {
  token: string;
  expiresAtTimestamp: number;
};

type ValidAuthTypes = "Bearer";

export type SignInRequest = {
  email: string;
  password: string;
};

export type SignInResponse = {
  accessToken: Token;
  refreshToken: Token;
  tokenType: ValidAuthTypes;
  user: UserResponse;
};

export type RenewAccessTokenResponse = {
  accessToken: Token;
  tokenType: ValidAuthTypes;
  user: UserResponse;
};
