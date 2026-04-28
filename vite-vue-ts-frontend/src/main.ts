import { createApp } from "vue";
import { router } from "./router/index";
import { createPinia } from "pinia";
import App from "./App.vue";
import AppNaiveUI from "./App-NaiveUI.vue";

import { useSessionStore } from "./stores/session";

const pinia = createPinia();

const useNaiveUI = true;

const app = createApp(useNaiveUI ? AppNaiveUI : App);
app.use(pinia);
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
    if (to.name === "login" || to.name === "register") {
      return { name: "home" };
    }
    return true;
  } else {
    if (to.name === "login" || to.name === "register") {
      return true;
    }
    return { name: "login" };
  }
});

app.use(router);
app.mount("#app");
