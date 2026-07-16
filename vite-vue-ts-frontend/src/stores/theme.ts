import { defineStore, acceptHMRUpdate } from "pinia";

type Theme = "dark" | "light";

const systemTheme: Theme = window.matchMedia("(prefers-color-scheme: dark)")
  .matches
  ? "dark"
  : "light";

interface State {
  theme: Theme;
}

export const useThemeStore = defineStore("themeStore", {
  persist: {
    key: "doneo.settings.theme",
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
    /*
    // TODO: REMOVE ?
    initSystemListener(): void {
      const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
      const listener = (e: MediaQueryListEvent) => {
        this.set(e.matches ? "dark" : "light");
      };
      mediaQuery.addEventListener("change", listener);
    },
    */
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useThemeStore, import.meta.hot));
}
