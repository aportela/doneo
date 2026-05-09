import { defineStore, acceptHMRUpdate } from "pinia";
import { AxiosError } from "axios";
import { type GetNewAccessTokenResponse as GetNewAccessTokenResponseInterface } from "../types/apiResponses";
import type { UserInterface } from "../types/models/user";
import { api } from "../composables/api";

interface State {
  session: {
    accessToken: {
      token: string | null;
      expiresAtTimestamp: number | null;
    };
    user: UserInterface | null;
  };
}

export const useSessionStore = defineStore("session", {
  state: (): State => ({
    session: {
      accessToken: {
        token: null,
        expiresAtTimestamp: null,
      },
      user: null,
    },
  }),
  getters: {
    sessionUserId: (state: State): string | null =>
      state.session.user?.id || null,
    sessionUserName: (state: State): string | null =>
      state.session.user?.name || null,
    hasAccessToken: (state: State): boolean =>
      state.session.accessToken.token !== null &&
      state.session.accessToken.expiresAtTimestamp !== null,
    accessToken: (state: State): string | null =>
      state.session.accessToken.token,
    accessTokenExpirationTimestamp: (state): number | null =>
      state.session.accessToken.expiresAtTimestamp,
  },
  actions: {
    async refreshAccessToken(): Promise<boolean> {
      try {
        const successResponse: GetNewAccessTokenResponseInterface =
          await api.auth.renewAccessToken();
        this.setAccessToken(
          successResponse.data.accessToken.token,
          successResponse.data.accessToken.expiresAtTimestamp,
        );
        this.setUser(successResponse.data.user);
        return true;
      } catch (e: unknown) {
        if (e instanceof AxiosError) {
          if (e.response?.status !== 401) {
            // 401 is the "normal" error response if we do not have a previous refresh token or refresh token is expired
            console.error("Invalid error http response code", e.status);
          }
        } else {
          console.error("Invalid error response", e);
        }
        return false;
      }
    },
    accessTokenExpiresBeforeInterval(interval: number) {
      if (this.session.accessToken.expiresAtTimestamp) {
        const currentTimestamp = Math.round(Date.now() / 1000);
        return (
          this.session.accessToken.expiresAtTimestamp - currentTimestamp <=
          interval
        );
      } else {
        return false;
      }
    },
    setAccessToken(token: string, expiresAtTimestamp: number): void {
      this.session.accessToken.token = token;
      this.session.accessToken.expiresAtTimestamp = expiresAtTimestamp;
    },
    removeAccessToken(): void {
      this.session.accessToken.token = null;
      this.session.accessToken.expiresAtTimestamp = null;
    },
    setUser(user: UserInterface): void {
      this.session.user = user;
    },
  },
});

export type SessionStoreType = ReturnType<typeof useSessionStore>;

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSessionStore, import.meta.hot));
}
