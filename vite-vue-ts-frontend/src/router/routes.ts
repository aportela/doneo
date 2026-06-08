import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/auth",
    component: () => import("../layouts/LoginLayout.vue"),
    children: [
      {
        name: "login",
        path: "login",
        component: () => import("../modules/auth/pages/LoginPage.vue"),
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
        name: "manageProjects",
        path: "projects",
        component: () =>
          import("../modules/projects/pages/ManageProjectsPage.vue"),
      },
      {
        name: "project",
        path: "projects/:projectId",
        component: () => import("../modules/projects/pages/ProjectPage.vue"),
        children: [
          {
            path: "tab/:tab",
            name: "projectTab",
            component: () =>
              import("../modules/projects/pages/ProjectPage.vue"),
          },
        ],
      },
      {
        name: "manageTasks",
        path: "tasks",
        component: () =>
          import("../modules/project-tasks/pages/ManageTasksPage.vue"),
      },
      {
        name: "task",
        path: "project/:projectId/tasks/:id",
        component: () => import("../modules/projects/pages/ProjectPage.vue"),
        children: [
          {
            path: "tab/:tab",
            name: "taskTab",
            component: () =>
              import("../modules/project-tasks/pages/TaskPage.vue"),
          },
        ],
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
        component: () => import("../modules/users/pages/ManageUsersPage.vue"),
      },
      {
        name: "manageRoles",
        path: "manage/roles",
        component: () => import("../modules/roles/pages/ManageRolesPage.vue"),
      },
      {
        name: "manageProjectTypes",
        path: "manage/project-types",
        component: () =>
          import("../modules/project-types/pages/ManageProjectTypesPage.vue"),
      },
      {
        name: "manageProjectStatuses",
        path: "manage/project-statuses",
        component: () =>
          import(
            "../modules/project-statuses/pages/ManageProjectStatusesPage.vue"
          ),
      },
      {
        name: "manageProjectPriorities",
        path: "manage/project-priorities",
        component: () =>
          import(
            "../modules/project-priorities/pages/ManageProjectPrioritiesPage.vue"
          ),
      },
      {
        name: "manageTaskPriorities",
        path: "manage/task-priorities",
        component: () =>
          import(
            "../modules/task-priorities/pages/ManageTaskPrioritiesPage.vue"
          ),
      },
      {
        name: "manageTaskStatuses",
        path: "manage/task-statuses",
        component: () =>
          import("../modules/task-statuses/pages/ManageTasktStatusesPage.vue"),
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
