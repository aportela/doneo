import { defineStore, acceptHMRUpdate } from "pinia";
import { createStorageEntry } from "../shared/composables/localStorage";

type Theme = "dark" | "light";

const systemTheme: Theme = window.matchMedia("(prefers-color-scheme: dark)")
  .matches
  ? "dark"
  : "light";

const localStorageColorScheme = createStorageEntry<Theme>(
  "userSettings.colorScheme",
  systemTheme,
);

const savedScheme: Theme = localStorageColorScheme.get();

interface State {
  colorScheme: Theme;
}

export const useColorSchemeStore = defineStore("colorSchemeStore", {
  persist: {
    key: "doneo.settings.theme",
  },
  state: (): State => ({
    colorScheme: savedScheme,
  }),
  getters: {
    light: (state): boolean => state.colorScheme === "light",
    dark: (state): boolean => state.colorScheme === "dark",
  },
  actions: {
    set(scheme: Theme): void {
      this.colorScheme = scheme;
      localStorageColorScheme.set(this.colorScheme);
    },
    toggle(): void {
      this.set(this.colorScheme === "light" ? "dark" : "light");
    },
    initSystemListener(): void {
      const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
      const listener = (e: MediaQueryListEvent) => {
        this.set(e.matches ? "dark" : "light");
      };
      mediaQuery.addEventListener("change", listener);
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useColorSchemeStore, import.meta.hot));
}
