import { axiosInstance } from "../../../api/client";
import type {
  SignInRequest,
  SignInResponse,
  RenewAccessTokenResponse,
} from "../types/dto";

export const authService = {
  async signIn(payload: SignInRequest): Promise<SignInResponse> {
    const { data } = await axiosInstance.post<SignInResponse>(
      "/auth/signin",
      payload,
    );
    return data;
  },
  async signOut(): Promise<void> {
    await axiosInstance.post("/auth/signout");
  },
  async renewAccessToken(): Promise<RenewAccessTokenResponse> {
    const { data } = await axiosInstance.post<RenewAccessTokenResponse>(
      "/auth/renew-access-token",
    );
    return data;
  },
};
