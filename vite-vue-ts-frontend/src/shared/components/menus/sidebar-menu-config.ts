import type { AppMenuItem } from "../../types/menu/types";

export const sidebarMenuConfig: AppMenuItem[] = [
  {
    type: "route",
    key: "home",
    label: "Home",
    icon: "home",
    route: "home",
  },
  {
    type: "group",
    key: "workspace",
    label: "Workspace",
    icon: "workspace",
    children: [
      {
        type: "route",
        key: "manageProjects1",
        label: "Projects",
        route: "manageProjects",
        icon: "projects",
      },
      {
        type: "route",
        key: "manageTasks1",
        label: "Tasks",
        route: "manageTasks",
        icon: "tasks",
      },
    ],
  },
  {
    type: "divider",
    key: "divider1",
  },
  {
    type: "group",
    key: "projectGroup",
    label: "Projects",
    children: [
      {
        type: "group",
        key: "project1key",
        label: "Project1",
        icon: "project",
        children: [
          {
            type: "route",
            key: "project1tasks",
            label: "Tasks",
            icon: "tasks",
            route: "manageTasks",
          },
          {
            type: "route",
            key: "project1pages",
            label: "Pages",
            icon: "tasks",
            route: "manageTasks",
          },
        ],
      },
    ],
  },
  {
    type: "divider",
    key: "divider2",
  },
  {
    type: "group",
    key: "settingsgroup",
    label: "Settings",
    icon: "settings",
    children: [
      {
        type: "route",
        key: "users",
        label: "Users",
        route: "manageUsers",
        icon: "users",
      },
      {
        type: "route",
        key: "roles",
        label: "Roles",
        route: "manageRoles",
        icon: "roles",
      },
      {
        type: "route",
        key: "projectTypes",
        label: "Project types",
        route: "manageProjectTypes",
        icon: "projectTypes",
      },
      {
        type: "route",
        key: "projectPriorities",
        label: "Project priorities",
        route: "manageProjectPriorities",
        icon: "priorities",
      },
      {
        type: "route",
        key: "projectStatuses",
        label: "Project statuses",
        route: "manageProjectStatuses",
        icon: "statuses",
      },
    ],
  },
  {
    type: "divider",
    key: "divider3",
  },
];
