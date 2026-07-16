import { defineStore, acceptHMRUpdate } from "pinia";
import { defaultDateTimeMask } from "../shared/composables/datetime";

type Theme = "dark" | "light";

interface State {
  notifications: boolean;
  datetimeMask: string;
  theme: Theme;
}

const systemTheme: Theme = window.matchMedia("(prefers-color-scheme: dark)")
  .matches
  ? "dark"
  : "light";

const storePersistenceKey = "doneo.settings";

export const useUserSettingsStore = defineStore("userSettings", {
  persist: {
    key: storePersistenceKey,
  },
  state: (): State => ({
    notifications: true,
    datetimeMask: defaultDateTimeMask,
    theme: systemTheme,
  }),
  getters: {
    hasNotificationsEnabled: (state): boolean => state.notifications === true,
    currentDatetimeMask: (state): string => state.datetimeMask,
    lightTheme: (state): boolean => state.theme === "light",
    darkTheme: (state): boolean => state.theme === "dark",
  },
  actions: {
    setNotifications(enabled: boolean): void {
      this.notifications = enabled;
    },
    toggleNotifications(): void {
      this.setNotifications(!this.notifications);
    },
    setDatetimeMask(mask: string): void {
      this.datetimeMask = mask;
    },
    setTheme(theme: Theme): void {
      this.theme = theme;
    },
    toggleTheme(): void {
      this.setTheme(this.theme === "light" ? "dark" : "light");
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useUserSettingsStore, import.meta.hot),
  );
}
