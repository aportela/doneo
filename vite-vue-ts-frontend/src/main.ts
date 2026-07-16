import { createApp } from "vue";
import { router } from "./router/index";
import { createPinia } from "pinia";
import piniaPluginPersistedstate from "pinia-plugin-persistedstate";
import { createI18n } from "vue-i18n";
import { messages } from "./i18n";
import App from "./App.vue";

import { useI18nStore } from "./stores/i18n";
import { useSessionStore } from "./stores/session";

import "./assets/styles/main.css";

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);

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

const i18nStore = useI18nStore();

const i18n = createI18n<{ message: MessageSchema }, MessageLanguages>({
  locale: i18nStore.currentLocale,
  legacy: false, // you must set `false`, to use Composition API
  globalInjection: true,
  messages,
});

app.use(i18n);

router.beforeEach(async (to, _from) => {
  if (!to.matched.length) {
    return { name: "notFound" };
  }
  const sessionStore = useSessionStore();
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
    return { name: "login", query: { redirect: to.fullPath } };
  }
});

app.use(router);
app.mount("#app");
