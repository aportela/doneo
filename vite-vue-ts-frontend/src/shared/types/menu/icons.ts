import type { Component } from "vue";

import {
  Home,
  FolderKanban,
  ListTodo,
  Users,
  UserKey,
  Settings,
  FolderCog,
  FileCog,
  Bookmark,
  Goal,
  Route,
  CircleUser,
  UserCog,
  LogOut,
  Search,
  Bell,
  BellOff,
  Sun,
  Moon,
  Notebook,
  Folder,
  FileText,
  PanelLeftOpen,
  PanelTopOpen,
} from "@lucide/vue";

export const MENU_ICON_SIZE = 20;

export const menuIcons: Record<string, Component> = {
  home: Home,

  workspace: Notebook,
  projects: FolderKanban,
  tasks: ListTodo,

  projectsGroup: Folder,

  users: Users,
  roles: UserKey,

  settings: Settings,

  projectSettings: FolderCog,
  taskSettings: FileCog,

  projectPages: FileText,

  projectTypes: Bookmark,
  priorities: Goal,
  statuses: Route,

  currentUser: CircleUser,
  profile: UserCog,
  logout: LogOut,

  search: Search,
  notifications: Bell,
  notificationsOff: BellOff,

  lightTheme: Sun,
  darkTheme: Moon,

  sideNavigation: PanelLeftOpen,

  topNavigation: PanelTopOpen,
};

export type MenuIconName = keyof typeof menuIcons;
