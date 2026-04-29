import { createApp } from "vue";
import { router } from "./router/index";
import { createPinia } from "pinia";
import { createI18n } from "vue-i18n";
import { messages } from "./i18n";
import App from "./App.vue";

import { useI18nStore } from "./stores/i18n";
import { useSessionStore } from "./stores/session";

import { useUserSettingsStore } from "./stores/userSettings";

const pinia = createPinia();

const app = createApp(App);
app.use(pinia);

export type MessageLanguages = keyof typeof messages;
// Type-define 'en-US' as the master schema for the resource
export type MessageSchema = (typeof messages)["en-US"];

// See https://vue-i18n.intlify.dev/guide/advanced/typescript.html#global-resource-schema-type-definition
/* eslint-disable @typescript-eslint/no-empty-object-type */
declare module "vue-i18n" {
  // define the locale messages schema
  export interface DefineLocaleMessage extends MessageSchema {}

  // define the datetime format schema
  export interface DefineDateTimeFormat {}

  // define the number format schema
  export interface DefineNumberFormat {}
}
/* eslint-enable @typescript-eslint/no-empty-object-type */

const userSettingsStore = useUserSettingsStore();
userSettingsStore.init();

const i18nStore = useI18nStore();

const i18n = createI18n<{ message: MessageSchema }, MessageLanguages>({
  locale: i18nStore.currentLocale,
  legacy: false, // you must set `false`, to use Composition API
  globalInjection: true,
  messages,
});

app.use(i18n);

const sessionStore = useSessionStore();
const accessTokenCheckInterval = 300; // check every 5 min (300 seconds)

if (!sessionStore.hasAccessToken) {
  await sessionStore.refreshAccessToken();
}

const refreshAccessTokenIfNeeded = async (): Promise<boolean> => {
  if (
    sessionStore.hasAccessToken &&
    sessionStore.accessTokenExpiresBeforeInterval(accessTokenCheckInterval)
  ) {
    return await sessionStore.refreshAccessToken();
  } else {
    return false;
  }
};

// create interval for checking && auto-refresh access token
setInterval(() => {
  refreshAccessTokenIfNeeded()
    .then(() => {})
    .catch((e: Error) => {
      console.error(
        "An unhandled exception occurred during access token refresh",
        e,
      );
    })
    .finally(() => {});
}, accessTokenCheckInterval * 1000);

router.beforeEach((to, _from) => {
  if (!to.matched.length) {
    return { name: "notFound" };
  }
  const loggedIn = sessionStore.hasAccessToken;
  if (loggedIn) {
    if (to.name === "login") {
      return { name: "home" };
    }
    return true;
  } else {
    if (to.name === "login") {
      return true;
    }
    return { name: "login" };
  }
});

app.use(router);
app.mount("#app");
