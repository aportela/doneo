import { axiosInstance } from "../client";
import type {
  SignInResponseInterface,
  RenewAccessTokenResponseInterface,
} from "../types/dto/auth";

export const authService = {
  async signIn(
    email: string,
    password: string,
  ): Promise<SignInResponseInterface> {
    const params = { email, password };
    const { data } = await axiosInstance.post<SignInResponseInterface>(
      "/auth/signin",
      params,
    );
    return data;
  },
  async signOut(): Promise<void> {
    await axiosInstance.post("/auth/signout");
  },
  async renewAccessToken(): Promise<RenewAccessTokenResponseInterface> {
    const { data } =
      await axiosInstance.post<RenewAccessTokenResponseInterface>(
        "/auth/renew-access-token",
      );
    return data;
  },
};
