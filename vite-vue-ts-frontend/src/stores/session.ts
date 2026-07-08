import { defineStore, acceptHMRUpdate } from "pinia";
import { User } from "../modules/users/models/user";

interface State {
  session: {
    accessToken: {
      token: string | null;
      expiresAt: number | null;
    };
    user: User | null;
  };
}

export const useSessionStore = defineStore("session", {
  state: (): State => ({
    session: {
      accessToken: {
        token: null,
        expiresAt: null,
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
      state.session.accessToken.expiresAt !== null,
    accessToken: (state: State): string | null =>
      state.session.accessToken.token,
    accessTokenExpirationTimestamp: (state): number | null =>
      state.session.accessToken.expiresAt,
  },
  actions: {
    accessTokenExpiresBeforeInterval(secondsInterval: number) {
      if (this.session.accessToken.expiresAt) {
        const currentTimestamp = Date.now();
        return (
          this.session.accessToken.expiresAt - currentTimestamp <=
          secondsInterval * 1000
        );
      } else {
        return false;
      }
    },
    setAccessToken(token: string, expiresAtTimestamp: number): void {
      this.session.accessToken.token = token;
      this.session.accessToken.expiresAt = expiresAtTimestamp;
    },
    removeAccessToken(): void {
      this.session.accessToken.token = null;
      this.session.accessToken.expiresAt = null;
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
