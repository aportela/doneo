import { defineStore, acceptHMRUpdate } from "pinia";

interface State {
  theme: Theme;
}

const storePersistenceKey = "doneo.settings.theme";

type Theme = "dark" | "light";

const systemTheme: Theme = window.matchMedia("(prefers-color-scheme: dark)")
  .matches
  ? "dark"
  : "light";

export const useThemeStore = defineStore("themeStore", {
  persist: {
    key: storePersistenceKey,
  },
  state: (): State => ({
    theme: systemTheme,
  }),
  getters: {
    light: (state): boolean => state.theme === "light",
    dark: (state): boolean => state.theme === "dark",
  },
  actions: {
    set(theme: Theme): void {
      this.theme = theme;
    },
    toggle(): void {
      this.set(this.theme === "light" ? "dark" : "light");
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useThemeStore, import.meta.hot));
}
