import { defineStore, acceptHMRUpdate } from "pinia";
import { defaultDateTimeMask } from "../shared/composables/datetime";

interface State {
  notifications: boolean;
  datetimeMask: string;
}

const storePersistenceKey = "doneo.settings";

export const useUserSettingsStore = defineStore("userSettings", {
  persist: {
    key: storePersistenceKey,
  },
  state: (): State => ({
    notifications: false,
    datetimeMask: defaultDateTimeMask,
  }),
  getters: {
    hasNotificationsEnabled: (state): boolean => state.notifications === true,
    currentDatetimeMask: (state): string => state.datetimeMask,
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
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(
    acceptHMRUpdate(useUserSettingsStore, import.meta.hot),
  );
}
