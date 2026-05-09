export type TokenInterface = {
  token: string;
  expiresAtTimestamp: number;
};

export interface UserInterface {
  id: string;
  name: string;
  password?: string;
  email: string;
  isSuperUser: boolean;
  createdAt: number;
  updatedAt: number;
  deletedAt: number;
  avatar: string;
}

type ValidAuthTypes = "Bearer";

export type SignInResponseInterface = {
  accessToken: TokenInterface;
  refreshToken: TokenInterface;
  tokenType: ValidAuthTypes;
  user: UserInterface;
};

export type RenewAccessTokenResponseInterface = {
  accessToken: TokenInterface;
  tokenType: ValidAuthTypes;
  user: UserInterface;
};
