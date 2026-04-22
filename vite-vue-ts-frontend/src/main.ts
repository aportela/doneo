import { createApp } from "vue";
import "@tabler/core/dist/css/tabler.min.css";
import "@tabler/core/dist/js/tabler.min.js";

import { router } from "./router/index";
import App from "./App.vue";

createApp(App).use(router).mount("#app");
