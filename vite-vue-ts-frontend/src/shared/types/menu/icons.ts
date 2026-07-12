import type { Component } from "vue";

import {
  IconPresentation,
  IconUserCircle,
  IconBug,
  IconSitemap,
  IconFileAnalytics,
  IconUserCheck,
  IconSettings,
  IconUsers,
  IconChartHistogram,
  IconBookmark,
  IconFlagBolt,
  IconAdjustmentsBolt,
  IconLogout,
  IconId,
  IconSearch,
  IconBell,
  IconBellOff,
  IconMoon,
  IconSun,
  IconLayoutSidebarLeftExpand,
  IconLayoutNavbarExpand,
  IconDesk,
} from "@tabler/icons-vue";

export const MENU_ICON_SIZE = 20;

export const menuIcons: Record<string, Component> = {
  home: IconPresentation,

  workspace: IconDesk,
  projects: IconSitemap,
  tasks: IconBug,

  reports: IconFileAnalytics,
  charts: IconChartHistogram,

  users: IconUsers,
  roles: IconUserCheck,

  settings: IconSettings,

  projectTypes: IconBookmark,
  priorities: IconFlagBolt,
  statuses: IconAdjustmentsBolt,

  currentUser: IconUserCircle,
  profile: IconId,
  logout: IconLogout,

  search: IconSearch,
  notifications: IconBell,
  notificationsOff: IconBellOff,

  lightTheme: IconSun,
  darkTheme: IconMoon,

  sideNavigation: IconLayoutSidebarLeftExpand,

  topNavigation: IconLayoutNavbarExpand,
};

export type MenuIconName = keyof typeof menuIcons;
