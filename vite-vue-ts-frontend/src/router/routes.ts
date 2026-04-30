import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/auth",
    component: () => import("../layouts/LoginRegisterLayout.vue"),
    children: [
      {
        name: "login",
        path: "login",
        component: () => import("../pages/LoginPage.vue"),
      },
    ],
  },
  {
    path: "/",
    name: "root",
    redirect: "/home",
    component: () => import("../layouts/MainLayout.vue"),
    children: [
      {
        name: "home",
        path: "home",
        component: () => import("../pages/HomePage.vue"),
      },
      {
        name: "projects",
        path: "projects",
        component: () => import("../pages/ProjectsPage.vue"),
      },
      {
        name: "workspaces",
        path: "workspaces",
        component: () => import("../pages/WorkSpacesPage.vue"),
      },
      {
        name: "tasks",
        path: "tasks",
        component: () => import("../pages/TasksPage.vue"),
      },
      {
        name: "reports",
        path: "reports",
        component: () => import("../pages/ReportsPage.vue"),
      },
      {
        name: "settings",
        path: "settings",
        component: () => import("../pages/SettingsPage.vue"),
      },
      {
        name: "users",
        path: "users",
        component: () => import("../pages/UsersPage.vue"),
      },
    ],
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: "/:catchAll(.*)*",
    name: "notFound",
    component: () => import("../layouts/ErrorNotFoundLayout.vue"),
  },
];

export default routes;
