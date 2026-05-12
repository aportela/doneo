import { authService } from "./auth";
import { useSessionStore } from "../../../stores/session";
import type { RenewAccessTokenResponse } from "../types/dto";
import { User } from "../../users/models/user";
import axios from "axios";

export const TokenManager = {
  async refreshAccessToken(
    sessionStore: ReturnType<typeof useSessionStore>,
  ): Promise<boolean> {
    try {
      const response: RenewAccessTokenResponse =
        await authService.renewAccessToken();
      sessionStore.setAccessToken(
        response.accessToken.token,
        response.accessToken.expiresAtTimestamp,
      );
      sessionStore.setUser(new User(response.user));
      return true;
    } catch (error: unknown) {
      if (axios.isAxiosError(error)) {
        if (error.status === 401) {
          // normal error (no cookie || invalid refresh token)
        } else {
          console.error("Invalid API response code", error);
        }
      } else {
        console.error("Uncaught error while refreshing token", error);
      }
      return false;
    }
  },
};
