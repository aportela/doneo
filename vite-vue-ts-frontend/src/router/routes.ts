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
        name: "project",
        path: "project/:id",
        component: () => import("../pages/ProjectPage.vue"),
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
        name: "manageUsers",
        path: "manage/users",
        component: () => import("../pages/ManageUsersPage.vue"),
      },
      {
        name: "manageProjectTypes",
        path: "manage/project-types",
        component: () => import("../pages/ManageProjectTypesPage.vue"),
      },
      {
        name: "manageProjectStatuses",
        path: "manage/project-statuses",
        component: () => import("../pages/ManageProjectStatusesPage.vue"),
      },
      {
        name: "manageProjectPriorities",
        path: "manage/project-priorities",
        component: () => import("../pages/ManageProjectPrioritiesPage.vue"),
      },
      {
        name: "profile",
        path: "profile",
        component: () => import("../pages/ProfilePage.vue"),
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
