import { defineStore, acceptHMRUpdate } from "pinia";
import { createStorageEntry } from "../composables/localStorage";

type navigationMode = "top" | "side";

const localStorageUserSettingsNavigationMode =
  createStorageEntry<navigationMode>("userSettings.navigationMode", "side");

const localStorageUserSettingsDisableNotifications =
  createStorageEntry<boolean>("userSettings.disableNotifications", false);

interface State {
  navigationMode: navigationMode;
  disableNotifications: boolean;
}

export const useUserSettingsStore = defineStore("userSettings", {
  state: (): State => ({
    navigationMode: localStorageUserSettingsNavigationMode.get(),
    disableNotifications: localStorageUserSettingsDisableNotifications.get(),
  }),
  getters: {
    currentNavigationMode: (state): navigationMode => state.navigationMode,
    topNavigationMode: (state): boolean => state.navigationMode === "top",
    sideNavigationMode: (state): boolean => state.navigationMode === "side",
    hasNotificationsEnabled: (state): boolean =>
      state.disableNotifications !== true,
  },
  actions: {
    setNavigationMode(mode: navigationMode): void {
      this.navigationMode = mode;
      localStorageUserSettingsNavigationMode.set(this.navigationMode);
    },
    toggleNavigationMode(): void {
      this.setNavigationMode(this.navigationMode === "top" ? "side" : "top");
    },
    setNotifications(enabled: boolean): void {
      this.disableNotifications = !enabled;
      localStorageUserSettingsDisableNotifications.set(
        this.disableNotifications,
      );
    },
    toggleNotifications(): void {
      this.setNotifications(!this.hasNotificationsEnabled);
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useUserSettingsStore, import.meta.hot),
  );
}
