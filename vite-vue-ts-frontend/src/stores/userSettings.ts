import { defineStore, acceptHMRUpdate } from "pinia";
import { createStorageEntry } from "../shared/composables/localStorage";
import { defaultDateTimeMask } from "../shared/composables/datetime";

type navigationMode = "top" | "side";

const localStorageUserSettingsNavigationMode =
  createStorageEntry<navigationMode>("userSettings.navigationMode", "side");

const localStorageUserSettingsDisableNotifications =
  createStorageEntry<boolean>("userSettings.disableNotifications", false);

const localStorageUserSettingsDatetimeMask = createStorageEntry<string>(
  "userSettings.datetimeMask",
  defaultDateTimeMask,
);

interface State {
  navigationMode: navigationMode;
  disableNotifications: boolean;
  datetimeMask: string;
}

export const useUserSettingsStore = defineStore("userSettings", {
  state: (): State => ({
    navigationMode: localStorageUserSettingsNavigationMode.get(),
    disableNotifications: localStorageUserSettingsDisableNotifications.get(),
    datetimeMask: localStorageUserSettingsDatetimeMask.get(),
  }),
  getters: {
    currentNavigationMode: (state): navigationMode => state.navigationMode,
    topNavigationMode: (state): boolean => state.navigationMode === "top",
    sideNavigationMode: (state): boolean => state.navigationMode === "side",
    hasNotificationsEnabled: (state): boolean =>
      state.disableNotifications !== true,
    currentDatetimeMask: (state): string => state.datetimeMask,
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
    setDatetimeMask(mask: string): void {
      this.datetimeMask = mask;
      localStorageUserSettingsDatetimeMask.set(this.datetimeMask);
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useUserSettingsStore, import.meta.hot),
  );
}
