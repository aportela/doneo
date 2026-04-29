import { defineStore, acceptHMRUpdate } from "pinia";
import { createStorageEntry } from "../composables/localStorage";

type Schemes = "dark" | "light";

const systemTheme: Schemes = window.matchMedia("(prefers-color-scheme: dark)")
  .matches
  ? "dark"
  : "light";

const localStorageColorScheme = createStorageEntry<Schemes>(
  "colorScheme",
  systemTheme,
);

const savedScheme: Schemes = localStorageColorScheme.get();

interface State {
  colorScheme: Schemes;
}

export const useColorSchemeStore = defineStore("colorSchemeStore", {
  state: (): State => ({
    colorScheme: savedScheme,
  }),
  getters: {
    light: (state): boolean => state.colorScheme === "light",
    dark: (state): boolean => state.colorScheme === "dark",
  },
  actions: {
    set(scheme: Schemes): void {
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
