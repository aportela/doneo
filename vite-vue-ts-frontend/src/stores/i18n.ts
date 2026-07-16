import { defineStore, acceptHMRUpdate } from "pinia";
import { useNavigatorLanguage } from "@vueuse/core";
import { availableSystemLocales } from "../i18n";
import { DEFAULT_LOCALE } from "../constants";

const getMatchedLocale = (locale: string): string | null => {
  if (availableSystemLocales.includes(locale)) {
    return locale;
  } else if (locale.length >= 2) {
    // try similar match, example: locale es-MX (spanish, mexico) return es-ES (spanish, spain) if availableSystemLocales only contains [ 'en-US', 'es-ES', 'gl-GL']
    const shortLocale = locale.substring(0, 2).toLocaleLowerCase();
    const match = availableSystemLocales.find(
      (locale) => locale.substring(0, 2).toLocaleLowerCase() === shortLocale,
    );
    return match ?? null;
  } else {
    return null;
  }
};

interface State {
  locale: string;
}

const { isSupported, language } = useNavigatorLanguage();

const storePersistenceKey = "doneo.settings.i18n";

export const useI18nStore = defineStore("i18nStore", {
  persist: {
    key: storePersistenceKey,
  },
  state: (): State => ({
    locale:
      getMatchedLocale(
        getMatchedLocale(isSupported ? new String(language).toString() : "") ||
          "",
      ) ?? DEFAULT_LOCALE,
  }),
  getters: {
    currentLocale: (state: State): string => state.locale,
  },
  actions: {
    setLocale(locale: string): boolean {
      const matched = getMatchedLocale(locale);
      if (matched !== null) {
        this.locale = matched;
        return true;
      } else {
        return false;
      }
    },
  },
});

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useI18nStore, import.meta.hot));
}
