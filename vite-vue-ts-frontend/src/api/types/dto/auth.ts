import type { User } from "./user";

export type TokenInterface = {
  token: string;
  expiresAtTimestamp: number;
};

type ValidAuthTypes = "Bearer";

export type SignInResponseInterface = {
  accessToken: TokenInterface;
  refreshToken: TokenInterface;
  tokenType: ValidAuthTypes;
  user: User;
};

export type RenewAccessTokenResponseInterface = {
  accessToken: TokenInterface;
  tokenType: ValidAuthTypes;
  user: User;
};
