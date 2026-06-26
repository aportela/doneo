import { defineStore, acceptHMRUpdate } from "pinia";
import { User } from "../modules/users/models/user";

interface State {
  session: {
    accessToken: {
      token: string | null;
      expiresAtTimestamp: number | null;
    };
    user: User | null;
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
    sessionUserEmail: (state: State): string | null =>
      state.session.user?.email || null,
    sessionUserIsAdmin: (state: State): boolean =>
      state.session.user?.permissions.isSuperUser ?? false,
    hasAccessToken: (state: State): boolean =>
      state.session.accessToken.token !== null &&
      state.session.accessToken.expiresAtTimestamp !== null,
    accessToken: (state: State): string | null =>
      state.session.accessToken.token,
    accessTokenExpirationTimestamp: (state): number | null =>
      state.session.accessToken.expiresAtTimestamp,
  },
  actions: {
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
    setUser(user: User): void {
      this.session.user = user;
    },
  },
});

export type SessionStoreType = ReturnType<typeof useSessionStore>;

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSessionStore, import.meta.hot));
}
