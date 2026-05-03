import { defineStore, acceptHMRUpdate } from "pinia";
import { createStorageEntry } from "../composables/localStorage";

const localStorageUserSettingsFluidLayout = createStorageEntry<boolean>(
  "userSettings.fluidLayout",
  true,
);

const localStorageUserSettingsDisableNotifications =
  createStorageEntry<boolean>("userSettings.disableNotifications", false);

interface State {
  fluidLayout: boolean;
  disableNotifications: boolean;
}

export const useUserSettingsStore = defineStore("userSettings", {
  state: (): State => ({
    fluidLayout: localStorageUserSettingsFluidLayout.get(),
    disableNotifications: localStorageUserSettingsDisableNotifications.get(),
  }),
  getters: {
    hasFluidLayout: (state): boolean => state.fluidLayout === true,
    hasNotificationsEnabled: (state): boolean =>
      state.disableNotifications !== true,
  },
  actions: {
    setFluidLayout(fluid: boolean): void {
      this.fluidLayout = fluid;
      localStorageUserSettingsFluidLayout.set(this.fluidLayout);
    },
    toggleFluidLayout(): void {
      this.setFluidLayout(!this.fluidLayout);
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
